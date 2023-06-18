package main

import (
	"fmt"
	"gatewaywork-go/business"
	"gatewaywork-go/register"
	"gatewaywork-go/workerman_go"
	"os"
	"sync"
)

var coroutine sync.WaitGroup

func A() {

}

func main() {

	Conf := workerman_go.ConfigGatewayWorker{
		RegisterListenAddr:             ":1238",
		RegisterListenPort:             ":1238",
		TLS:                            false,
		TlsKeyPath:                     "",
		TlsPemPath:                     "",
		RegisterPublicHostForComponent: "127.0.0.1:1238",
		GatewayPublicHostForClient:     "",
		GatewayListenAddr:              "",
		GatewayListenPort:              "",
		SkipVerify:                     false,
		SignKey:                        "da!!bskdhaskld#1238asjiocy89123",
	}

	if register_enable := os.Getenv("register_enable"); register_enable == "1" {
		coroutine.Add(1)
		go StartRegister(&Conf)
	}

	if business_enable := os.Getenv("business_enable"); business_enable == "1" {
		coroutine.Add(1)
		go StartBusiness(&Conf)
	}

	coroutine.Wait()
}

func StartBusiness(Conf *workerman_go.ConfigGatewayWorker) {
	defer coroutine.Done()
	service := business.NewBusiness("Business处理器", Conf)
	err := service.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func StartRegister(Conf *workerman_go.ConfigGatewayWorker) {
	defer coroutine.Done()
	service := register.NewRegister("Business处理器", Conf)
	err := service.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
