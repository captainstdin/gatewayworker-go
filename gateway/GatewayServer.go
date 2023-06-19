package gateway

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

// GatewayServer 需要有公网端口2727，那些人会连2727？{Sdk:http://:2727/action,Business:ws://:2727, Client:ws://:2727}
type GatewayServer struct {
	//监听地址，可以是:8080，同时监听ipv6 4
	ListenAddr string

	//服务配置
	Config workerman_go.ConfigGatewayWorker

	//记录了 Client
	ConnectionMap map[uint64]*workerman_go.TcpConnection
	//读写锁
	ConnectionMapRWLock *sync.RWMutex

	ComponentClient map[uint64]*ComponentClient

	Register map[uint64]*ComponentClient

	Name string
}

func (g *GatewayServer) Run() error {

	g.onWorkerStart(g)

	//todo 1-2: http://启动一个gin服务器，

	g.RunGinServer(g.Config.GatewayListenAddr, g.Config.GatewayListenPort)

	//todo 2：准备完毕，连接Register，并且
	return nil
}

// onWorkerStart 连接到register注册中心
func (g *GatewayServer) onWorkerStart(worker *GatewayServer) {
	//TODO implement me

	Scheme := "wss"
	if !g.Config.TLS {
		Scheme = "ws"
	}

	// 设置WebSocket客户端配置
	wsConfig := &websocket.Config{
		Location: &url.URL{
			Scheme: Scheme,
			Host:   g.Config.RegisterPublicHostForComponent,
			Path:   workerman_go.RegisterForBusniessWsPath,
		},
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		TlsConfig: &tls.Config{InsecureSkipVerify: g.Config.SkipVerify},
		Version:   websocket.ProtocolVersionHybi13,
		Origin:    &url.URL{Scheme: "http", Host: g.Config.RegisterPublicHostForComponent},
	}

	g.com.Lock()

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

}

// OnWorkerStart 内置服务启动回调
func (g *GatewayServer) OnWorkerStart(worker workerman_go.Worker) {
	//TODO implement me
	panic("implement me")
}

// InnerOnConnect 内置Sdk || Business || Client 连接上来
func (g *GatewayServer) InnerOnConnect(connection workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

// InnerOnMessage 内置Sdk || Business || Client 消息或者指令
func (g *GatewayServer) InnerOnMessage(connection workerman_go.TcpConnection, msg []byte) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) InnerOnClose(connection workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")

	//todo 判断是 wsClient还是Business用户

	//todo 当client_id下线（连接断开）时会自动与uid解绑，开发者无需在onClose事件调用Gateway::unbindUid。
}
