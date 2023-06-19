package business

import (
	"gatewaywork-go/workerman_go"
	"sync"
)

type Business struct {
	ListenAddr string
	//(看自己把)用户 定义 启动 事件
	OnWorkerStart func(Worker *Business)
	//(必须处理)用户处理业务逻辑（从gateway转发过来的）
	OnMessage func(TcpConnection workerman_go.TcpConnection, msg []byte)
	//(没必要处理)当gateway或者sdk连接
	OnConnect func(conn workerman_go.TcpConnection)
	//(没必要处理)当gateway或者sdk断开
	OnClose func(conn workerman_go.TcpConnection)
	//gatewayMapRWMutex 组件-网关-并发注册注销锁
	gatewayMapRWMutex *sync.RWMutex
	//GatewayList 组件-网关-列表，网关是公网的，一定是唯一的连接
	gatewayMap map[string]*ComponentGateway

	//RegisterList 组件-业务处理-列表-并发注册注销锁
	registerMap map[uint64]*ComponentRegister
	//registerMapRWMutex 组件-业务处理-并发注册注销锁
	registerMapRWMutex *sync.RWMutex

	//集群配置模块
	Config *workerman_go.ConfigGatewayWorker

	//服务名
	Name string
}

// ConnectGatewayServerByRegisterBroadcast 监听注册发现广播
func (b *Business) ConnectGatewayServerByRegisterBroadcast() {
	go func() {

	}()
}

func NewBusiness(name string, Conf *workerman_go.ConfigGatewayWorker) *Business {

	return &Business{
		ListenAddr:         "",
		OnWorkerStart:      nil,
		OnMessage:          nil,
		OnConnect:          nil,
		OnClose:            nil,
		gatewayMapRWMutex:  &sync.RWMutex{},
		gatewayMap:         make(map[string]*ComponentGateway, 0),
		registerMap:        make(map[uint64]*ComponentRegister, 0),
		registerMapRWMutex: &sync.RWMutex{},
		Config:             Conf,
		Name:               name,
	}
}
