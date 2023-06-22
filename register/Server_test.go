package register

import (
	"gatewaywork-go/workerman_go"
	"testing"
)

func TestServer_Run(t *testing.T) {

	conf := &workerman_go.ConfigGatewayWorker{
		RegisterEnable:                 "1",
		GatewayEnable:                  "0",
		BusinessEnable:                 "0",
		RegisterListenAddr:             ":1238",
		RegisterListenPort:             "",
		TLS:                            false,
		TlsKeyPath:                     "",
		TlsPemPath:                     "",
		RegisterPublicHostForComponent: "",
		GatewayPublicHostForClient:     "",
		GatewayListenAddr:              "",
		GatewayListenPort:              "",
		SkipVerify:                     false,
		SignKey:                        "",
	}
	S := NewServer("Register", conf)
	S.Run()
}
