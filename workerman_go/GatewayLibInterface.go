package workerman_go

// GatewayLibInterface SDK接口  gpt3.5-turbo
type GatewayLibInterface interface {

	//SendToAll 向所有客户端或者client_id_array指定的客户端发送$send_data数据。如果指定的$client_id_array中的client_id不存在则自动丢弃
	SendToAll(data []byte, client_id_array []string, exclude_client_id []string)

	//SendToClient 向客户端client_id发送$send_data数据。如果client_id对应的客户端不存在或者不在线则自动丢弃发送数据
	SendToClient(client_id string, send_data string)

	//CloseClient 断开与client_id对应的客户端的连接
	CloseClient(client_id string)

	//IsOnline 判断$client_id是否还在线
	IsOnline(client_id string) bool

	// 将client_id与uid绑定，以便通过Gateway::sendToUid($uid)发送数据，通过Gateway::isUidOnline($uid)用户是否在线。 uid解释：这里uid泛指用户id或者设备id，用来唯一确定一个客户端用户或者设备。
	BindUid(client_id string, uid string)

	//UnbindUid 将client_id与uid解绑。
	UnbindUid(client_id string, uid string)

	//IsUidOnline 判断$uid是否在线，此方法需要配合Gateway::bindUid($client_uid, $uid)使用。
	IsUidOnline(uid string) int

	//GetClientIdByUid 返回一个数组，数组元素为与uid绑定的所有在线的client_id。如果没有在线的client_id则返回一个空数组。
	GetClientIdByUid(uid string) []string

	//GetUidByClientId 返回client_id绑定的uid，如果client_id没有绑定uid，则返回null。
	GetUidByClientId(client_id string) string

	//SendToUid 向uid绑定的所有在线client_id发送数据。 默认uid与client_id是一对多的关系，如果当前uid下绑定了多个client_id，则多个client_id对应的客户端都会收到消息，这类似于PC QQ和手机QQ同时在线接收消息。
	SendToUid(uid string, message string)

	//JoinGroup 将client_id加入某个组，以便通过Gateway::sendToGroup发送数据。
	JoinGroup(client_id string, group string)

	//LeaveGroup 将client_id从某个组中删除，不再接收该分组广播(Gateway::sendToGroup)发送的数据。
	LeaveGroup(client_id string, group string)

	//Ungroup 取消分组，或者说解散分组。
	Ungroup(group string)

	//SendToGroup 向某个分组的所有在线client_id发送数据。
	SendToGroup(group string, message string, exclude_client_id []string)

	//GetClientIdCountByGroup 获取某分组当前在线成连接数（多少client_id在线）。
	GetClientIdCountByGroup(group string) int

	//GetClientSessionsByGroup 获取某个分组所有在线client_id信息。
	GetClientSessionsByGroup(group string) map[string]SessionKv

	//GetAllClientIdCount 获取当前在线连接总数（多少client_id在线）。
	GetAllClientIdCount() int

	//GetAllClientSessions 获取当前所有在线client_id信息。
	GetAllClientSessions() map[string]SessionKv

	//SetSession 设置某个client_id对应的session。如果对应client_id已经下线或者不存在，则会被忽略。
	SetSession(client_id string, data SessionKv)

	//UpdateSession 更新某个client_id对应的session。如果对应client_id已经下线或者不存在，则会被忽略。
	UpdateSession(client_id string, data SessionKv)

	//GetSession 获取某个client_id对应的session。
	GetSession(client_id string) SessionKv

	//GetClientIdListByGroup 获取某个分组所有在线client_id列表。
	GetClientIdListByGroup(group string) []string

	//GetAllClientIdList 获取全局所有在线client_id列表。
	GetAllClientIdList() []string

	//GetUidListByGroup 获取某个分组所有在线uid列表。
	GetUidListByGroup(group string) []string

	//GetUidCountByGroup 获取某个分组下的在线uid数量。
	GetUidCountByGroup(group string) int

	//GetAllUidList 获取全局所有在线uid列表。
	GetAllUidList() []string

	//GetAllUidCount 获取全局所有在线uid数量。
	GetAllUidCount() int

	//GetAllGroupIdList 获取全局所有在线group id列表。
	GetAllGroupIdList() []string

	// GetAllGroupCount 额外加的 GetAllGroupCount
	GetAllGroupCount() int
}

type SessionKv map[string]string
