package workerman_go

// ProtocolRegister 内部组件协议 - 注册结构体 ,目前 register在onMessage中使用，
type ProtocolRegister struct {
	//组件类型
	ComponentType int `json:"component_type"`
	//组件的名称
	Name                                string                              `json:"name"`
	ProtocolPublicGatewayConnectionInfo ProtocolPublicGatewayConnectionInfo `json:"protocol_public_gateway_connection_info"`
	Data                                string                              `json:"data"`
	//这个字段由register标记
	Authed string `json:"authed"`
}

// ProtocolRegisterBroadCastComponentGateway 注册中心发出的广播 网关地址的指令
type ProtocolRegisterBroadCastComponentGateway struct {
	Msg         string                                `json:"msg"`
	Data        string                                `json:"data"`
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
	//CommandComponentAuthRequest 请求认证
	CommandComponentAuthRequest

	//CommandComponentGatewayList  business接受最新 []gateway列表指令
	CommandComponentGatewayList
	// CommandServiceBroadcastBusiness 广播 business指令
	CommandServiceBroadcastBusiness

	//CommandGatewayForwardUserOnMessage 转发用户发来消息
	CommandGatewayForwardUserOnMessage
	CommandGatewayForwardUserOnClose
	CommandGatewayForwardUserOnConnect
	CommandGatewayForwardUserOnError
)
