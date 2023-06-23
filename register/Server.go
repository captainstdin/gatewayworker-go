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
	keyAuth          = "Auth"
	keyComponentType = "ComponentType"

	keyGatewayLanInfo = "GatewayLanInfo"
)

type Server struct {
	workerman_go.Worker
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
		OnWorkerStart:   onWorkerStart,
		OnConnect:       onConnect,
		OnMessage:       onMessage,
		OnClose:         nil,
		Ctx:             ctx,
		CtxF:            cf,
		Config:          conf,
	}

	return &Server{w}
}

func onWorkerStart(worker *workerman_go.Worker) {
	startInfo := bytes.Buffer{}
	startInfo.WriteByte('[')
	startInfo.WriteString(worker.Name)
	startInfo.WriteString("] Starting  server at  ->【")
	startInfo.WriteString(worker.ListenAddress)
	startInfo.WriteString(worker.ListenPath)
	startInfo.WriteString("】 Listening...")
	log.Println(strconv.Quote(startInfo.String()))
}

func onConnect(conn workerman_go.InterfaceConnection) {
	//非阻塞
	log.Printf("[worker-business]new component connected! %s", conn.GetRemoteAddress())
	SendSignData(workerman_go.ProtocolRegister{
		ComponentType:                       0,
		Name:                                "",
		ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
		Data:                                "Business.request.auth",
		Authed:                              "0",
	}, conn)

}

func onMessage(conn workerman_go.InterfaceConnection, buff []byte) {

	Data, err := workerman_go.ParseAndVerifySignJsonTime(buff, conn.Worker().Config.SignKey)
	if err != nil {
		return
	}

	switch Data.Cmd {
	case workerman_go.CommandComponentAuthRequest:
		var RegisterInfo workerman_go.ProtocolRegister
		json.Unmarshal(Data.Json, &RegisterInfo)
		//回复
		RegisterInfo.Authed = "1"
		RegisterInfo.Data = "register say passed"
		SendSignData(RegisterInfo, conn)

		conn.Worker().ConnectionsLock.Lock()
		//设置已验证
		conn.Set(keyAuth, true)
		//设置类型
		conn.Set(keyComponentType, RegisterInfo.ComponentType)

		if RegisterInfo.ComponentType == workerman_go.ComponentIdentifiersTypeGateway {
			conn.Set(keyGatewayLanInfo, RegisterInfo.ProtocolPublicGatewayConnectionInfo)
		}

		conn.Worker().ConnectionsLock.Unlock()
		BroadcastOnBusinessConnected(conn, &RegisterInfo)

	}

}

func BroadcastOnBusinessConnected(conn workerman_go.InterfaceConnection, registerInfo *workerman_go.ProtocolRegister) {
	conn.Worker().ConnectionsLock.RLock()
	defer conn.Worker().ConnectionsLock.RUnlock()

	var gatewayList []workerman_go.ProtocolPublicGatewayConnectionInfo

	var businessList []workerman_go.InterfaceConnection

	for _, item := range conn.Worker().Connections {
		//过滤通过认证的conns
		if v, ok := item.Get(keyAuth); ok && v.(bool) {
			//过滤设置了组件类型的conn
			if componentType, ok2 := item.Get(keyComponentType); ok2 {
				//判断组件类型
				switch componentType.(int) {
				case workerman_go.ComponentIdentifiersTypeBusiness:
					businessList = append(businessList, item)
				case workerman_go.ComponentIdentifiersTypeGateway:
					if gatewayLanInfo, ok3 := item.Get(keyGatewayLanInfo); ok3 {
						gatewayList = append(gatewayList, gatewayLanInfo.(workerman_go.ProtocolPublicGatewayConnectionInfo))
					}

				}
			}

		}
	}

	//广播给business连接gatewaylist列表
	for _, item := range businessList {
		SendSignData(workerman_go.ProtocolRegisterBroadCastComponentGateway{
			Msg:         "BroadcastOnBusinessConnected",
			Data:        "",
			GatewayList: gatewayList,
		}, item)
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

	conn.Send(timeByte.ToByte())
}
