package gateway

import (
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
)

type UserClient struct {
	root *GatewayServer
	//生成的在当前内部组件中标志目标gateway所在地
	ClientId *workerman_go.ClientToken

	//组件名称
	Name string

	//组件类型
	ComponentType int

	//连接地址
	Address string
	Port    workerman_go.Port

	FdWs *websocket.Conn
}

func (u UserClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) Send(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) GetRemoteIp() string {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) Pipe(connection *workerman_go.TcpWsConnection) {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (u UserClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
