package gateway

import (
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

// ComponentClient 每个连接上来的ws Client主要是  component组件(business)与 WebSocket用户
type ComponentClient struct {
	root *GatewayServer
	//生成的在当前内部组件中标志目标gateway所在地
	ClientId *workerman_go.ClientToken

	//组件名称
	Name string

	//组件类型
	ComponentType int

	//连接地址
	Address string
	Port    workerman_go.Port

	FdWs *websocket.Conn
}

// Listen  连接
func (g *ComponentClient) Listen(f func(sign *workerman_go.GenerateComponentSign, client *ComponentClient)) {

	for true {

		msg := make([]byte, 10240)
		n, err := g.FdWs.Read(msg)
		if err != nil {
			log.Println("【ComponentClient】离线", err)
			g.onClose(g)
			return
		}

		data, err := workerman_go.ParseAndVerifySignJsonTime(msg[:n], g.root.Config.SignKey)
		if err != nil {
			log.Println("【ComponentClient 协议异常】离线", err)
			g.onClose(g)
			return
		}
		f(data, g)
	}
}

// onClose 连接
func (g *ComponentClient) onClose(tcpConnect *ComponentClient) {

	g.root.ComponentsLock.Lock()

	defer g.root.ComponentsLock.Unlock()

	delete(g.root.ComponentsMap, tcpConnect.ClientId.ClientGatewayNum)
}

func (g *ComponentClient) Close() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) send(b interface{}) error {

	timeByte, err := workerman_go.GenerateSignTimeByte(workerman_go.CommandComponentAuthRequest, b, g.root.Config.SignKey, func() time.Duration {
		return 10 * time.Second
	})
	if err != nil {
		return err
	}

	_, err2 := g.FdWs.Write(timeByte.ToByte())
	g.onClose(g)
	if err2 != nil {
		return err2
	}
	return nil
}

func (g *ComponentClient) Send(data interface{}) error {

	switch data.(type) {
	case []byte:
	case workerman_go.ProtocolRegister:
		err := g.send(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *ComponentClient) ListenFor() {

}

func (g *ComponentClient) GetRemoteIp() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetRemotePort() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Pipe(connection *workerman_go.TcpConnection) {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetClientId() string {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) GetClientIdInfo() *workerman_go.ClientToken {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Get(str string) (interface{}, bool) {
	//TODO implement me
	panic("implement me")
}

func (g *ComponentClient) Set(str string, v interface{}) {
	//TODO implement me
	panic("implement me")
}
