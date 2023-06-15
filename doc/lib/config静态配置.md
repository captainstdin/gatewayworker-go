# 全局配置

## 常见问题
>### 1： 如何同时让组件之间通讯支持`ipv4`&&`ipv6`？

#### 回答： 以下两种均支持`ipv4`同时支持`ipv6`

设置监听地址`register_listen_addr`为`0.0.0.0:端口号`

如果想只允许本机访问 `127.0.0.1:端口号`

>### 2： 监听地址的写法和`ipv4`与`ipv6`的关系

1. `0.0.0.0:端口` 外网&&本地&&局域网用户可以使用`ipv4`&&`ipv6`访问
2. `127.0.01:端口` 仅本机可以使用`ipv4`&&`ipv6`访问
3. `192.168.3.15:端口` 仅局域网用户可以使用`ipv4`访问
4. `[2409:8a28:3e21:2810::fc]:端口` 外网用户可以使用`ipv6`访问
5. `[::1]:端口` 外网&&本地&&局域网用户使用`ipv6`访问

### 如果你是开发者：

1： 如何同时让组件之间通讯支持`ipv4`&&`ipv6`？

文件参考：`/workerman_go/ConfigGatewayWorker.go`
```
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

	//RegisterPublicHostForComponent   其他组件通过这个公网地址经行连接 例如 {"baidIpv6.com","baidu.com:1237"}
	RegisterPublicHostForComponent string `json:"register_public_host_for_component"`

	//gateway 其他组件通过这个地址经行连接 例如 {"baidIpv6.com","baidu.com:2727"}
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
```
 