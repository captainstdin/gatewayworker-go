package business

import (
	"gatewaywork-go/workerman_go"
	"net"
)

// ComponentGateway 组件的连接实体
type ComponentGateway struct {
	Address string

	Port   workerman_go.Port
	ConnWs *net.Conn

	Authd bool
}

// Connect 连接网关
func (g *ComponentGateway) Connect() {

	ok := false

	for ok == false {

	}
}

func (g *ComponentGateway) ListenMessage() {

	//阻塞循环连接
	g.Connect()

}
