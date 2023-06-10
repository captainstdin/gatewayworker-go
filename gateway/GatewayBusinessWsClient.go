package gateway

import (
	"gatewaywork-go/workerman_go"
	"net"
)

// GatewayBusinessClient Business客户端,主要是把gin.ws -> business
type GatewayBusinessClient struct {
	//session 设置，应该没啥用
	Session *workerman_go.SessionKv

	//ip类型
	IpType workerman_go.IpType

	//ip4
	Ipv4 net.IP

	//ip6
	Ipv6 net.IP

	//端口号
	Port workerman_go.Port

	//实体
	FdWs net.Conn

	//root
	GatewayServer *GatewayServer
	//生成的在当前内部组件中标志目标gateway所在地
	ClientToken *workerman_go.ClientToken
}

func (bc *GatewayBusinessClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) Send(data interface{}) error {
	//TODO implement me

	_, err := bc.FdWs.Write(data.([]byte))
	if err != nil {
		bc.GatewayServer.InnerOnClose(bc)

		return err
	}
	return nil
}

func (bc *GatewayBusinessClient) GetRemoteIp() string {
	//TODO implement me

	return ""
}

func (bc *GatewayBusinessClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) Pipe(connection *workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (bc *GatewayBusinessClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
