package business

import (
	"gatewaywork-go/workerman_go"
	"net"
)

type ComponentClient struct {
	IpType workerman_go.IpType
	Ipv4   net.IPNet
	Ipv6   net.IPNet
	Port   workerman_go.Port
}
