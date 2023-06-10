package workerman_go

type Worker interface {
	//启动服务
	Run() error

	InnerOnWorkerStart(worker Worker)

	InnerOnConnect(connection TcpConnection)

	InnerOnMessage(connection TcpConnection, msg []byte)

	InnerOnClose(connection TcpConnection)
}
