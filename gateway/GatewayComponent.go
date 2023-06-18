package gateway

import (
	"gatewaywork-go/workerman_go"
	"net"
)

// ComponentClient 每个连接上来的ws Client主要是  component组件(business)与 WebSocket用户
type ComponentClient struct {
	ClientId string
	//是否是用户,true是用户，
	IsClient bool

	//session 设置
	Session *workerman_go.SessionKv

	//ip类型
	IpType workerman_go.IpType

	//ip4
	Ipv4 net.IP

	//ip6
	Ipv6 net.IP

	Port workerman_go.Port

	//生成的在当前内部组件中标志目标gateway所在地
	ClientToken *workerman_go.ClientToken
}

func (g *ComponentClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Send(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetRemoteIp() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Pipe(connection *workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
