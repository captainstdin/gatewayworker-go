package workerman_go

// ProtocolCommandName 必须是每一个protocol指令必须共有的，发送协议或者接受信息的时候来区分
const ProtocolCommandName = "command"

// ProtocolRegister 内部组件协议 - 注册结构体 ,目前 register在onMessage中使用，
type ProtocolRegister struct {
	Command    string `json:"command"`
	IsBusiness string `json:"is_business"`
	IsGateway  string `json:"is_gateway"`
	Data       string `json:"data"`
	//这个字段由register标记
	Authed string `json:"authed"`
}

// 服务类型
const ComponentType = "ComponentType"
const ComponentAuthed = "ComponentAuthed"

const ComponentLastHeartbeat = "ComponentLastHeartbeat"
const ConstSignFieldName = "sign"
const (
	ConstSignBy             = iota
	ConstSignTokenTimeStamp //timestamp
	//ComponentTypeBusiness business类型的服务
	ComponentTypeBusiness
	//ComponentTypeGateway gateway网管类型服务
	ComponentTypeGateway

	//CommandComponentHeartbeat 心跳指令
	CommandComponentHeartbeat
	//CommandComponentAuthRequest 请求认证
	CommandComponentAuthRequest
	//CommandComponentAuthResponse 认证回响
	CommandComponentAuthResponse
	// CommandServiceBroadcastBusiness 广播 business指令
	CommandServiceBroadcastBusiness
)
