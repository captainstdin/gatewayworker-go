package business

import (
	"crypto/tls"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-make(chan struct{})
	}()

	//阻塞当前线程
	wg.Wait()
	return nil
}

func (b *Business) InnerOnWorkerStart(worker workerman_go.Worker) {
	//todo 连接 注册发现

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
	}

	b.RegisterListRWMutex.Lock()
	defer b.RegisterListRWMutex.Unlock()
	_, wsConnWithRegisterErr := websocket.DialConfig(wsConfig)

	if wsConnWithRegisterErr != nil {
		log.Fatalf("[%s]无法连接  注册发现 {%s%s} ", b.Name, Scheme, b.Config.RegisterPublicHostForComponent)
	}

	//registerServer:=Server{
	//	IpType: wsConnWithRegister.RemoteAddr(),
	//	Ipv4:   net.IPNet{},
	//	Ipv6:   net.IPNet{},
	//	Port:   0,
	//}
	//b.RegisterList[]
}

// InnerOnConnect 当gateway连接
func (b *Business) InnerOnConnect(connection workerman_go.TcpConnection) {
	//TODO  sdk ||gateway连接[发送身份请求，增加协程踢人定时器]
}

// InnerOnMessage 当 gateway 发送指令，需要回复gateway
func (b *Business) InnerOnMessage(connection workerman_go.TcpConnection, msg []byte) {
	//TODO   gateway 发送指令处理回复，进行业务处理
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
