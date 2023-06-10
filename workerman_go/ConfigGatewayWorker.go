package workerman_go

type ConfigGatewayWorker struct {

	//注册发现，监听的地址和权限，{仅本地,仅局域网,全部}
	RegisterListenAddr string `json:"register_listen_addr"`
	//注册发现，监听的端口
	RegisterListenPort string `json:"register_listen_port"`

	//组件之间内部通讯是否开启wss
	TLS bool `json:"tls"`

	//组件之间内部通讯 key路径
	TlsKeyPath string `json:"tls_key_path"`
	//组件之间内部通讯 pem 路径
	TlsPemPath string `json:"tls_pem_path"`

	//公网地址
	RegisterPublicHostForRegister string `json:"register_public_addr"`

	//gateway 公网连接地址
	GatewayPublicHostForClient string `json:"gateway_public_host_for_client"`

	//gateway 监听地址
	GatewayListenAddr string `json:"gateway_listen_addr"`
	//gateway 监听端口
	GatewayListenPort string `json:"gateway_listen_port"`

	//是否跳过证书验证，自签证书 请:=true
	SkipVerify bool `json:"skip_verify"`

	//三个组件内部互相通讯的签名密钥
	SignKey string `json:"sign_key"`
}
