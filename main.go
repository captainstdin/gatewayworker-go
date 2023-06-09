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
		service := register.NewRegister()
		service.Run()
	}()
}
