package workerman_go

type TcpConnection interface {
	Close()

	Send(data interface{})

	GetRemoteIp() string

	GetRemotePort() string

	PauseRecv()
	ResumeRecv()

	Pipe(connection *TcpConnection)
}
