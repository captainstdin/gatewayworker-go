package workerman_go

type TcpConnection interface {
	Close()

	Send(data interface{}) error

	GetRemoteIp() string

	GetRemotePort() string

	PauseRecv()
	ResumeRecv()

	Pipe(connection *TcpConnection)

	GetClientId() string

	GetClientIdInfo() *ClientToken

	Get(str string) (interface{}, bool)

	Set(str string, v interface{})
}
