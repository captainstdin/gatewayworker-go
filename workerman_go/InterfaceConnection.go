package workerman_go

import (
	"net"
)

type InterfaceConnection interface {
	Close()

	Send(data interface{}) error

	GetRemoteAddress() string

	GetRemoteIp() (net.IP, error)

	GetRemotePort() (uint16, error)

	Pipe(connection *TcpWsConnection)

	//PauseRecv ResumeRecv 暂未实现
	PauseRecv()
	ResumeRecv()

	GetClientId() string

	GetClientIdInfo() *GatewayIdInfo

	Get(str string) (interface{}, bool)

	Set(str string, v interface{})

	Worker() *Worker

	TcpWsConnection() *TcpWsConnection
}
