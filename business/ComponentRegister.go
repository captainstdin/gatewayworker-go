package business

import (
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"sync"
	"time"
)

// ComponentRegister 组件的连接实体--注册中心
type ComponentRegister struct {
	addr   string //包含Port
	ConnWs *websocket.Conn
	//读写锁
	RWLock *sync.RWMutex

	root *Business

	ClientId *workerman_go.ClientToken
}

// 人工关闭的时候，或者事件触发
func (r *ComponentRegister) onClose(register *ComponentRegister) {
	r.root.registerMapRWMutex.Lock()
	//删除list
	register.ConnWs.Close()
	delete(r.root.registerMap, uint64(register.ClientId.ClientGatewayNum))
	r.root.registerMapRWMutex.Unlock()
}

func (r *ComponentRegister) onMessage(WsConn *ComponentRegister, data *workerman_go.GenerateComponentSign) {

	cmd := data.Cmd

	switch cmd {

	case workerman_go.CommandComponentAuthRequest:

		buffObj, err := workerman_go.GenerateSignTimeByte(workerman_go.CommandComponentAuthRequest, workerman_go.ProtocolRegister{
			ComponentType:                       workerman_go.ComponentIdentifiersTypeBusiness,
			Name:                                "",
			ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
			Data:                                "workerman_go.ComponentIdentifiersTypeBusiness.auth",
			Authed:                              "",
		}, r.root.Config.SignKey, func() time.Duration {
			return time.Second * 60
		})
		if err != nil {
			return
		}

		WsConn.ConnWs.Write(buffObj.ToByte())

	case workerman_go.CommandComponentGatewayList:
		//获取 网关列表
		var gatewayList workerman_go.ProtocolRegisterBroadCastComponentGateway
		err := json.Unmarshal(data.Json, &gatewayList)
		if err != nil {
			log.Println("解析 register发来的gatewayList数据错误：", err.Error())
			return
		}

		for _, gatewayInstance := range gatewayList.GatewayList {
			r.root.IncGatewayConn(gatewayInstance)
		}

	}

}

// IncGatewayConn 收到gateway列表广播，去连接gateway
func (b *Business) IncGatewayConn(gateway workerman_go.ProtocolPublicGatewayConnectionInfo) {
	//锁定gateway
	b.gatewayMapRWMutex.Lock()
	instance := &ComponentGateway{
		root:    b,
		Name:    gateway.GatewayAddr,
		Address: gateway.GatewayAddr,
		ConnWs:  nil,
		Authd:   false,
	}

	//这里的地址一定是唯一的
	b.gatewayMap[gateway.GatewayAddr] = instance
	b.gatewayMapRWMutex.Unlock()

	//读锁连接
	b.gatewayMapRWMutex.RLock()
	for _, gatewayInstance := range b.gatewayMap {
		if gatewayInstance.Authd == true {
			continue
		}
		gatewayInstance.Connect()
	}
	b.gatewayMapRWMutex.RUnlock()
}

func (r *ComponentRegister) ListenMessageSync() {
	for true {
		CmdMsg := make([]byte, 1024*10)
		n, err := r.ConnWs.Read(CmdMsg)

		if err != nil {
			r.onClose(r)
			log.Println("与register连接断开：", err)
			return
		}

		DataObj, err := workerman_go.ParseAndVerifySignJsonTime(CmdMsg[:n], r.root.Config.SignKey)
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		//阻塞
		r.onMessage(r, DataObj)
	}

}
