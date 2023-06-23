package business

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func (b *Business) Run() error {

	b.InnerOnWorkerStart(b)

	if b.OnWorkerStart != nil {
		b.OnWorkerStart(b)
	}

	return nil
}

// InnerOnWorkerStart 启动后，应该连接register,等待获得gateway地址，然后去连接gateway
func (b *Business) InnerOnWorkerStart(worker *Business) {

	Scheme := "wss"
	if !b.Config.TLS {
		Scheme = "ws"
	}

	// 设置WebSocket客户端配置
	wsConfig := &websocket.Config{
		Location: &url.URL{
			Scheme: Scheme,
			Host:   b.Config.RegisterPublicHostForComponent,
			Path:   workerman_go.RegisterForBusniessWsPath,
		},
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		TlsConfig: &tls.Config{InsecureSkipVerify: b.Config.SkipVerify},
		Version:   websocket.ProtocolVersionHybi13,
		Origin:    &url.URL{Scheme: "http", Host: b.Config.RegisterPublicHostForComponent},
	}

	b.registerMapRWMutex.Lock()

	wsRegister, wsConnWithRegisterErr := websocket.DialConfig(wsConfig)

	if wsConnWithRegisterErr != nil {
		for wsConnWithRegisterErr != nil {
			wsRegister, wsConnWithRegisterErr = websocket.DialConfig(wsConfig)

			log.Println(wsConnWithRegisterErr)
			log.Printf("[%s]无法连接  注册发现 {%s://%s%s}  10秒后重连.. ", b.Name, Scheme, b.Config.RegisterPublicHostForComponent, workerman_go.RegisterForBusniessWsPath)
			t := time.NewTicker(time.Second * 10)
			<-t.C
		}
	}

	//创建register实例
	RegisterConn := &ComponentRegister{
		root:     b,
		addr:     b.Config.RegisterPublicHostForComponent,
		ConnWs:   wsRegister,
		RWLock:   &sync.RWMutex{},
		ClientId: nil,
	}

	ok := false
	for ok == false {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := b.registerMap[num.Uint64()]; !exist {
			//设置列表实例
			RegisterConn.ClientId = &workerman_go.ClientToken{
				ClientGatewayNum: num.Uint64(),
			}
			b.registerMap[num.Uint64()] = RegisterConn
			ok = true
		}
	}

	//解锁
	b.registerMapRWMutex.Unlock()

	startInfo := bytes.Buffer{}
	startInfo.WriteByte('[')
	startInfo.WriteString(b.Name)
	startInfo.WriteString("] Starting  server with Connected ->【")
	startInfo.WriteString(b.Config.RegisterPublicHostForComponent)
	startInfo.WriteString("】 Listening...")
	log.Println(strconv.Quote(startInfo.String()))
	//阻塞式监听register消息
	RegisterConn.ListenMessageSync()

	//开始监听注册中心发来的指令和回复
}

// InnerOnConnect 当gateway连接
func (b *Business) InnerOnConnect(connection workerman_go.TcpWsConnection) {
	//TODO  gateway连接[发送身份请求，增加协程踢人定时器]
}

// InnerOnMessage 当 gateway 发送指令，需要回复gateway
func (b *Business) InnerOnMessage(connection workerman_go.TcpWsConnection, msg []byte) {
	//TODO   gateway 发送指令处理回复，进行业务处理

	if b.OnMessage != nil {
		b.OnMessage(connection, msg)
	}
}

func (b *Business) InnerOnClose(connection workerman_go.TcpWsConnection) {
	//todo 当 gatewya断开的时候，需要进行规律重连 ，然后从列表踢掉

	if b.OnClose != nil {
		b.OnClose(connection)
	}

}
