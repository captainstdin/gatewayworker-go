package gateway

import (
	"gatewaywork-go/workerman_go"
	"net"
)

// GatewayClient 每个连接上来的ws Client主要是  component组件(business)与 WebSocket用户
type GatewayClient struct {
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

func (g *GatewayClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) Send(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) GetRemoteIp() string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) Pipe(connection *workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
