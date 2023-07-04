package workerman_go

// ProtocolRegister 内部组件协议 - 注册结构体 ,目前 register在onMessage中使用，
type ProtocolRegister struct {
	//组件类型
	ComponentType int `json:"component_type"`
	//组件的名称
	Name                                string                              `json:"name"`
	ProtocolPublicGatewayConnectionInfo ProtocolPublicGatewayConnectionInfo `json:"protocol_public_gateway_connection_info"`
	Data                                string                              `json:"Data"`
	//这个字段由register标记
	Authed string `json:"authed"`
}

// ProtocolRegisterBroadCastComponentGateway 注册中心发出的广播 网关地址的指令
type ProtocolRegisterBroadCastComponentGateway struct {
	Msg         string                                `json:"msg"`
	Data        string                                `json:"Data"`
	GatewayList []ProtocolPublicGatewayConnectionInfo `json:"gateway_list"`
}
type ProtocolPublicGatewayConnectionInfo struct {
	GatewayAddr string `json:"gateway_addr"`
	GatewayPort string `json:"gateway_port"`
}

type ClientWs struct {
	ClientId string `json:"client_id"`
}

type ProtocolForwardUserOnClose struct {
	ClientId string `json:"client_id"`
}

type ProtocolForwardUserOnMessage struct {
	ClientId string `json:"client_id"`
	Message  string `json:"message"`
}
type ProtocolForwardUserOnConnect struct {
	ClientId string `json:"client_id"`
}

// ProtocolHeartbeat 心跳
type ProtocolHeartbeat struct {
	A int `json:"a"`
}

const (
	//ComponentIdentifiersTypeBusiness business类型的服务
	ComponentIdentifiersTypeBusiness = iota
	//ComponentIdentifiersTypeGateway gateway网管类型服务
	ComponentIdentifiersTypeGateway

	//CommandComponentHeartbeat 心跳指令
	CommandComponentHeartbeat
	//CommandComponentAuthRequest 结构体 ProtocolRegister 请求认证
	CommandComponentAuthRequest

	//CommandComponentGatewayList  business接受最新 []gateway列表指令 ProtocolRegisterBroadCastComponentGateway
	CommandComponentGatewayList

	//CommandGatewayForwardUserOnMessage 转发用户发来消息  ProtocolForwardUserOnMessage
	CommandGatewayForwardUserOnMessage
	// CommandGatewayForwardUserOnClose 对应的结构体 ProtocolForwardUserOnClose
	CommandGatewayForwardUserOnClose
	// CommandGatewayForwardUserOnConnect 对应的结构体 ProtocolForwardUserOnConnect
	CommandGatewayForwardUserOnConnect
	//CommandGatewayForwardUserOnError 暂未使用
	CommandGatewayForwardUserOnError

	// 官方28个接口
	GatewayCommandSendToAll
	GatewayCommandSendToClient
	GatewayCommandCloseClient
	GatewayCommandIsOnline
	GatewayCommandBindUid
	GatewayCommandUnbindUid
	GatewayCommandIsUidOnline
	GatewayCommandGetClientIdByUid
	GatewayCommandGetUidByClientId
	GatewayCommandSendToUid
	GatewayCommandJoinGroup
	GatewayCommandLeaveGroup
	GatewayCommandUngroup
	GatewayCommandSendToGroup
	GatewayCommandGetClientIdCountByGroup
	GatewayCommandGetClientSessionsByGroup
	GatewayCommandGetAllClientIdCount
	GatewayCommandGetAllClientSessions
	GatewayCommandSetSession
	GatewayCommandUpdateSession
	GatewayCommandGetSession
	GatewayCommandGetClientIdListByGroup
	GatewayCommandGetAllClientIdList
	GatewayCommandGetUidListByGroup
	GatewayCommandGetUidCountByGroup
	GatewayCommandGetAllUidList
	GatewayCommandGetAllUidCount
	GatewayCommandGetAllGroupIdList
	GatewayCommandGetAllGroupCount
)

type GcmdSendToAll struct {
	Data            []byte   `json:"Data"`
	ClientIdArray   []string `json:"client_id_array"`
	ExcludeClientId []string `json:"exclude_client_id"`
}

// GcmdSendToClient GatewayCommandSendToClient gpt-3.5-turbo
type GcmdSendToClient struct {
	ClientId string `json:"client_id"`
	SendData string `json:"send_data"`
}

// GcmdCloseClient GatewayCommandCloseClient gpt-3.5-turbo
type GcmdCloseClient struct {
	ClientId string `json:"client_id"`
}

// GcmdIsOnline GatewayCommandIsOnline gpt-3.5-turbo
type GcmdIsOnline struct {
	ClientId string `json:"client_id"`
}

