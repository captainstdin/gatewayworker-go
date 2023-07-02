package gateway

import (
	"gatewaywork-go/workerman_go"
	"strconv"
)

//warning 需要严重提醒，任何 发来过来的协议中的ClientID不是真正的 Hex字符串，而是 uint64 Num表现形式
// gatewayApi 所有的 ClientID均为 uint64 ClientToken.ClientGatewayNum的uint64表现形式，例如 +00000000000000001

type gatewayApi struct {
	Server *Server
	WsConn *workerman_go.TcpWsConnection
}

func (g *gatewayApi) SendToAll(data []byte, client_id_array []string, exclude_client_id []string) {

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	setClientIdArray := make(map[uint64]struct{}) // 创建一个空的哈希集合
	for _, item := range client_id_array {
		parseUint, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			continue
		}
		setClientIdArray[parseUint] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
	}

	setExcludeClientId := make(map[uint64]struct{}) // 创建一个空的哈希集合
	for _, item := range exclude_client_id {
		parseUint, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			continue
		}
		setExcludeClientId[parseUint] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
	}

	for keyuint64, conn := range g.Server.Connections {
		//白名单有的话，就不执行排除了，
		if _, found := setClientIdArray[keyuint64]; found {
			conn.Send(data)
			continue
		}

		//如果不在黑名单，就发送
		if _, found := setExcludeClientId[keyuint64]; found == false {
			conn.Send(data)
		}

	}
}

func (g *gatewayApi) SendToClient(client_id string, send_data []byte) {
	parseUint, err := strconv.ParseUint(client_id, 10, 64)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[parseUint]; ok {
		conn.Send(send_data)
	}

}

func (g *gatewayApi) CloseClient(client_id string) {

	parseUint, err := strconv.ParseUint(client_id, 10, 64)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[parseUint]; ok {
		conn.Close()
	}
}

func (g *gatewayApi) IsOnline(client_id string) int {
	parseUint, err := strconv.ParseUint(client_id, 10, 64)
	if err != nil {
		return 0
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if _, ok := g.Server.Connections[parseUint]; ok {
		return 1
	}
	return 0
}

func (g *gatewayApi) BindUid(client_id string, uid string) {

	//todo
}

func (g *gatewayApi) UnbindUid(client_id string, uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) IsUidOnline(uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetClientIdByUid(uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetUidByClientId(client_id string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) SendToUid(uid string, message string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) JoinGroup(client_id string, group string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) LeaveGroup(client_id string, group string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) Ungroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) SendToGroup(group string, message string, exclude_client_id []string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetClientIdCountByGroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetClientSessionsByGroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllClientIdCount() int {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllClientSessions() map[string]interface{} {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) SetSession(client_id string, data workerman_go.SessionKv) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) UpdateSession(client_id string, data workerman_go.SessionKv) {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetSession(client_id string) workerman_go.SessionKv {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetClientIdListByGroup(group string) []string {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllClientIdList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetUidListByGroup(group string) []string {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetUidCountByGroup(group string) int {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllUidList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllUidCount() int {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllGroupIdList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *gatewayApi) GetAllGroupCount() int {
	//TODO implement me
	panic("implement me")
}
