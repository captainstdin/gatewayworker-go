package main

import (
	"fmt"
	"gatewaywork-go/register"
	"gatewaywork-go/workerman_go"
	"sync"
)

var coroutine sync.WaitGroup

func main() {

	Conf := workerman_go.ConfigGatewayWorker{
		RegisterListenAddr:            ":1238",
		RegisterListenPort:            ":1238",
		TLS:                           false,
		TlsKeyPath:                    "",
		TlsPemPath:                    "",
		RegisterPublicHostForRegister: "127.0.0.1:1237",
		GatewayPublicHostForClient:    "",
		GatewayListenAddr:             "",
		GatewayListenPort:             "",
		SkipVerify:                    false,
		SignKey:                       "da!!bskdhaskld#1238asjiocy89123",
	}

	StartRegister(&Conf)
	coroutine.Wait()
}

func StartRegister(Conf *workerman_go.ConfigGatewayWorker) {
	coroutine.Add(1)
	go func() {
		defer coroutine.Done()
		service := register.NewRegister("Business处理器", Conf)
		err := service.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}
