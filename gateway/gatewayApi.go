package gateway

import (
	"encoding/json"
	"gatewaywork-go/workerman_go"
)

//warning 需要严重提醒，任何 发来过来的协议中的ClientID不是真正的 Hex字符串，而是 uint64 Num表现形式
//这里的调用者都是 SDK或者routeUser.go
// gatewayApi 所有的 ClientID均为 base64的成品ID (由于不在使用hexString，而是Base64 string ClientId)
// gatewayApi 所有的返回的ClientID均为 base64的成品ID
// gatewayApi  distributed-api 标记的注释，是一个分布式接口，SDK或者Business应该遍历gateway服务器调用或者接收

type gatewayApi struct {
	Server *Server
	WsConn *workerman_go.TcpWsConnection
}

const (
	constUid    = "uid"
	constGroups = "groups"
)

// SendToAll distributed-api
func (g *gatewayApi) SendToAll(data []byte, client_id_array []string, exclude_client_id []string) {

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	setClientIdArray := make(map[uint64]struct{}) // 创建一个空的哈希集合
	for _, item := range client_id_array {
		id, err := workerman_go.ParseGatewayClientId(item)
		if err != nil {
			continue
		}
		setClientIdArray[id.ClientGatewayNum] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
	}

	setExcludeClientId := make(map[uint64]struct{}) // 创建一个空的哈希集合
	for _, item := range exclude_client_id {
		id, err := workerman_go.ParseGatewayClientId(item)
		if err != nil {
			continue
		}
		setExcludeClientId[id.ClientGatewayNum] = struct{}{} // 将列表中的元素作为键存储到哈希集合中
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

// SendToClient 发送给指定client，发送期间，ConnectionsLock只读
func (g *gatewayApi) SendToClient(client_id string, send_data string) {
	id, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[id.ClientGatewayNum]; ok {
		conn.Send(send_data)
	}

}

// CloseClient 关闭client，
func (g *gatewayApi) CloseClient(client_id string) {

	id, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if conn, ok := g.Server.Connections[id.ClientGatewayNum]; ok {
		conn.Close()
	}
	return

	// conn.close ->触发for{  锁connections -> delete -> 释放锁connections ;return }
}

func (g *gatewayApi) IsOnline(client_id string) int {
	id, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return 0
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	if _, ok := g.Server.Connections[id.ClientGatewayNum]; ok {
		return 1
	}
	return 0
}

func (g *gatewayApi) BindUid(client_id string, uid string) {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return
	}

	// 把conn绑定到uidConnnections[]上面去
	// 格式{"uid1":{"uint64x1":conn,"uint64x2":conn,},"uid2":....}

	g.Server.uidConnectionsLock.Lock()
	g.Server.uidConnections[uid][clientInfo.ClientGatewayNum] = conn
	g.Server.uidConnectionsLock.Unlock()

	//设置conn的data ，内置conn 的data锁
	conn.Set(constUid, uid)

}

func (g *gatewayApi) UnbindUid(client_id string, uid string) {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.uidConnectionsLock.Lock()
	delete(g.Server.uidConnections[uid], clientInfo.ClientGatewayNum)
	g.Server.uidConnectionsLock.Unlock()
}

// IsUidOnline distributed-api
func (g *gatewayApi) IsUidOnline(uid string) int {
	g.Server.uidConnectionsLock.RLock()
	defer g.Server.uidConnectionsLock.RUnlock()

	if mapCids, ok := g.Server.uidConnections[uid]; ok {
		if len(mapCids) > 0 {
			return 1
		}
		return 0
	}
	return 0
}

// GetClientIdByUid  distributed-api
func (g *gatewayApi) GetClientIdByUid(uid string) []string {
	var listCid []string
	g.Server.uidConnectionsLock.RLock()
	defer g.Server.uidConnectionsLock.RUnlock()
	if mapCids, ok := g.Server.uidConnections[uid]; ok {
		if len(mapCids) > 0 {
			for _, conn := range mapCids {
				listCid = append(listCid, conn.GetClientId())
			}
			return listCid
		}
		return nil
	}
	return nil

}

func (g *gatewayApi) GetUidByClientId(client_id string) string {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return ""
	}

	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放uidConnectionsLock锁
		return ""
	}

	v, okV := conn.Get(constUid)
	if okV {
		return v
	}
	return ""
}

