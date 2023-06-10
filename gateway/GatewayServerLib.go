package gateway

import "gatewaywork-go/workerman_go"

func (g *GatewayServer) SendToAll(data []byte, client_id_array []string, exclude_client_id []string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) SendToClient(client_id string, send_data []byte) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) CloseClient(client_id string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) IsOnline(client_id string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) BindUid(client_id string, uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) UnbindUid(client_id string, uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) IsUidOnline(uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetClientIdByUid(uid string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetUidByClientId(client_id string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) SendToUid(uid string, message string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) JoinGroup(client_id string, group string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) LeaveGroup(client_id string, group string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) Ungroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) SendToGroup(group string, message string, exclude_client_id []string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetClientIdCountByGroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetClientSessionsByGroup(group string) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllClientIdCount() int {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllClientSessions() map[string]interface{} {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) SetSession(client_id string, data workerman_go.SessionKv) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) UpdateSession(client_id string, data workerman_go.SessionKv) {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetSession(client_id string) workerman_go.SessionKv {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetClientIdListByGroup(group string) []string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllClientIdList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetUidListByGroup(group string) []string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetUidCountByGroup(group string) int {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllUidList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllUidCount() int {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllGroupIdList() []string {
	//TODO implement me
	panic("implement me")
}

func (g *GatewayServer) GetAllGroupCount() int {
	//TODO implement me
	panic("implement me")
}
