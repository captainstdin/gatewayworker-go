package workerman_go

// ProtocolCommandName 必须是每一个protocol指令必须共有的，发送协议或者接受信息的时候来区分
const ProtocolCommandName = "command"

// ProtocolRegister 内部组件协议 - 注册结构体 ,目前 register在onMessage中使用，
type ProtocolRegister struct {
	Command       int `json:"command"`
	ComponentType int `json:"component_type"`

	//组件的名称
	Name                                string                              `json:"name"`
	ProtocolPublicGatewayConnectionInfo ProtocolPublicGatewayConnectionInfo `json:"protocol_public_gateway_connection_info"`

	Data string `json:"data"`
	//这个字段由register标记
	Authed string `json:"authed"`
}

// ProtocolRegisterBroadCastComponentGateway 注册中心发出的广播 网关地址的指令
type ProtocolRegisterBroadCastComponentGateway struct {
	Command     int                                   `json:"command"`
	Data        string                                `json:"data"`
	GatewayList []ProtocolPublicGatewayConnectionInfo `json:"gateway_list"`
}

type ProtocolPublicGatewayConnectionInfo struct {
	GatewayAddr string `json:"gateway_addr"`
	GatewayPort string `json:"gateway_port"`
}

// 服务类型
const ComponentType = "ComponentType"
const ComponentAuthed = "ComponentAuthed"
const (
	ComponentIdentifiersAuthed = "ComponentIdentifiersAuthed"
	ComponentIdentifiersType   = "ComponentIdentifiersType"
	//ComponentIdentifiersTypeBusiness business类型的服务
	ComponentIdentifiersTypeBusiness = iota
	//ComponentIdentifiersTypeGateway gateway网管类型服务
	ComponentIdentifiersTypeGateway
)
const ComponentLastHeartbeat = "ComponentLastHeartbeat"
const ConstSignFieldName = "sign"
const (
	ConstSignBy             = iota
	ConstSignTokenTimeStamp //timestamp

	//CommandComponentHeartbeat 心跳指令
	CommandComponentHeartbeat
	//CommandComponentAuthRequest 请求认证
	CommandComponentAuthRequest
	//CommandComponentAuthResponse 认证回响
	CommandComponentAuthResponse

	//CommandComponentGatewayListResponse  business接受最新 []gateway列表指令
	CommandComponentGatewayListResponse
	// CommandServiceBroadcastBusiness 广播 business指令
	CommandServiceBroadcastBusiness
)
