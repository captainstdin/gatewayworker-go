package workerman_go

// ProtocolRegister 内部组件注册结构体 ,目前 register在onMessage中使用，
type ProtocolRegister struct {
	Command    int    `json:"command"`
	IsBusiness int    `json:"is_business"`
	IsGateway  int    `json:"is_gateway"`
	Data       string `json:"data"`
	Sign       int    `json:"sign"`
	TimeStamp  int    `json:"time_stamp"`
}

const (
	//请求认证
	CommandServiceAuthRequest = iota
	//认证回响
	CommandServiceAuthResponse
	//广播 business指令
	CommandServiceBroadcastBusiness
)
