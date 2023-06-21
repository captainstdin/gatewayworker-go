package gateway

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

	ComponentsMap map[uint64]*ComponentClient

	ComponentsLock sync.RWMutex

	Register map[uint64]*ComponentClient

	Name string
}

func (g *GatewayServer) Run() error {

	g.RunGinServer(g.Config.GatewayListenAddr, g.Config.GatewayListenPort)

	//启动，连接register
	g.onWorkerStart(g)

	return nil
}

func (g *GatewayServer) onConnectForward(connection workerman_go.TcpConnection) {

	//todo 转发到哈希路由business上
}

func (g *GatewayServer) onCloseForward(connection workerman_go.TcpConnection) {

	//todo 转发到哈希路由business上
}

// onMessageForward 用户连接上来的时候，原生msg，加密后转发给business
func (g *GatewayServer) onMessageForward(connection workerman_go.TcpConnection, msg []byte) {

	//todo 转发到哈希路由business上
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

	//注册服务中心实例

	wsRegister, wsConnWithRegisterErr := websocket.DialConfig(wsConfig)

	if wsConnWithRegisterErr != nil {
		for wsConnWithRegisterErr != nil {
			wsRegister, wsConnWithRegisterErr = websocket.DialConfig(wsConfig)

			log.Println(wsConnWithRegisterErr)
			log.Printf("[%s]无法连接  注册发现 {%s://%s%s}  10秒后重连.. ", g.Name, Scheme, g.Config.RegisterPublicHostForComponent, workerman_go.RegisterForBusniessWsPath)
			t := time.NewTicker(time.Second * 10)
			<-t.C
		}
	}

	//锁
	g.ComponentsLock.Lock()

	//成功连接register
	RegisterInstance := &ComponentClient{
		root:    g,
		FdWs:    wsRegister,
		Address: g.Config.RegisterPublicHostForComponent,
		Port:    0,
		ClientId: &workerman_go.ClientToken{
			IPType:            0,
			ClientGatewayIpv4: nil,
			ClientGatewayIpv6: nil,
			ClientGatewayPort: 0,
			ClientGatewayNum:  getUniqueKey(g.ComponentsMap).Uint64(),
		},
		Name:          "Register", //临时起名
		ComponentType: workerman_go.ComponentIdentifiersTypeBusiness,
	}

	g.Register[RegisterInstance.ClientId.ClientGatewayNum] = RegisterInstance

	g.ComponentsLock.Unlock()

	//监听
	g.Register[RegisterInstance.ClientId.ClientGatewayNum].Listen(func(sign *workerman_go.GenerateComponentSign, client *ComponentClient) {
		//处理下脸上register的时候，发送注册
		switch sign.Cmd {
		case workerman_go.CommandComponentAuthRequest:
			//如果要求认证，就返回
			client.Send(workerman_go.ProtocolRegister{
				ComponentType:                       workerman_go.ComponentIdentifiersTypeGateway,
				Name:                                g.Name,
				ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
				Data:                                "gateway.auth",
				Authed:                              "0",
			})
		}
	})

	RegisterInstance.Send(workerman_go.ProtocolRegister{
		ComponentType:                       workerman_go.ComponentIdentifiersTypeGateway,
		Name:                                g.Name,
		ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
		Data:                                "gateway.auth",
		Authed:                              "0",
	})

}

func getUniqueKeyByUserClient(mapData map[uint64]*workerman_go.TcpConnection) *big.Int {
	for {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := mapData[num.Uint64()]; !exist {
			//设置列表实例
			return num
		}
	}

}

// 获取唯一key
func getUniqueKey(mapData map[uint64]*ComponentClient) *big.Int {
	for {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := mapData[num.Uint64()]; !exist {
			//设置列表实例
			return num
		}
	}
}

// OnWorkerStart 内置服务启动回调
func (g *GatewayServer) OnWorkerStart(worker workerman_go.Worker) {
	//TODO implement me
	panic("implement me")
}

// InnerOnConnect 内置Sdk || Business  连接上来
func (g *GatewayServer) InnerOnConnect(connection workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

// InnerOnMessage 内置Sdk || Business  消息或者指令
func (g *GatewayServer) InnerOnMessage(connection *ComponentClient, msg []byte) {
	//TODO implement me

	//todo 写一个定时器30  秒后验证关闭未验证的business
	go func(com *ComponentClient) {
		timer := time.NewTimer(30 * time.Second)
		for true {
			select {
			case <-timer.C:
				g.ComponentsLock.RLock()
				com.onClose(com)

				if _, ok := g.ComponentsMap[com.ClientId.ClientGatewayNum]; ok == true {
					delete(g.ComponentsMap, com.ClientId.ClientGatewayNum)
				}
				g.ComponentsLock.RUnlock()
			}
		}
	}(connection)

}

func (g *GatewayServer) InnerOnClose(connection workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")

	//todo 判断是 wsClient还是Business用户

	//todo 当client_id下线（连接断开）时会自动与uid解绑，开发者无需在onClose事件调用Gateway::unbindUid。
}
