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
		//return 后关闭close再次
		defer clientConn.Close()

		s.ConnectionsLock.Lock()
		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())

		gatewayNum := genPrimaryKeyUint64(s.Connections)

		//todo 填写gatewayclientHex中的公网或者可连接的地址，

		ConnectionUser := &workerman_go.TcpWsConnection{

			RemoteAddress: ctx.Request.RemoteAddr,
			RequestCtx:    ctx,
			Ctx:           TcpWsCtx,
			CtxF:          TcpWsCancel,
			GatewayIdInfo: &workerman_go.GatewayIdInfo{
				ClientGatewayAddr: s.Config.GatewayPublicHostForClient,
				ClientGatewayNum:  gatewayNum,
			},
			Name:      "defaultUser",
			Address:   "",
			Port:      0,
			FdWs:      clientConn,
			OnConnect: nil,
			OnMessage: nil,
			OnClose:   nil,
		}

		s.Connections[gatewayNum] = ConnectionUser
		s.ConnectionsLock.Unlock()
		//todo 发送OnConnect 给business

		channelBuff := make(chan []byte)

		// <-ConnectionUser.Ctx.Done() 的时候关闭channel
		defer func(c chan []byte) {
			close(c)
		}(channelBuff)

		//阻塞式把 reader读取数据到channel
		go userChannelBuff(channelBuff, ConnectionUser)

		for {
			select {
			case <-channelBuff:
				//todo forward 转发给固定的 Business
			case <-ConnectionUser.Ctx.Done():

				//异步收到通知， 等待connlistlock锁定用完后，抢占
				s.ConnectionsLock.Lock()
				//删除列表
				delete(s.Connections, ConnectionUser.TcpWsConnection().GatewayIdInfo.ClientGatewayNum)
				s.ConnectionsLock.Unlock()

				//主动cancel()关闭协程，或者 read err触发
				return
				//触发defer 关闭channel
			}

		}
	})
}
