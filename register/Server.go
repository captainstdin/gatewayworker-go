package register

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	keyGatewayLanInfo = "GatewayLanInfo"
)

type Server struct {
	workerman_go.Worker

	//business
	_workerConnections map[uint64]workerman_go.InterfaceConnection

	//gateway
	_gatewayConnections map[uint64]workerman_go.InterfaceConnection
}

func NewServer(name string, conf *workerman_go.ConfigGatewayWorker) *Server {

	ctx, cf := context.WithCancel(context.Background())
	w := workerman_go.Worker{
		Connections:      map[uint64]workerman_go.InterfaceConnection{},
		ConnectionsLock:  &sync.RWMutex{},
		ListenAddress:    conf.RegisterListenAddr,
		ListenPath:       workerman_go.RegisterForComponent,
		Name:             name,
		Tls:              false,
		TlsPem:           "",
		TlsKey:           "",
		OnWorkerStart:    nil,
		OnConnect:        nil,
		OnMessage:        nil,
		OnClose:          nil,
		Ctx:              ctx,
		CtxF:             cf,
		Config:           conf,
		ExtraHttpHandles: make(map[string]func(ctx *gin.Context)),
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

	server.Worker.ExtraHttpHandles[workerman_go.RegisterForComponent] = func(ctx *gin.Context) {

		buff, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return
		}
		Command, parseError := workerman_go.ParseAndVerifySignJsonTime(buff, conf.SignKey)

		if parseError != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errCode": http.StatusBadRequest,
				"errMsg":  parseError.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"errCode": http.StatusOK,
			"data":    handleHttpCmd(Command, server),
		})

	}

	return server
}

// handleHttpCmd 处理http SDK
func handleHttpCmd(Command *workerman_go.GenerateComponentSign, server *Server) any {

	fmt.Sprintf("%+v", Command.Cmd)
	switch Command.Cmd {
	case workerman_go.CommandComponentAuthRequest:
		var cmd workerman_go.ProtocolRegister
		json.Unmarshal(Command.Json, &cmd)

		server.ConnectionsLock.RLock()
		defer server.ConnectionsLock.RUnlock()

		var gatewayList []workerman_go.ProtocolPublicGatewayConnectionInfo

		gatewayList = append(gatewayList, workerman_go.ProtocolPublicGatewayConnectionInfo{
			GatewayAddr: "192.168.2.2",
			GatewayPort: "",
		})
		//便利每一个gatewayconn，吧他们的公网连接信息 整理出来
		for _, item := range server._gatewayConnections {
			if gatewayLanInfo, ok := item.Get(keyGatewayLanInfo); ok {
				var gatewayAddress workerman_go.ProtocolPublicGatewayConnectionInfo
				err := json.Unmarshal([]byte(gatewayLanInfo), &gatewayAddress)
				//如果有错误就跳过
				if err != nil {
					continue
				}
				gatewayList = append(gatewayList, gatewayAddress)
			}
		}
		return gatewayList
	case workerman_go.GatewayCommandSendToClient:
		var cmd workerman_go.GcmdSendToClient
		json.Unmarshal(Command.Json, &cmd)

		return nil
	}

	return nil
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
		defer conn.Worker().ConnectionsLock.Unlock()

		switch RegisterInfo.ComponentType {
		case workerman_go.ComponentIdentifiersTypeBusiness:
			s._workerConnections[conn.GetClientIdInfo().ClientGatewayNum] = conn
		case workerman_go.ComponentIdentifiersTypeGateway:
			marshal, marshalErr := json.Marshal(RegisterInfo.ProtocolPublicGatewayConnectionInfo)
			if marshalErr == nil {
				return
			}
			conn.Set(keyGatewayLanInfo, string(marshal))
			s._workerConnections[conn.GetClientIdInfo().ClientGatewayNum] = conn
		}
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

	//便利每一个gatewayconn，吧他们的公网连接信息 整理出来
	for _, item := range s._gatewayConnections {
		if gatewayLanInfo, ok := item.Get(keyGatewayLanInfo); ok {

			var gatewayAddress workerman_go.ProtocolPublicGatewayConnectionInfo

			err := json.Unmarshal([]byte(gatewayLanInfo), &gatewayAddress)
			//如果有错误就跳过
			if err != nil {
				continue
			}
			gatewayList = append(gatewayList, gatewayAddress)
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
