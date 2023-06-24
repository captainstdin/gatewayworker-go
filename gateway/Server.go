package gateway

import (
	"context"
	"gatewaywork-go/workerman_go"
	"sync"
)

type Server struct {
	workerman_go.Worker

	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]*workerman_go.AsyncTcpWsConnection //key是remoteAddress

	//已连接的Gateway
	ConnectedBusinessLock *sync.RWMutex
	ConnectedBusinessMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

}

func (s *Server) Run() {

}

func (s *Server) OnWorkerStart(server *workerman_go.Worker) {

}

func (s *Server) OnConnect(connection workerman_go.InterfaceConnection) {

}
func (s *Server) OnMessage(connection workerman_go.InterfaceConnection, buff []byte) {

}
func (s *Server) OnClose(connection workerman_go.InterfaceConnection) {

}

func NewServer(name string, conf *workerman_go.ConfigGatewayWorker) *Server {

	ctx, cf := context.WithCancel(context.Background())
	w := workerman_go.Worker{
		Connections:     map[uint64]workerman_go.InterfaceConnection{},
		ConnectionsLock: &sync.RWMutex{},
		ListenAddress:   conf.RegisterListenAddr,
		ListenPath:      workerman_go.RegisterForComponent,
		Name:            name,
		Tls:             false,
		TlsPem:          "",
		TlsKey:          "",
		OnWorkerStart:   nil,
		OnConnect:       nil,
		OnMessage:       nil,
		OnClose:         nil,
		Ctx:             ctx,
		CtxF:            cf,
		Config:          conf,
	}

	server := &Server{
		Worker:                w,
		ConnectedBusinessLock: &sync.RWMutex{},
		ConnectedRegisterLock: &sync.RWMutex{},
		ConnectedRegisterMap:  make(map[string]*workerman_go.AsyncTcpWsConnection),
		ConnectedBusinessMap:  make(map[string]workerman_go.InterfaceConnection),
	}
	server.Worker.OnWorkerStart = server.OnWorkerStart
	server.Worker.OnConnect = server.OnConnect
	server.Worker.OnMessage = server.OnMessage
	server.Worker.OnClose = server.OnClose
	return server
}
