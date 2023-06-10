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

type Business struct {
	ListenAddr string
	//(看自己把)用户 定义 启动 事件
	OnWorkerStart func(Worker workerman_go.Worker)
	//(必须处理)用户处理业务逻辑（从gateway转发过来的）
	OnMessage func(TcpConnection workerman_go.TcpConnection, msg []byte)
	//(没必要处理)当gateway或者sdk连接
	OnConnect func(conn workerman_go.TcpConnection)
	//(没必要处理)当gateway或者sdk断开
	OnClose func(conn workerman_go.TcpConnection)

	GatewayList map[uint64]*GatewayServer

	RegisterList map[uint64]*BusinessServer

	RegisterListRWMutex *sync.RWMutex

	Config workerman_go.ConfigGatewayWorker

	Name string
}
type BusinessServer struct {
	IpType workerman_go.IpType
	Ipv4   net.IPNet
	Ipv6   net.IPNet
	Port   workerman_go.Port
}

type GatewayServer struct {
	IpType workerman_go.IpType
	Ipv4   net.IPNet
	Ipv6   net.IPNet
	Port   workerman_go.Port
}

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

// 监听注册发现广播
func (b *Business) ConnectGatewayServerByRegisterBroadcast() {

	go func() {
		for true {

		}
	}()
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
			Host:   b.Config.RegisterPublicHostForRegister,
			Path:   workerman_go.RegisterBusniessWsPath,
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
		log.Fatalf("[%s]无法连接  注册发现 {%s%s} ", b.Name, Scheme, b.Config.RegisterPublicHostForRegister)
	}

	//registerServer:=BusinessServer{
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

// SendAuth 发送具有时效限制且有sign字段的json字符串
func (b *Business) SendAuth() {

	//dataMap := make(map[string]string)
	//dataMap[]
	//workerman_go.GenerateSignTime()
}
