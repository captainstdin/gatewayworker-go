package workerman_go

import (
	"fmt"
	"net"
	"strconv"
)

type ConfigGatewayWorker struct {
	RegisterEnable string `json:"register_enable"`
	GatewayEnable  string `json:"gateway_enable"`
	BusinessEnable string `json:"business_enable"`

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
type IpType uint8

const IpTypeV4 = IpType(4)
const IpTypeV6 = IpType(6)

func (c *ConfigGatewayWorker) GetGatewayPublicAddress() (net.IP, int, IpType, error) {

	host, port, err := net.SplitHostPort(c.GatewayPublicHostForClient)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("解析地址时出错: %s\n", err)
	}
	// 检查主机部分是否是IPv4或IPv6地址
	ip := net.ParseIP(host)

	portInt, _ := strconv.Atoi(port)
	if ip != nil {
		// 是一个有效的IP地址
		if ip.To4() != nil {
			return ip, portInt, IpTypeV4, nil
		}

		if ip.To16() != nil {
			return ip, portInt, IpTypeV6, nil
		}
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("解析域名时出错:%s\n", err)
	}

	domainIp := ips[0]

	if domainIp.To4() != nil {
		return domainIp, portInt, IpTypeV4, nil
	}

	if domainIp.To16() != nil {
		return domainIp, portInt, IpTypeV6, nil
	}

	return nil, 0, 0, fmt.Errorf("解析域名后，无法判断ipv4/6: %s\n", err)

}
