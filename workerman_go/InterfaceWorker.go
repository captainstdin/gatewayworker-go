package workerman_go

type InterfaceWorker interface {
	//Run 启动服务
	Run() error

	onWorkerStart(worker InterfaceWorker)

	onConnect(connection InterfaceConnection)

	onMessage(connection InterfaceConnection, msg []byte)

	onClose(connection InterfaceConnection)
}
