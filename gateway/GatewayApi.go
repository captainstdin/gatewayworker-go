package gateway

import "gatewaywork-go/workerman_go"

type gatewayApi struct {
	Server *Server
	WsConn *workerman_go.TcpWsConnection
}

func (g *gatewayApi) SendToAll(data []byte, client_id_array []string, exclude_client_id []string) {

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	set_client_id_array := make(map[string]struct{}) // 创建一个空的哈希集合
	for _, item := range client_id_array {
		set_client_id_array[item] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
	}

	set_exclude_client_id := make(map[string]struct{}) // 创建一个空的哈希集合
	for _, item := range exclude_client_id {
		set_exclude_client_id[item] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
	}

	for _, conn := range g.Server.Connections {
		//白名单有的话，就不执行排除了，
		if _, found := set_client_id_array[conn.GetClientId()]; found {
			conn.Send(data)
			continue
		}

		//如果不在黑名单，就发送
		if _, found := set_exclude_client_id[conn.GetClientId()]; found == false {
			conn.Send(data)
		}

	}
}

func (g *gatewayApi) SendToClient(client_id string, send_data []byte) {

	c, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[c.ClientGatewayNum]; ok {
		conn.Send(send_data)
	}

}

func (g *gatewayApi) CloseClient(client_id string) {

	c, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[c.ClientGatewayNum]; ok {
		conn.Close()
	}
}

func (g *gatewayApi) IsOnline(client_id string) int {
	c, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return 0
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if _, ok := g.Server.Connections[c.ClientGatewayNum]; ok {
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
