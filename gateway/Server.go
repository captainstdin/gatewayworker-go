package gateway

import (
	"context"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"sync"
	"time"
)

type Server struct {
	workerman_go.Worker

	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

	//已连接的Gateway
	ConnectedBusinessLock *sync.RWMutex
	ConnectedBusinessMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

}

func (s *Server) SendSignData(data any, conn workerman_go.InterfaceConnection) {

	var CMDInt int
	switch data.(type) {
	case workerman_go.ProtocolRegister:
		CMDInt = workerman_go.CommandComponentAuthRequest
	case workerman_go.ProtocolPublicGatewayConnectionInfo:

	}

	timeByte, err := workerman_go.GenerateSignTimeByte(CMDInt, data, s.Config.SignKey, func() time.Duration {
		return time.Second * 120
	})
	if err != nil {
		return
	}
	conn.Send(timeByte.ToByte())
}

func (s *Server) Run() {

	business := workerman_go.NewAsyncTcpWsConnection(s.Config.RegisterPublicHostForComponent)

	business.OnConnect = func(connection *workerman_go.AsyncTcpWsConnection) {
		s.ConnectedBusinessLock.Lock()
		s.ConnectedBusinessMap[connection.GetRemoteAddress()] = connection
		s.ConnectedBusinessLock.Unlock()
	}

	business.OnMessage = func(connection *workerman_go.TcpWsConnection, buff []byte) {
		//sdk或者指令
		Data, parse := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
		if parse != nil {
			return
		}
		switch Data.Cmd {
		case workerman_go.CommandComponentAuthRequest:
			var reg workerman_go.ProtocolRegister
			err := json.Unmarshal(buff, &reg)
			if err != nil {
				return
			}

			if reg.Authed == "1" {
				s.ConnectedRegisterLock.RLock()
				s.ConnectedRegisterMap[connection.GetRemoteAddress()] = connection
				s.ConnectedRegisterLock.Unlock()
				return
			}

			s.SendSignData(workerman_go.ProtocolRegister{}, connection)

		}
	}

	business.OnClose = func(connection *workerman_go.TcpWsConnection) {
	}

	go business.Connect()

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
		ConnectedRegisterMap:  make(map[string]workerman_go.InterfaceConnection),
		ConnectedBusinessMap:  make(map[string]workerman_go.InterfaceConnection),
	}
	server.Worker.OnWorkerStart = server.OnWorkerStart
	server.Worker.OnConnect = server.OnConnect
	server.Worker.OnMessage = server.OnMessage
	server.Worker.OnClose = server.OnClose
	return server
}