// SendToUid distributed-api
func (g *gatewayApi) SendToUid(uid string, message string) {
	clientIdList := g.GetClientIdByUid(uid)
	for _, clientId := range clientIdList {
		g.SendToClient(clientId, message)
	}
}

type groupsKv struct {
	Groups []string `json:"groups"`
}

func (g *gatewayApi) JoinGroup(client_id string, group string) {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.Lock()
	defer g.Server.ConnectionsLock.Unlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return
	}

	g.Server.groupConnectionsLock.Lock()
	defer g.Server.groupConnectionsLock.Unlock()

	v, ok := conn.TcpWsConnection().Data[constGroups]
	if !ok {
		//释放groupConnectionsLock锁
		//释放ConnectionsLock锁
		return
	}

	var oldGroupKv groupsKv
	err = json.Unmarshal([]byte(v), &oldGroupKv)

	if err != nil {
		//释放groupConnectionsLock锁
		//释放ConnectionsLock锁
		return
	}

	for _, groupOld := range oldGroupKv.Groups {
		if groupOld == groupOld {
			//释放groupConnectionsLock锁
			//释放ConnectionsLock锁
			return
		}
	}

	oldGroupKv.Groups = append(oldGroupKv.Groups, group)

	marshal, errJson := json.Marshal(oldGroupKv)
	if errJson != nil {
		//释放groupConnectionsLock锁
		//释放ConnectionsLock锁
		return
	}

	conn.TcpWsConnection().DataLock.Lock()
	defer conn.TcpWsConnection().DataLock.Unlock()

	//防止死锁，不得调用.set() .get()
	conn.TcpWsConnection().Data[constGroups] = string(marshal)
}

func (g *gatewayApi) LeaveGroup(client_id string, group string) {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.Lock()
	defer g.Server.ConnectionsLock.Unlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return
	}

	joined, ok := conn.TcpWsConnection().Data[constGroups]

	if !ok {
		return
	}
	g.Server.groupConnectionsLock.Lock()
	defer g.Server.groupConnectionsLock.Unlock()

	var oldJoined groupsKv
	err = json.Unmarshal([]byte(joined), &oldJoined)
	if err != nil {
		return
	}

	var newJoined []string

	for _, item := range oldJoined.Groups {
		if item == group {
			continue
		}
		newJoined = append(newJoined, item)
	}

	newJoinedStr, err2 := json.Marshal(groupsKv{Groups: newJoined})

	if err2 != nil {
		return
	}

	conn.TcpWsConnection().DataLock.Lock()
	defer conn.TcpWsConnection().DataLock.Unlock()

	//重塑 conn.groups[]
	conn.TcpWsConnection().Data[constGroups] = string(newJoinedStr)

	delete(g.Server.groupConnections[group], clientInfo.ClientGatewayNum)
}

// Ungroup distributed-api
func (g *gatewayApi) Ungroup(group string) {
	g.Server.groupConnectionsLock.Lock()
	defer g.Server.groupConnectionsLock.Unlock()
	delete(g.Server.groupConnections, group)
}

// SendToGroup distributed-api
func (g *gatewayApi) SendToGroup(group string, message string, exclude_client_id []string) {

	groupMap, ok := g.Server.groupConnections[group]

	if !ok {
		return
	}

	excludeClientId := make(map[uint64]struct{})

	for _, clientId := range exclude_client_id {
		clientInfo, err := workerman_go.ParseGatewayClientId(clientId)
		if err != nil {
			continue
		}
		excludeClientId[clientInfo.ClientGatewayNum] = struct{}{}
	}

	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()

	for _, conn := range groupMap {
		_, exist := excludeClientId[conn.GetClientIdInfo().ClientGatewayNum]
		if exist {
			continue
		}
		conn.Send(message)
	}

}

// GetClientIdCountByGroup distributed-api
func (g *gatewayApi) GetClientIdCountByGroup(group string) int {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()

	groupMap, ok := g.Server.groupConnections[group]

	if !ok {
		return 0
	}
	return len(groupMap)
}

// GetClientSessionsByGroup distributed-api
func (g *gatewayApi) GetClientSessionsByGroup(group string) map[string]workerman_go.SessionKv {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()

	result := make(map[string]workerman_go.SessionKv)

	groupMap, ok := g.Server.groupConnections[group]

	if !ok {
		return nil
	}

	for _, conn := range groupMap {
		result[conn.GetClientId()] = conn.TcpWsConnection().Data
	}

	return result
}

