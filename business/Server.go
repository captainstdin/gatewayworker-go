package business

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gatewaywork-go/workerman_go"
	"sync"
	"time"
)

type Server struct {
	//已连接的 注册服务
	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]*workerman_go.TcpWsConnection //key是remoteAddress

	//已连接的Gateway
	ConnectedGatewayLock *sync.RWMutex
	ConnectedGatewayMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress
	//配置中心
	Config *workerman_go.ConfigGatewayWorker

	OnConnect func(conn workerman_go.InterfaceConnection)
	OnMessage func(conn workerman_go.InterfaceConnection, buff []byte)
	OnClose   func(conn workerman_go.InterfaceConnection)

	Ctx       context.Context
	CtxCancel context.CancelFunc
}

// 尝试连接gateway，并且设置好监听gateway,并且加入ConnectedGatewayMap[]
func (s *Server) connectGateway(gatewayInfo *workerman_go.ProtocolRegisterBroadCastComponentGateway) {

	var newGateway []string

	//需要连接的新加入集群设备
	for _, gateway := range gatewayInfo.GatewayList {
		s.ConnectedGatewayLock.Lock()
		//需要连接的新加入集群设备
		connected, ok := s.ConnectedGatewayMap[gateway.GatewayAddr]
		s.ConnectedGatewayLock.Unlock() //这里不解锁，会死锁，因为Close会触发Onclose()事件，从而造成竞争资源阻塞
		if !ok {
			newGateway = append(newGateway, gateway.GatewayAddr)
			continue
		}
		//未包含 `已连接的Gateway` 即为被踢出集群
		connected.Close()
	}

	//开协程连接gateway了

	for _, gatewayAddress := range newGateway {

		gateway := workerman_go.NewAsyncTcpWsConnection(gatewayAddress)
		gateway.OnConnect = func(connection *workerman_go.TcpWsConnection) {
			// ### 2. 发送身份认证请求
			s.sendToAsyncData(workerman_go.ProtocolRegister{
				ComponentType:                       0,
				Name:                                "",
				ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
				Data:                                "",
				Authed:                              "",
			}, gateway)
		}

		gateway.OnClose = func(connection *workerman_go.TcpWsConnection) {
			//被断开的时候

			s.ConnectedGatewayLock.Lock()
			defer s.ConnectedGatewayLock.Unlock()
			delete(s.ConnectedGatewayMap, connection.GetRemoteAddress())
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
				var reg workerman_go.ProtocolRegister
				err := json.Unmarshal(buff, &reg)
				if err != nil {
					return
				}
				if reg.Authed == "1" {
					s.ConnectedGatewayLock.RLock()
					s.ConnectedGatewayMap[connection.GetRemoteAddress()] = connection
					s.ConnectedGatewayLock.Unlock()
					return
				}
				s.sendToAsyncData(workerman_go.ProtocolRegister{}, gateway)
			case workerman_go.CommandGatewayForwardUserOnConnect:
				//用户连接
				var OnConnect workerman_go.ProtocolForwardUserOnConnect
				err := json.Unmarshal(buff, &OnConnect)
				if err != nil {
					return
				}
				//todo

			case workerman_go.CommandGatewayForwardUserOnMessage:
				//用户消息
				var OnMessage workerman_go.ProtocolForwardUserOnMessage
				err := json.Unmarshal(buff, &OnMessage)
				if err != nil {
					return
				}
				//todo

			case workerman_go.CommandGatewayForwardUserOnClose:
				//用户关闭
				var OnClose workerman_go.ProtocolForwardUserOnClose
				err := json.Unmarshal(buff, &OnClose)
				if err != nil {
					return
				}
				//todo
			}

		}

		go gateway.Connect()
	}
}

// 同步发送SDK 到gateway
func (s *Server) sendToAsyncData(data any, conn workerman_go.InterfaceConnection) {
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

func (s *Server) Run() error {

	//todo ### 1. `AsyncWebsocket`连接 `register注册发现`

	urlRegister := fmt.Sprintf("%s%s", s.Config.RegisterPublicHostForComponent, workerman_go.RegisterForComponent)
	register := workerman_go.NewAsyncTcpWsConnection(urlRegister)
	register.OnConnect = func(connection *workerman_go.TcpWsConnection) {
		s.sendToAsyncData(workerman_go.ProtocolRegister{}, register)
	}

	register.OnClose = func(connection *workerman_go.TcpWsConnection) {
		//todo 断网 循环连接
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
				s.ConnectedRegisterMap[register.GetRemoteAddress()] = connection
				s.ConnectedRegisterLock.Unlock()
				return
			}
			//要求认证
			s.sendToAsyncData(workerman_go.ProtocolRegister{}, register)
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

	<-s.Ctx.Done()

	return errors.New("business exit(0)")

}

func NewServer(name string, conf *workerman_go.ConfigGatewayWorker) *Server {

	ctx, cfunc := context.WithCancel(context.Background())
	s := &Server{
		Ctx:                   ctx,
		CtxCancel:             cfunc,
		ConnectedRegisterLock: nil,
		ConnectedRegisterMap:  nil,
		ConnectedGatewayLock:  nil,
		ConnectedGatewayMap:   nil,
		Config:                nil,
		OnConnect:             nil,
		OnMessage:             nil,
		OnClose:               nil,
	}
	return s
}
