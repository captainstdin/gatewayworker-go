package gateway

import (
	"gatewaywork-go/workerman_go"
	"testing"
)

var Conf = workerman_go.ConfigGatewayWorker{
	RegisterListenAddr:             "0.0.0.0:1238",
	RegisterListenPort:             ":1238",
	TLS:                            false,
	TlsKeyPath:                     "",
	TlsPemPath:                     "",
	RegisterPublicHostForComponent: "127.0.0.1:1238",
	GatewayPublicHostForClient:     "",
	GatewayListenAddr:              ":2727",
	SkipVerify:                     false,
	SignKey:                        "da!!bskdhaskld#1238asjiocy89123",
}

func TestNewGatewayServer(t *testing.T) {

	gateway := NewGatewayServer("gateway", &Conf)

	gateway.Run()
}
