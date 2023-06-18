package workerman_go

type Worker interface {
	//启动服务
	Run() error

	OnWorkerStart(worker Worker)

	OnConnect(connection TcpConnection)

	OnMessage(connection TcpConnection, msg []byte)

	OnClose(connection TcpConnection)
}
