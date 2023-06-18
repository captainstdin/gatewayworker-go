package business

import (
	"crypto/rand"
	"crypto/tls"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"math/big"
	"net"
	"net/url"
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

// InnerOnWorkerStart 启动后，应该连接register,获得gateway地址，然后去连接
func (b *Business) InnerOnWorkerStart(worker workerman_go.Worker) {

	Scheme := "wss://"
	if !b.Config.TLS {
		Scheme = "ws://"
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
	}

	b.registerMapRWMutex.Lock()

	wsRegister, wsConnWithRegisterErr := websocket.DialConfig(wsConfig)

	if wsConnWithRegisterErr != nil {
		for wsConnWithRegisterErr != nil {
			wsRegister, wsConnWithRegisterErr = websocket.DialConfig(wsConfig)
			log.Println("[%s]无法连接  注册发现 {%s%s}  10秒后重连.. ", b.Name, Scheme, b.Config.RegisterPublicHostForComponent)
			t := time.NewTicker(time.Second * 10)
			<-t.C
		}
	}

	//创建register实例
	RegisterConn := &ComponentRegister{
		addr:   b.Config.RegisterPublicHostForComponent,
		ConnWs: wsRegister,
		RWLock: &sync.RWMutex{},
	}

	ok := false
	for ok == false {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := b.registerMap[num.Uint64()]; !exist {
			//设置列表实例
			b.registerMap[num.Uint64()] = RegisterConn
			ok = true
		}
	}
	//解锁
	b.registerMapRWMutex.Unlock()

	//开始监听注册中心发来的指令和回复
}

// InnerOnConnect 当gateway连接
func (b *Business) InnerOnConnect(connection workerman_go.TcpConnection) {
	//TODO  gateway连接[发送身份请求，增加协程踢人定时器]
}

// InnerOnMessage 当 gateway 发送指令，需要回复gateway
func (b *Business) InnerOnMessage(connection workerman_go.TcpConnection, msg []byte) {
	//TODO   gateway 发送指令处理回复，进行业务处理

	//

	if b.OnMessage != nil {
		b.OnMessage(connection, msg)
	}
}

func (b *Business) InnerOnClose(connection workerman_go.TcpConnection) {
	//todo 当 gatewya断开的时候，需要进行规律重连 ，然后从列表踢掉

	if b.OnClose != nil {
		b.OnClose(connection)
	}
}
