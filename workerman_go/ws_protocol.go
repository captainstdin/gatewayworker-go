package workerman_go

type ProtocolJsonRegister struct {
	Command    int `json:"command"`
	IsBusiness int `json:"is_business"`
	IsGateway  int `json:"is_gateway"`
}

const (
	//请求认证
	CommandServiceAuthRequest = iota
	//认证回响
	CommandServiceAuthResponse
	//广播 business指令
	CommandServiceBroadcastBusiness
)
