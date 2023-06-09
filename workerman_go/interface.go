package workerman_go

type Worker interface {
	//启动服务
	Run() error
}
