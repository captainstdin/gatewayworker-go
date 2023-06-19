package workerman_go

import (
	"net"
	"testing"
)

// TestClientId_GenerateGatewayClientId 测试 Auth:GPT-3.5-turbo
func TestClientId_GenerateGatewayClientId(t *testing.T) {

	test4 := ClientToken{
		IPType:            IpTypeV4,
		ClientGatewayIpv4: net.ParseIP("192.168.0.1"),
		ClientGatewayIpv6: nil,
		ClientGatewayPort: Port(3306),
		ClientGatewayNum:  GatewayNum(1),
	}

	t.Log("Gen :Ipv4 ClientID", test4.GenerateGatewayClientId())

	id, _ := ParseGatewayClientId(test4.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id)

	test6 := ClientToken{
		IPType:            IpTypeV6,
		ClientGatewayIpv4: nil,
		ClientGatewayIpv6: net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		ClientGatewayPort: Port(3306),
		ClientGatewayNum:  GatewayNum(1),
	}

	t.Log("Gen :Ipv6 ClientID", test6.GenerateGatewayClientId())
	id2, _ := ParseGatewayClientId(test6.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id2)
	test6_compressed := ClientToken{
		IPType:            IpTypeV4,
		ClientGatewayIpv4: nil,
		ClientGatewayIpv6: net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		ClientGatewayPort: Port(3306),
		ClientGatewayNum:  GatewayNum(1),
	}
	t.Log("Gen :Ipv6 ClientID", test6_compressed.GenerateGatewayClientId(), "[compressed")
	id3, _ := ParseGatewayClientId(test6_compressed.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id3)
}

func TestGenPrimaryKeyUint64(t *testing.T) {
	//type args struct {
	//	mapData map[uint64]interface{}
	//}
	//
	//mapdata := make(map[uint64]interface{})
	//
	//mapdata[uint64(123)] = "business.ComponentGateway{}"
	//
	//tests := []struct {
	//	name string
	//	args args
	//	want uint64
	//}{
	//	{name: "你好", args: mapdata, want: uint64(123)},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := GenPrimaryKeyUint64(tt.args.mapData); got != tt.want {
	//			t.Errorf("GenPrimaryKeyUint64() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
