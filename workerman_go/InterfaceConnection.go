package workerman_go

import (
	"context"
	"github.com/gorilla/websocket"
	"net"
)

type InterfaceConnection interface {
	Close()

	Send(data interface{}) error

	GetRemoteIp() (net.IP, error)

	GetRemotePort() (uint16, error)

	Pipe(connection *TcpConnection)

	//PauseRecv ResumeRecv 暂未实现
	PauseRecv()
	ResumeRecv()

	GetClientId() string

	GetClientIdInfo() *ClientToken

	Get(str string) (interface{}, bool)

	Set(str string, v interface{})

	GotCtxWithF() (context.Context, context.CancelFunc)

	GotFd() *websocket.Conn
}
