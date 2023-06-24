package business

import (
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"sync"
	"time"
)

type Server struct {
	//已连接的 注册服务
	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]*workerman_go.AsyncTcpWsConnection

	//已连接的Gateway
	ConnectedGatewayLock *sync.RWMutex
	ConnectedGatewayMap  map[string]*workerman_go.AsyncTcpWsConnection
	//配置中心
	Config *workerman_go.ConfigGatewayWorker

	OnConnect func(conn workerman_go.InterfaceConnection)
	OnMessage func(conn workerman_go.InterfaceConnection, buff []byte)
	OnClose   func(conn workerman_go.InterfaceConnection)
}

// 尝试连接gateway
func (s *Server) connectGateway(gatewayInfo *workerman_go.ProtocolRegisterBroadCastComponentGateway) {

	var newGateway []string

	s.ConnectedGatewayLock.Lock()

	//需要连接的新加入集群设备
	for _, gateway := range gatewayInfo.GatewayList {
		//需要连接的新加入集群设备
		connected, ok := s.ConnectedGatewayMap[gateway.GatewayAddr]
		if !ok {
			newGateway = append(newGateway, gateway.GatewayAddr)
			continue
		}

		//未包含 `已连接的Gateway` 即为被踢出集群
		connected.Close()
	}
	s.ConnectedGatewayLock.Unlock()

	//开协程连接gateway了

	for _, gatewayAddress := range newGateway {

		gateway := workerman_go.NewAsyncTcpWsConnection(gatewayAddress)
		gateway.OnConnect = func(connection *workerman_go.AsyncTcpWsConnection) {
			// ### 2. 发送身份认证请求
			s.sendSignData(workerman_go.ProtocolRegister{
				ComponentType:                       0,
				Name:                                "",
				ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
				Data:                                "",
				Authed:                              "",
			}, gateway)
		}

		gateway.OnMessage = func(connection *workerman_go.TcpWsConnection, buff []byte) {
			//business 接受gateway 事件
			Data, parseErr := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
			if parseErr != nil {
				return
			}

			switch Data.Cmd {
			case workerman_go.CommandComponentAuthRequest:
				//请求认证
				s.sendSignData(workerman_go.ProtocolRegister{}, gateway)

			case workerman_go.CommandGatewayForwardUserOnConnect:
				//用户连接
			case workerman_go.CommandGatewayForwardUserOnMessage:
				//用户消息
			case workerman_go.CommandGatewayForwardUserOnClose:
				//用户关闭
			}

		}
	}
}

func (s *Server) sendSignData(data any, conn *workerman_go.AsyncTcpWsConnection) {
	var CMDInt int
	switch data.(type) {
	case workerman_go.ProtocolRegister:
		CMDInt = workerman_go.CommandComponentAuthRequest
		//todo 暂未实现接受用户事件
	}

	SignData, err := workerman_go.GenerateSignTimeByte(CMDInt, data, s.Config.SignKey, func() time.Duration {
		return time.Second * 120
	})
	if err != nil {
		return
	}

	err = conn.Send(SignData.ToByte())
	if err != nil {
		conn.Close()
		return
	}

}

func (s *Server) Run() {

	//todo ### 1. `AsyncWebsocket`连接 `register注册发现`

	urlRegister := fmt.Sprintf("%s%s", s.Config.RegisterPublicHostForComponent, workerman_go.RegisterForComponent)
	register := workerman_go.NewAsyncTcpWsConnection(urlRegister)
	register.OnConnect = func(connection *workerman_go.AsyncTcpWsConnection) {
		s.sendSignData(workerman_go.ProtocolRegister{}, register)
	}

	register.OnMessage = func(connection *workerman_go.TcpWsConnection, buff []byte) {

		Data, parseErr := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
		if parseErr != nil {
			return
		}

		switch Data.Cmd {
		case workerman_go.CommandComponentAuthRequest:
			//发送身份认证请求

			var registerMsg workerman_go.ProtocolRegister
			err := json.Unmarshal(Data.Json, &registerMsg)
			if err != nil {
				return
			}

			//认证通过，加入映射
			if registerMsg.Authed == "1" {
				s.ConnectedRegisterLock.Lock()
				s.ConnectedRegisterMap[register.GetRemoteAddress()] = register
				s.ConnectedRegisterLock.Unlock()
				return
			}
			//要求认证
			s.sendSignData(workerman_go.ProtocolRegister{}, register)
		case workerman_go.CommandComponentGatewayList:
			//todo ### 3. （发生多次）等待 `register注册发现` 返回 `[]Gateway` 列表
			var gatewayInfo workerman_go.ProtocolRegisterBroadCastComponentGateway

			err := json.Unmarshal(Data.Json, &gatewayInfo)
			if err != nil {
				return
			}
			//协程连接每一个gateway,或者Close()不存在的gateway
			s.connectGateway(&gatewayInfo)
		}

	}

}
