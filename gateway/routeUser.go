package gateway

import (
	"context"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var forwardTimeout = 60

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

		s.ConnectedBusinessLock.RLock()
		workerArray := make([]string, len(s.ConnectedBusinessMap))
		num := gatewayNum & 0xF
		remainder := int(num) % len(workerArray)
		signDataToBusiness, _ := workerman_go.GenerateSignTimeByte(workerman_go.CommandGatewayForwardUserOnConnect, workerman_go.ProtocolForwardUserOnConnect{
			ClientId: ConnectionUser.GatewayIdInfo.GenerateGatewayClientId(),
		}, s.Config.SignKey, func() time.Duration {
			return time.Duration(forwardTimeout) * time.Second
		})

		s.ConnectedBusinessMap[workerArray[remainder]].Send(signDataToBusiness)
		s.ConnectedBusinessLock.Unlock()

		channelBuff := make(chan []byte)

		// <-ConnectionUser.Ctx.Done() 的时候关闭channel
		defer func() {
			//关闭管道
			close(channelBuff)

			//解除所有的group映射
			gStr, ok := ConnectionUser.Get(constGroups)
			if ok {
				var groupS groupsKv
				json.Unmarshal([]byte(gStr), &groupS)
				s.groupConnectionsLock.Lock()
				//具体的conn名下所有的group，并且删除映射关系
				for _, group := range groupS.Groups {
					delete(s.groupConnections[group], ConnectionUser.TcpWsConnection().GetClientIdInfo().ClientGatewayNum)
				}
				s.groupConnectionsLock.Unlock()

			}

			//解除所有的uid
			uid, ok2 := ConnectionUser.Get(constUid)
			if ok2 {
				s.uidConnectionsLock.Lock()
				delete(s.uidConnections, uid)
				s.uidConnectionsLock.Unlock()
			}

			//异步收到通知， 等待connlistlock锁定用完后，抢占
			s.ConnectionsLock.Lock()
			//删除列表
			delete(s.Connections, ConnectionUser.GetClientIdInfo().ClientGatewayNum)
			s.ConnectionsLock.Unlock()
		}()

		//阻塞式把 reader读取数据到channel
		go userChannelBuff(channelBuff, ConnectionUser)

		for {
			select {
			case msgBuff := <-channelBuff:
				s.ConnectedBusinessLock.RLock()
				workerArray := make([]string, len(s.ConnectedBusinessMap))
				num := gatewayNum & 0xF
				//workerarray数量为7的时候， 可能的值就是0-6
				remainder := int(num) % len(workerArray)
				signDataToBusiness, signDataToBusinessErr := workerman_go.GenerateSignTimeByte(workerman_go.CommandGatewayForwardUserOnMessage, workerman_go.ProtocolForwardUserOnMessage{
					ClientId: ConnectionUser.GatewayIdInfo.GenerateGatewayClientId(),
					Message:  string(msgBuff),
				}, s.Config.SignKey, func() time.Duration {
					return time.Duration(forwardTimeout) * time.Second
				})
				if signDataToBusinessErr != nil {
					continue
				}
				s.ConnectedBusinessMap[workerArray[remainder]].Send(signDataToBusiness)
				s.ConnectedBusinessLock.Unlock()

			case <-ConnectionUser.Ctx.Done():

				//主动cancel()关闭协程，或者 read err触发

				s.ConnectedBusinessLock.RLock()
				workerArray := make([]string, len(s.ConnectedBusinessMap))
				num := gatewayNum & 0xF
				remainder := int(num) % len(workerArray)
				signDataToBusiness, signDataToBusinessErr := workerman_go.GenerateSignTimeByte(workerman_go.CommandGatewayForwardUserOnClose, workerman_go.ProtocolForwardUserOnClose{
					ClientId: ConnectionUser.GatewayIdInfo.GenerateGatewayClientId(),
				}, s.Config.SignKey, func() time.Duration {
					return time.Duration(forwardTimeout) * time.Second
				})
				if signDataToBusinessErr != nil {
					continue
				}
				s.ConnectedBusinessMap[workerArray[remainder]].Send(signDataToBusiness)
				s.ConnectedBusinessLock.Unlock()

				return
				//触发defer 关闭channel
			}

		}
	})
}
