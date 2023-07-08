package gateway

import (
	"encoding/json"
	"gatewaywork-go/workerman_go"
)

func (s *Server) connectBusness() {

	register := workerman_go.NewAsyncTcpWsConnection(s.Config.RegisterPublicHostForComponent)

	register.OnConnect = func(connection *workerman_go.TcpWsConnection) {
		s.ConnectedRegisterLock.Lock()
		s.ConnectedRegisterMap[connection.GetRemoteAddress()] = connection
		s.ConnectedRegisterLock.Unlock()
	}

	register.OnMessage = func(connection *workerman_go.TcpWsConnection, buff []byte) {
		//business的指令
		Data, parse := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
		if parse != nil {
			return
		}
		switch Data.Cmd {
		//身份认证
		case workerman_go.CommandComponentAuthRequest:
			var reg workerman_go.ProtocolRegister
			err := json.Unmarshal(buff, &reg)
			if err != nil {
				return
			}
			if reg.Authed == "1" {
				//认证成功，把register加入内存列表
				s.ConnectedRegisterLock.RLock()
				s.ConnectedRegisterMap[connection.GetRemoteAddress()] = connection
				s.ConnectedRegisterLock.Unlock()
				return
			}
			//发送身份注册认证
			s.sendSignData(workerman_go.ProtocolRegister{}, connection)
		}
	}

	register.OnClose = func(connection *workerman_go.TcpWsConnection) {
		//todo 重连business
	}

	//启动连接 gateway连接register
	err := register.Connect()
	if err != nil {
		register.Connect()
	}
}
