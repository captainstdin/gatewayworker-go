package gateway

import (
	"context"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"log"
)

// componentBuffChannel 转换为reader.io为chann
func componentBuffChannel(c chan []byte, connection *workerman_go.TcpWsConnection) {
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

// listenComponent 监听sdk指令
func (s *Server) listenComponent() {

	//SDk或者Business处理器
	s.gin.GET(workerman_go.GatewayForBusinessWsPath, func(ctx *gin.Context) {
		clientConn, err := upgraderWs.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("【business】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}

		s.ConnectedBusinessLock.Lock()
		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())

		ConnectionBusiness := &workerman_go.TcpWsConnection{
			RemoteAddress: ctx.Request.RemoteAddr,
			RequestCtx:    ctx,
			Ctx:           TcpWsCtx,
			CtxF:          TcpWsCancel,
			GatewayIdInfo: &workerman_go.GatewayIdInfo{},
			Name:          "default",
			Address:       "",
			Port:          0,
			FdWs:          clientConn,
			OnConnect:     nil,
			OnMessage:     nil,
			OnClose:       nil,
		}

		s.ConnectedBusinessMap[ctx.Request.RemoteAddr] = ConnectionBusiness
		s.ConnectedBusinessLock.Unlock()

		//onConnected

		channelBuff := make(chan []byte)

		//return的时候关闭channel
		defer func(c chan []byte) {
			close(c)
		}(channelBuff)

		go componentBuffChannel(channelBuff, ConnectionBusiness)

		sdk := gatewayApi{
			Server: s,
			WsConn: ConnectionBusiness,
		}
		for {
			select {
			case buff := <-channelBuff:
				//SDK 或者business 发来的指令

				Command, parseError := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
				if parseError != nil {
					return
				}

				switch Command.Cmd {
				case workerman_go.GatewayCommandSendToAll:
					var cmd workerman_go.GcmdSendToAll
					json.Unmarshal(Command.Json, &cmd)
					sdk.SendToAll(cmd.Data, cmd.ClientIdArray, cmd.ExcludeClientId)
				case workerman_go.GatewayCommandSendToClient:
					var cmd workerman_go.GcmdSendToClient
					json.Unmarshal(Command.Json, &cmd)
					sdk.SendToClient(cmd.ClientId, cmd.SendData)

				case workerman_go.GatewayCommandCloseClient:
					var cmd workerman_go.GcmdCloseClient
					json.Unmarshal(Command.Json, &cmd)
					sdk.CloseClient(cmd.ClientId)

				case workerman_go.GatewayCommandIsOnline:
					var cmd workerman_go.GcmdIsOnline
					json.Unmarshal(Command.Json, &cmd)
					sdk.IsOnline(cmd.ClientId)

				case workerman_go.GatewayCommandBindUid:
					var cmd workerman_go.GcmdBindUid
					json.Unmarshal(Command.Json, &cmd)
					sdk.BindUid(cmd.ClientId, cmd.Uid)

				case workerman_go.GatewayCommandUnbindUid:
					var cmd workerman_go.GcmdUnbindUid
					json.Unmarshal(Command.Json, &cmd)
					sdk.UnbindUid(cmd.ClientId, cmd.Uid)

				case workerman_go.GatewayCommandIsUidOnline:
					var cmd workerman_go.GcmdIsUidOnline
					json.Unmarshal(Command.Json, &cmd)
					sdk.IsUidOnline(cmd.Uid)

				case workerman_go.GatewayCommandGetClientIdByUid:
					var cmd workerman_go.GcmdGetClientIdByUid
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetClientIdByUid(cmd.Uid)

				case workerman_go.GatewayCommandGetUidByClientId:
					var cmd workerman_go.GcmdGetUidByClientId
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetUidByClientId(cmd.ClientId)

				case workerman_go.GatewayCommandSendToUid:
					var cmd workerman_go.GcmdSendToUid
					json.Unmarshal(Command.Json, &cmd)
					sdk.SendToUid(cmd.Uid, cmd.Message)

				case workerman_go.GatewayCommandJoinGroup:
					var cmd workerman_go.GcmdJoinGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.JoinGroup(cmd.ClientId, cmd.Group)

				case workerman_go.GatewayCommandLeaveGroup:
					var cmd workerman_go.GcmdLeaveGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.LeaveGroup(cmd.ClientId, cmd.Group)

				case workerman_go.GatewayCommandUngroup:
					var cmd workerman_go.GcmdUngroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.Ungroup(cmd.Group)

				case workerman_go.GatewayCommandSendToGroup:
					var cmd workerman_go.GcmdSendToGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.SendToGroup(cmd.Group, cmd.Message, cmd.ExcludeClientId)

				case workerman_go.GatewayCommandGetClientIdCountByGroup:
					var cmd workerman_go.GcmdGetClientIdCountByGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetClientIdCountByGroup(cmd.Group)

				case workerman_go.GatewayCommandGetClientSessionsByGroup:
					var cmd workerman_go.GcmdGetClientSessionsByGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetClientSessionsByGroup(cmd.Group)

				case workerman_go.GatewayCommandGetAllClientIdCount:
					var cmd workerman_go.GcmdGetAllClientIdCount
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllClientIdCount()

				case workerman_go.GatewayCommandGetAllClientSessions:
					var cmd workerman_go.GcmdGetAllClientSessions
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllClientSessions()

				case workerman_go.GatewayCommandSetSession:
					var cmd workerman_go.GcmdSetSession
					json.Unmarshal(Command.Json, &cmd)
					sdk.SetSession(cmd.ClientId, cmd.Data)

				case workerman_go.GatewayCommandUpdateSession:
					var cmd workerman_go.GcmdUpdateSession
					json.Unmarshal(Command.Json, &cmd)
					sdk.UpdateSession(cmd.ClientId, cmd.Data)

				case workerman_go.GatewayCommandGetSession:
					var cmd workerman_go.GcmdGetSession
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetSession(cmd.ClientId)

				case workerman_go.GatewayCommandGetClientIdListByGroup:
					var cmd workerman_go.GcmdGetClientIdListByGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetClientIdListByGroup(cmd.Group)
				case workerman_go.GatewayCommandGetAllClientIdList:
					var cmd workerman_go.GcmdGetAllClientIdList
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllClientIdList()
				case workerman_go.GatewayCommandGetUidListByGroup:
					var cmd workerman_go.GcmdGetUidListByGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetUidListByGroup(cmd.Group)

				case workerman_go.GatewayCommandGetUidCountByGroup:
					var cmd workerman_go.GcmdGetUidCountByGroup
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetUidCountByGroup(cmd.Group)

				case workerman_go.GatewayCommandGetAllUidList:
					var cmd workerman_go.GcmdGetAllUidList
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllUidList()

				case workerman_go.GatewayCommandGetAllUidCount:
					var cmd workerman_go.GcmdGetAllUidCount
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllUidCount()

				case workerman_go.GatewayCommandGetAllGroupIdList:
					var cmd workerman_go.GcmdGetAllGroupIdList
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllGroupIdList()

				case workerman_go.GatewayCommandGetAllGroupCount:
					var cmd workerman_go.GcmdGetAllGroupCount
					json.Unmarshal(Command.Json, &cmd)
					sdk.GetAllGroupCount()
				}

			case <-ConnectionBusiness.Ctx.Done():
				//主动关闭协程，或者errBuff触发
				return
			}

		}

	})
}
