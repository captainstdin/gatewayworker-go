package business

import (
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"sync"
)

// ComponentRegister 组件的连接实体--注册中心
type ComponentRegister struct {
	addr   string //包含Port
	ConnWs *websocket.Conn
	//读写锁
	RWLock *sync.RWMutex

	root *Business
}

func (r *ComponentRegister) onClose(register *ComponentRegister) {

}

func (r *ComponentRegister) onMessage(data *workerman_go.GenerateComponentSign) {

	cmd := data.Cmd

	switch cmd {

	case workerman_go.CommandComponentGatewayList:
		//获取 网关列表

		var gatewayList workerman_go.ProtocolRegisterBroadCastComponentGateway
		err := json.Unmarshal(data.Json, &gatewayList)
		if err != nil {
			log.Println("解析 register发来的gatewayList数据错误：", err.Error())
			return
		}

		for _, gatewayInstance := range gatewayList.GatewayList {

			r.root.gatewayMapRWMutex.Lock()

			fmt.Println(gatewayInstance.GatewayAddr)
			r.root.gatewayMapRWMutex.Unlock()
		}

	}

}

func (r *ComponentRegister) ListenMessage() {

	for true {

		CmdMsg := make([]byte, 1024*10)
		n, err := r.ConnWs.Read(CmdMsg)

		fmt.Println(string(CmdMsg[:n]))
		if err != nil {
			r.onClose(r)
			return
		}

		DataObj, err := workerman_go.ParseAndVerifySignJsonTime(CmdMsg[:n], r.root.Config.SignKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		r.onMessage(DataObj)
	}

}
