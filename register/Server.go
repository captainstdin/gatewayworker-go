package register

import (
	"bytes"
	"context"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	keyGatewayLanInfo = "GatewayLanInfo"
)

type Server struct {
	workerman_go.Worker

	_workerConnections map[uint64]workerman_go.InterfaceConnection

	_gatewayConnections map[uint64]workerman_go.InterfaceConnection
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
		Worker:              w,
		_workerConnections:  make(map[uint64]workerman_go.InterfaceConnection),
		_gatewayConnections: make(map[uint64]workerman_go.InterfaceConnection),
	}
	server.Worker.OnWorkerStart = server.OnWorkerStart
	server.Worker.OnConnect = server.OnConnect
	server.Worker.OnMessage = server.OnMessage

	server.Worker.OnClose = server.OnClose
	return server
}

func (s *Server) OnWorkerStart(worker *workerman_go.Worker) {
	startInfo := bytes.Buffer{}
	startInfo.WriteByte('[')
	startInfo.WriteString(worker.Name)
	startInfo.WriteString("] Starting  server at  ->【")
	startInfo.WriteString(worker.ListenAddress)
	startInfo.WriteString(worker.ListenPath)
	startInfo.WriteString("】 Listening...")
	log.Println(strconv.Quote(startInfo.String()))
}

func (s *Server) OnConnect(conn workerman_go.InterfaceConnection) {
	//非阻塞
	SendSignData(workerman_go.ProtocolRegister{
		ComponentType:                       0,
		Name:                                "",
		ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
		Data:                                "Business.request.auth",
		Authed:                              "0",
	}, conn)

}

func (s *Server) OnMessage(conn workerman_go.InterfaceConnection, buff []byte) {

	Data, err := workerman_go.ParseAndVerifySignJsonTime(buff, conn.Worker().Config.SignKey)
	if err != nil {
		return
	}

	switch Data.Cmd {
	case workerman_go.CommandComponentAuthRequest:
		var RegisterInfo workerman_go.ProtocolRegister
		jsonErr := json.Unmarshal(Data.Json, &RegisterInfo)
		if jsonErr != nil {
			return
		}
		//回复
		RegisterInfo.Authed = "1"
		RegisterInfo.Data = "register say passed"
		SendSignData(RegisterInfo, conn)

		conn.Worker().ConnectionsLock.Lock()

		switch RegisterInfo.ComponentType {
		case workerman_go.ComponentIdentifiersTypeBusiness:
			s._workerConnections[conn.GetClientIdInfo().ClientGatewayNum] = conn
		case workerman_go.ComponentIdentifiersTypeGateway:
			conn.Set(keyGatewayLanInfo, RegisterInfo.ProtocolPublicGatewayConnectionInfo)
			s._workerConnections[conn.GetClientIdInfo().ClientGatewayNum] = conn
		}

		conn.Worker().ConnectionsLock.Unlock()
		s.broadcastOnBusinessConnected(conn)

	}

}

func (s *Server) OnClose(conn workerman_go.InterfaceConnection) {
	conn.Worker().ConnectionsLock.Lock()

	_, isGateway := s._gatewayConnections[conn.GetClientIdInfo().ClientGatewayNum]

	delete(s._gatewayConnections, conn.GetClientIdInfo().ClientGatewayNum)
	delete(s._workerConnections, conn.GetClientIdInfo().ClientGatewayNum)

	conn.Worker().ConnectionsLock.Unlock()

	if isGateway {
		//如果有gateway离线，则广播全部business
		s.broadcastOnBusinessConnected(nil)
	}
}

func (s *Server) broadcastOnBusinessConnected(conn workerman_go.InterfaceConnection) {
	s.ConnectionsLock.RLock()
	defer s.ConnectionsLock.RUnlock()

	var gatewayList []workerman_go.ProtocolPublicGatewayConnectionInfo

	for _, item := range s._gatewayConnections {
		if gatewayLanInfo, ok := item.Get(keyGatewayLanInfo); ok {
			gatewayList = append(gatewayList, gatewayLanInfo.(workerman_go.ProtocolPublicGatewayConnectionInfo))
		}
	}

	//如果是指定发送的
	if conn != nil {
		SendSignData(workerman_go.ProtocolRegisterBroadCastComponentGateway{
			Msg:         "BroadcastOnBusinessConnected",
			Data:        "",
			GatewayList: gatewayList,
		}, conn)
		return
	}

	//广播给business连接gatewaylist列表
	for _, WorkerConn := range s._workerConnections {
		SendSignData(workerman_go.ProtocolRegisterBroadCastComponentGateway{
			Msg:         "BroadcastOnBusinessConnected",
			Data:        "",
			GatewayList: gatewayList,
		}, WorkerConn)
	}

}

func SendSignData(data any, conn workerman_go.InterfaceConnection) {
	timeOut := time.Duration(workerman_go.TimeOutSecond) * time.Second
	var CommandInt int
	switch data.(type) {
	case workerman_go.ProtocolRegister:
		CommandInt = workerman_go.CommandComponentAuthRequest
	case workerman_go.ProtocolRegisterBroadCastComponentGateway:
		CommandInt = workerman_go.CommandComponentGatewayList
	}

	timeByte, err := workerman_go.GenerateSignTimeByte(CommandInt, data, conn.Worker().Config.SignKey, func() time.Duration {
		return timeOut
	})
	if err != nil {
		log.Println(err)
		return
	}

	err = conn.Send(timeByte.ToByte())
	if err != nil {
		log.Println("[send error]:", err)
		return
	}
}
