package register

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"workerman_go/workerman_go"
)

type RegisterClient struct {
	RegisterService *Register
	Address         string
	Port            string
	FdWs            *websocket.Conn
	Data            map[string]string
	Request         *http.Request
	//是否通过认证
	Authed bool
	//TokenStructString
	ClientToken workerman_go.ClientToken

	//ip类型
	ServiceType uint8
}

func (conn *RegisterClient) SendCommand(v interface{}) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = conn.FdWs.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		//close conn
		conn.FdWs.Close()
		return
	}

}

// 主动关闭接口
func (rc *RegisterClient) Close() {
	rc.FdWs.Close()
	rc.RegisterService._OnClose(rc)
}

func (conn *RegisterClient) Send(data interface{}) {

}

func (conn *RegisterClient) getRemoteIp() string {
	return ""
}

func (conn *RegisterClient) getRemotePort() string {

	return ""
}

func (conn *RegisterClient) PauseRecv() {

}
func (conn *RegisterClient) ResumeRecv() {

}

func (conn *RegisterClient) Pipe(connection *workerman_go.TcpConnection) {

}
