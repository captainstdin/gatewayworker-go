package workerman_go

import (
	"testing"
)

// TestClientId_GenerateGatewayClientId 测试 Auth:GPT-3.5-turbo
func TestClientId_GenerateGatewayClientId(t *testing.T) {

	test4 := GatewayIdInfo{
		ClientGatewayAddr: "192.168.0.1",
		//ClientGatewayPort: Port(3306),
		ClientGatewayNum: 1,
	}
	t.Log("Gen :Ipv4 ClientID", test4.GenerateGatewayClientId())
	id, _ := ParseGatewayClientId(test4.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id)

	test6 := GatewayIdInfo{
		//IPType:            IpTypeV6,
		//ClientGatewayIpv4: nil,
		//ClientGatewayIpv6: net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		ClientGatewayAddr: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		//ClientGatewayPort: Port(3306),
		ClientGatewayNum: 1,
	}
	t.Log("Gen :Ipv6 ClientID", test6.GenerateGatewayClientId())
	id2, _ := ParseGatewayClientId(test6.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id2)

	test6_compressed := GatewayIdInfo{
		//IPType:            IpTypeV4,
		//ClientGatewayIpv4: nil,
		//ClientGatewayIpv6: net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		ClientGatewayAddr: "2001:db8:85a3::8a2e:370:7334",
		//ClientGatewayPort: Port(3306),
		ClientGatewayNum: 1,
	}
	t.Log("Gen :Ipv6 ClientID", test6_compressed.GenerateGatewayClientId(), "[compressed")
	id3, _ := ParseGatewayClientId(test6_compressed.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id3)

	test_host := GatewayIdInfo{
		//IPType:            IpTypeV4,
		//ClientGatewayIpv4: nil,
		//ClientGatewayIpv6: net.ParseIP("2001:db8:85a3::8a2e:370:7334"),
		ClientGatewayAddr: "www.example.com",
		//ClientGatewayPort: Port(3306),
		ClientGatewayNum: 1,
	}
	t.Log("Gen :domain(www.example.com) ClientID", test6.GenerateGatewayClientId())
	id4, _ := ParseGatewayClientId(test_host.GenerateGatewayClientId())
	t.Logf("Parse :%+v", id4)

}

func TestGenPrimaryKeyUint64(t *testing.T) {

}
