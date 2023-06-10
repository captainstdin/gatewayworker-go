package main

import (
	"gatewaywork-go/register"
	"sync"
)

var coroutine sync.WaitGroup

func main() {

	StartRegister()
	coroutine.Wait()
}

func StartRegister() {
	coroutine.Add(1)
	go func() {
		defer coroutine.Done()
		service := register.NewRegister("", "1237", "Business处理器")
		service.Run()
	}()
}
