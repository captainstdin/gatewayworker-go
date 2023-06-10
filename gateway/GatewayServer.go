package gateway

import (
	"gatewaywork-go/workerman_go"
	"sync"
)

// GatewayServer 需要有公网端口2727，那些人会连2727？{Sdk:http://:2727/action,Business:ws://:2727}
type GatewayServer struct {
	//监听地址，可以是:8080，同时监听ipv6 4
	ListenAddr string
	//(非必填)启动服务器
	OnWorkerStart func(Worker *GatewayServer)

	//(非必填)Sdk || Business || Client 连接上来
	OnConnect func(conn workerman_go.TcpConnection)
	//(非必填)SSdk || Business || Client 发送指令
	OnMessage func(Worker workerman_go.TcpConnection, msg []byte)

	//(非必填)SSdk || Business || Client 断开连接
	OnClose func(Worker workerman_go.TcpConnection)

	//服务配置
	GatewayWorkerConfig workerman_go.ConfigGatewayWorker

	//记录了 Sdk || Business || Client
	ConnectionList map[uint64]*workerman_go.TcpConnection
	//读写锁
	ConnectionListRWLock *sync.RWMutex

	Name string
}

func (g *GatewayServer) Run() error {
	//todo 1：开启监听本地ws://:2727 ,http://:2727

	//http://启动一个gin服务器，

	//todo 2：准备完毕，连接Register，并且
	return nil
}

// InnerOnWorkerStart 内置服务启动回调
func (g *GatewayServer) InnerOnWorkerStart(worker workerman_go.Worker) {
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
}