// GetAllClientIdCount    distributed-api
func (g *gatewayApi) GetAllClientIdCount() int {
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()
	return len(g.Server.Connections)
}

// GetAllClientSessions  clientID=>array(...) distributed-api   distributed-api
func (g *gatewayApi) GetAllClientSessions() map[string]workerman_go.SessionKv {
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()
	ClientSessionKv := make(map[string]workerman_go.SessionKv)
	for _, conn := range g.Server.Connections {
		ClientSessionKv[conn.GetClientId()] = conn.TcpWsConnection().Data
	}
	return ClientSessionKv
}

func (g *gatewayApi) SetSession(client_id string, data workerman_go.SessionKv) {

	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.Lock()
	defer g.Server.ConnectionsLock.Unlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return
	}

	//清空
	conn.TcpWsConnection().Data = make(map[string]string)
	for k, v := range data {
		conn.Set(k, v)
	}

}

func (g *gatewayApi) UpdateSession(client_id string, data workerman_go.SessionKv) {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return
	}
	g.Server.ConnectionsLock.Lock()
	defer g.Server.ConnectionsLock.Unlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return
	}

	for k, v := range data {
		conn.Set(k, v)
	}
}

func (g *gatewayApi) GetSession(client_id string) workerman_go.SessionKv {
	clientInfo, err := workerman_go.ParseGatewayClientId(client_id)
	if err != nil {
		return nil
	}
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()

	conn, ok := g.Server.Connections[clientInfo.ClientGatewayNum]
	if !ok {
		//找不到client_id，释放ConnectionsLock锁
		return nil
	}
	return conn.TcpWsConnection().Data
}

// GetClientIdListByGroup distributed-api
func (g *gatewayApi) GetClientIdListByGroup(group string) []string {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()
	var clientIdList []string
	groupMap, ok := g.Server.groupConnections[group]
	if !ok {
		return nil
	}

	for _, conn := range groupMap {
		clientIdList = append(clientIdList, conn.GetClientId())
	}
	return clientIdList
}

// GetAllClientIdList  distributed-api
func (g *gatewayApi) GetAllClientIdList() []string {
	g.Server.ConnectionsLock.RLock()
	defer g.Server.ConnectionsLock.RUnlock()
	var clientIdList []string
	for _, conn := range g.Server.Connections {
		clientIdList = append(clientIdList, conn.GetClientId())
	}
	return clientIdList
}

// GetUidListByGroup distributed-api
func (g *gatewayApi) GetUidListByGroup(group string) []string {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()
	var uidList []string
	groupMap, ok := g.Server.groupConnections[group]
	if !ok {
		return nil
	}

	for _, conn := range groupMap {
		uid, ok2 := conn.TcpWsConnection().Get(constUid)
		if !ok2 {
			continue
		}
		uidList = append(uidList, uid)
	}
	return uidList
}

// GetUidCountByGroup distributed-api
func (g *gatewayApi) GetUidCountByGroup(group string) int {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()
	groupMap, ok := g.Server.groupConnections[group]
	//群组都不存在
	if !ok {
		return 0
	}

	var sum int = 0
	for _, conn := range groupMap {
		_, ok2 := conn.TcpWsConnection().Get(constUid)
		if !ok2 {
			//没有uid
			continue
		}
		sum++
	}
	return sum
}

// GetAllUidList distributed-api
func (g *gatewayApi) GetAllUidList() []string {
	g.Server.uidConnectionsLock.RLock()
	defer g.Server.uidConnectionsLock.RUnlock()
	var uidList []string
	for uid, _ := range g.Server.uidConnections {
		uidList = append(uidList, uid)
	}
	return uidList
}

// GetAllUidCount distributed-api
func (g *gatewayApi) GetAllUidCount() int {
	g.Server.uidConnectionsLock.RLock()
	defer g.Server.uidConnectionsLock.RUnlock()
	return len(g.Server.uidConnections)
}

// GetAllGroupIdList distributed-api
func (g *gatewayApi) GetAllGroupIdList() []string {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()

	var groupList []string
	for groupId, _ := range g.Server.groupConnections {
		groupList = append(groupList, groupId)
	}

	return groupList
}

// GetAllGroupCount distributed-api
func (g *gatewayApi) GetAllGroupCount() int {
	g.Server.groupConnectionsLock.RLock()
	defer g.Server.groupConnectionsLock.RUnlock()

	return len(g.Server.groupConnections)
}
