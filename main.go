package main

import "workerman_go/register"

func main() {

	go func() {
		register.NewRegister()
	}()
	select {}
}
