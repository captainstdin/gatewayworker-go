package business

import (
	"gatewaywork-go/workerman_go"
	"net"
	"sync"
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

	//GatewayList 组件-网关-列表
	GatewayList map[uint64]*GatewayServer

	//RegisterList 组件-业务处理-列表-并发注册注销锁
	RegisterList map[uint64]*ComponentClient

	//GatewayList 组件-网关-并发注册注销锁
	GatewayListRWMutex *sync.RWMutex

	//RegisterList 组件-业务处理-并发注册注销锁
	RegisterListRWMutex *sync.RWMutex

	Config workerman_go.ConfigGatewayWorker

	Name string
}

type GatewayServer struct {
	IpType workerman_go.IpType
	Ipv4   net.IPNet
	Ipv6   net.IPNet
	Port   workerman_go.Port
}

// 监听注册发现广播
func (b *Business) ConnectGatewayServerByRegisterBroadcast() {
	go func() {

	}()
}

// SendAuth 发送具有时效限制且有sign字段的json字符串
func (b *Business) SendAuth() {

	//dataMap := make(map[string]string)
	//dataMap[]
	//workerman_go.GenerateSignTime()
}
