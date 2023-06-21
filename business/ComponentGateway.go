package business

import (
	"context"
	"crypto/tls"
	"fmt"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/url"
	"time"
)

// ComponentGateway 组件的连接实体
type ComponentGateway struct {
	root *Business

	Name string

	//连接地址远程
	Address string

	Port workerman_go.Port

	//远程实例
	ConnWs *websocket.Conn

	//是否在gateway认证通过
	Authd bool

	Ctx    context.Context
	Cancel context.CancelFunc
}

// Connect  business 收到列表后 尝试连接网关
func (g *ComponentGateway) Connect() {
	Scheme := ""

	if g.root.Config.TLS {
		Scheme = "wss"
	} else {
		Scheme = "ws"
	}
	// 设置business WebSocket客户端配置
	wsConfig := &websocket.Config{
		Location: &url.URL{
			Scheme: Scheme,
			Host:   g.Address,
			Path:   workerman_go.GatewayForBusinessWsPath,
		},
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		TlsConfig: &tls.Config{InsecureSkipVerify: g.root.Config.SkipVerify},
		Version:   websocket.ProtocolVersionHybi13,
		Origin:    &url.URL{Scheme: "http", Host: g.root.Config.RegisterPublicHostForComponent},
	}

	//锁gateway列表
	g.root.gatewayMapRWMutex.Lock()

	wsGateway, wsConnWithRegisterErr := websocket.DialConfig(wsConfig)

	if wsConnWithRegisterErr != nil {
		for wsConnWithRegisterErr != nil {
			wsGateway, wsConnWithRegisterErr = websocket.DialConfig(wsConfig)

			log.Println(wsConnWithRegisterErr)
			log.Printf("[%s]无法连接  gateway {%s://%s}  10秒后重连.. ", g.Name, Scheme, workerman_go.GatewayForBusinessWsPath)
			t := time.NewTicker(time.Second * 10)
			<-t.C
		}
	}

	//设置conn Fd
	g.ConnWs = wsGateway
	//解锁
	g.root.gatewayMapRWMutex.Unlock()

	//发送身份认证

	signSturct, err := workerman_go.GenerateSignTimeByte(workerman_go.CommandComponentAuthRequest, workerman_go.ProtocolRegister{
		ComponentType:                       workerman_go.ComponentIdentifiersTypeBusiness,
		Name:                                g.root.Name,
		ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
		Data:                                "workerman_go.ComponentIdentifiersTypeBusiness.auth",
		Authed:                              "0",
	}, g.root.Config.SignKey, func() time.Duration {
		return time.Second * 10
	})
	if err != nil {
		return
	}

	//可能发送失败,就是断开连接
	g.ConnWs.Write(signSturct.ToByte())

	//监听gateway的指令，例如请求认证 -> 用户消息
	g.ListenMessageSync()
}

// onClose 人工关闭，或者事件触发
func (g *ComponentGateway) onClose(gateway *ComponentGateway) {

	g.root.gatewayMapRWMutex.Lock()
	defer g.root.gatewayMapRWMutex.Unlock()

	gateway.ConnWs.Close()
	delete(g.root.gatewayMap, g.Address)
}

func (g *ComponentGateway) onMessage(gateway *ComponentGateway, msg string) {

}

// ListenMessageSync 阻塞循环连接
func (g *ComponentGateway) ListenMessageSync() {

	mapKey := g.Address
	//开个定时器，
	go func(gateway *ComponentGateway) {

		timerTick := time.NewTicker(time.Second * 10)
		for true {
			select {
			case <-gateway.Ctx.Done():
				return
			case <-timerTick.C:

				if gateway.Authd == false {
					//这里写发送认证指令
				}
			}
		}
	}(g)

	//todo 应对手动关闭，应该给与ctx信号，在人工删除前就return
	for true {
		CmdMsg := make([]byte, 1024*10)
		if _, ok := g.root.gatewayMap[mapKey]; ok == false {
			//已被删除
			return
		}
		n, err := g.ConnWs.Read(CmdMsg)
		if err != nil {
			g.onClose(g)
			log.Println("与gateway连接断开：", err)
			return
		}

		DataObj, err := workerman_go.ParseAndVerifySignJsonTime(CmdMsg[:n], g.root.Config.SignKey)
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		//阻塞处理gateway数据指令
		switch DataObj.Cmd {
		case workerman_go.CommandComponentAuthRequest:
		//要求认证
		case workerman_go.CommandGatewayForwardUserOnConnect:
			//用户连接gateway转发
			if g.root.OnMessage != nil {
				//g.root.OnMessage()
			}
		case workerman_go.CommandGatewayForwardUserMessage:
			//用户消息gateway

		}
	}

}
