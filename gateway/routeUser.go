package gateway

import (
	"context"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"log"
)

func userChannelBuff(c chan []byte, connection *workerman_go.TcpWsConnection) {
	for {
		_, buff, errBuff := connection.FdWs.ReadMessage()
		if errBuff != nil {
			//如果发现已经断开，就通知协程结束
			connection.TcpWsConnection().CtxF()
			//todo  forward 转发给Business  onclose()
			return
		}
		c <- buff
	}
}

func (s *Server) listenUser() {

	//普通用户
	s.gin.GET(workerman_go.GatewayForUserPath, func(ctx *gin.Context) {
		clientConn, err := upgraderWs.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("【clientUser】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}

		s.ConnectionsLock.Lock()
		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())

		ConnectionUser := &workerman_go.TcpWsConnection{
			RemoteAddress: ctx.Request.RemoteAddr,
			RequestCtx:    ctx,
			Ctx:           TcpWsCtx,
			CtxF:          TcpWsCancel,
			ClientToken:   &workerman_go.ClientToken{},
			Name:          "defaultUser",
			Address:       "",
			Port:          0,
			FdWs:          clientConn,
			OnConnect:     nil,
			OnMessage:     nil,
			OnClose:       nil,
		}

		s.Connections[genPrimaryKeyUint64(s.Connections)] = ConnectionUser
		s.ConnectionsLock.Unlock()
		//todo 发送OnConnect 给business

		channelBuff := make(chan []byte)

		//return的时候关闭channel
		defer func(c chan []byte) {
			close(c)
		}(channelBuff)

		go userChannelBuff(channelBuff, ConnectionUser)

		for {
			select {
			case <-channelBuff:
				//todo forward 转发给固定的 Business
			case <-ConnectionUser.Ctx.Done():
				//主动关闭协程，或者errBuff触发
				return
			}

		}
	})
}
