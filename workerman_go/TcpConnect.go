package workerman_go

type TcpConnection interface {
	Close()

	Send(data interface{})

	getRemoteIp() string

	getRemotePort() string

	PauseRecv()
	ResumeRecv()

	Pipe(connection *TcpConnection)
}
