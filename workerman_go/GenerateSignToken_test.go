package workerman_go

import (
	"testing"
	"time"
)

func TestGenerateSignTimeByte(t *testing.T) {
	timeByte, err := GenerateSignTimeByte(CommandComponentAuthRequest, ProtocolRegister{
		ComponentType:                       0,
		Name:                                "register",
		ProtocolPublicGatewayConnectionInfo: ProtocolPublicGatewayConnectionInfo{},
		Data:                                "aaa",
		Authed:                              "1",
	}, "123123", func() time.Duration {
		return time.Second * 60
	})
	if err != nil {
		return
	}
	t.Log(timeByte.ToString())
}

func TestParseAndVerifySignJsonTime1(t *testing.T) {
	key := "123123"
	timeByte, err := GenerateSignTimeByte(CommandComponentAuthRequest, ProtocolRegister{
		ComponentType:                       0,
		Name:                                "register",
		ProtocolPublicGatewayConnectionInfo: ProtocolPublicGatewayConnectionInfo{},
		Data:                                "aaa",
		Authed:                              "1",
	}, key, func() time.Duration {
		return time.Second * 60
	})
	if err != nil {
		return
	}

	//fmt.Println(timeByte.ToString())

	jsonTime, err := ParseAndVerifySignJsonTime(timeByte.ToByte(), key)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("解析成功，内容如下JSON")
	t.Log(string(jsonTime.Json))
	//fmt.Printf("%s", jsonTime.Json)
}