// GcmdBindUid GatewayCommandBindUid gpt-3.5-turbo
type GcmdBindUid struct {
	ClientId string `json:"client_id"`
	Uid      string `json:"uid"`
}

// GcmdUnbindUid GatewayCommandUnbindUid gpt-3.5-turbo
type GcmdUnbindUid struct {
	ClientId string `json:"client_id"`
	Uid      string `json:"uid"`
}

// GcmdIsUidOnline GatewayCommandIsUidOnline gpt-3.5-turbo
type GcmdIsUidOnline struct {
	Uid string `json:"uid"`
}

// GcmdGetClientIdByUid GatewayCommandGetClientIdByUid gpt-3.5-turbo
type GcmdGetClientIdByUid struct {
	Uid string `json:"uid"`
}

// GcmdGetUidByClientId GatewayCommandGetUidByClientId gpt-3.5-turbo
type GcmdGetUidByClientId struct {
	ClientId string `json:"client_id"`
}

// GcmdSendToUid GatewayCommandSendToUid gpt-3.5-turbo
type GcmdSendToUid struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

// GcmdJoinGroup GatewayCommandJoinGroup gpt-3.5-turbo
type GcmdJoinGroup struct {
	ClientId string `json:"client_id"`
	Group    string `json:"group"`
}

// GcmdLeaveGroup GatewayCommandLeaveGroup gpt-3.5-turbo
type GcmdLeaveGroup struct {
	ClientId string `json:"client_id"`
	Group    string `json:"group"`
}

// GcmdUngroup GatewayCommandUngroup gpt-3.5-turbo
type GcmdUngroup struct {
	Group string `json:"group"`
}

// GcmdSendToGroup GatewayCommandSendToGroup gpt-3.5-turbo
type GcmdSendToGroup struct {
	Group           string   `json:"group"`
	Message         string   `json:"message"`
	ExcludeClientId []string `json:"exclude_client_id"`
}

// GcmdGetClientIdCountByGroup GatewayCommandGetClientIdCountByGroup gpt-3.5-turbo
type GcmdGetClientIdCountByGroup struct {
	Group string `json:"group"`
}

// GcmdGetClientSessionsByGroup GatewayCommandGetClientSessionsByGroup gpt-3.5-turbo
type GcmdGetClientSessionsByGroup struct {
	Group string `json:"group"`
}

// GcmdGetAllClientIdCount GatewayCommandGetAllClientIdCount gpt-3.5-turbo
type GcmdGetAllClientIdCount struct{}

// GcmdGetAllClientSessions GatewayCommandGetAllClientSessions gpt-3.5-turbo
type GcmdGetAllClientSessions struct{}

// GcmdSetSession GatewayCommandSetSession gpt-3.5-turbo
type GcmdSetSession struct {
	ClientId string    `json:"client_id"`
	Data     SessionKv `json:"Data"`
}

// GcmdUpdateSession GatewayCommandUpdateSession gpt-3.5-turbo
type GcmdUpdateSession struct {
	ClientId string    `json:"client_id"`
	Data     SessionKv `json:"Data"`
}

// GcmdGetSession GatewayCommandGetSession gpt-3.5-turbo
type GcmdGetSession struct {
	ClientId string `json:"client_id"`
}

// GcmdGetClientIdListByGroup GatewayCommandGetClientIdListByGroup gpt-3.5-turbo
type GcmdGetClientIdListByGroup struct {
	Group string `json:"group"`
}

// GcmdGetAllClientIdList GatewayCommandGetAllClientIdList gpt-3.5-turbo
type GcmdGetAllClientIdList struct{}

// GcmdGetUidListByGroup GatewayCommandGetUidListByGroup gpt-3.5-turbo
type GcmdGetUidListByGroup struct {
	Group string `json:"group"`
}

// GcmdGetUidCountByGroup GatewayCommandGetUidCountByGroup gpt-3.5-turbo
type GcmdGetUidCountByGroup struct {
	Group string `json:"group"`
}

// GcmdGetAllUidList GatewayCommandGetAllUidList gpt-3.5-turbo
type GcmdGetAllUidList struct{}

// GcmdGetAllUidCount GatewayCommandGetAllUidCount gpt-3.5-turbo
type GcmdGetAllUidCount struct{}

// GcmdGetAllGroupIdList GatewayCommandGetAllGroupIdList gpt-3.5-turbo
type GcmdGetAllGroupIdList struct{}

// GcmdGetAllGroupCount GatewayCommandGetAllGroupCount gpt-3.5-turbo
type GcmdGetAllGroupCount struct{}
