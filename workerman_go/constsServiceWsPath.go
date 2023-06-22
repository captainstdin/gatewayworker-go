package workerman_go

// 用于注册中心的注册地址
const (
	RegisterForBusniessWsPath = "/component/ws/register/business"
	RegisterForGatewayWsPath  = "/component/ws/register/gateway"

	RegisterForComponent = "/component/ws/register"
)

// 用于gateway的内部定义监听地址
const (
	//GatewayForBusinessWsPath 用于business连接到gateway的websocket地址
	GatewayForBusinessWsPath = "/component/ws/gateway/full_duplex_business"

	GatewayForSdkWsPath = "/component/ws/gateway/sdk"
)
