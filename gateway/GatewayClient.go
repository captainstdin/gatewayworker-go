package gateway

import (
	"gatewaywork-go/workerman_go"
)

type GatewayClient struct {
}

func (g GatewayClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) Send(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) GetRemoteIp() string {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) Pipe(connection *workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (g GatewayClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
