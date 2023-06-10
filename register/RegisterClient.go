package register

import (
	"encoding/json"
	"errors"
	"gatewaywork-go/workerman_go"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

const (
	AuthedName = iota
	ServiceTypeName
)

// 是否通过认证
type Authed bool

// ip类型
type ServiceType uint8

type RegisterClient struct {
	RegisterService *Register
	Address         string
	Port            string
	FdWs            *websocket.Conn
	DataRWMutex     *sync.RWMutex
	Data            map[string]interface{}
	Request         *http.Request
	//TokenStructString
	ClientToken workerman_go.ClientToken
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
	rc.RegisterService.InnerOnClose(rc)
}

func (conn *RegisterClient) Send(data interface{}) error {

	if cmd, cmdOk := data.(workerman_go.ProtocolRegister); cmdOk {
		str, _ := json.Marshal(cmd)
		err := conn.FdWs.WriteMessage(websocket.TextMessage, str)
		if err != nil {
			conn.Close()
			return err
		}
		return nil
	}

	if str, strOk := data.(string); strOk {
		err := conn.FdWs.WriteMessage(websocket.TextMessage, []byte(str))
		if err != nil {
			conn.Close()
			return err

		}
		return nil
	}

	if byteStr, byteOk := data.([]byte); byteOk {
		err := conn.FdWs.WriteMessage(websocket.TextMessage, byteStr)
		if err != nil {
			conn.Close()
			return err
		}
		return nil
	}
	return errors.New("conn.Send(Unknown protocol message)")
}

func (conn *RegisterClient) GetRemoteIp() string {
	return ""
}

func (conn *RegisterClient) GetRemotePort() string {

	return ""
}

func (conn *RegisterClient) PauseRecv() {

}
func (conn *RegisterClient) ResumeRecv() {

}

func (conn *RegisterClient) Pipe(connection *workerman_go.TcpConnection) {

}

func (conn *RegisterClient) GetClientId() string {
	return conn.ClientToken.GenerateGatewayClientId()
}

func (conn *RegisterClient) GetClientIdInfo() *workerman_go.ClientToken {
	return &conn.ClientToken
}

func (conn *RegisterClient) Get(str string) (interface{}, bool) {
	//读锁，防止读的时候写
	conn.DataRWMutex.RLock()
	defer conn.DataRWMutex.RLock()
	item, ok := conn.Data[str]
	return item, ok
}

func (conn *RegisterClient) Set(str string, v interface{}) {
	//写锁，防止读
	conn.DataRWMutex.Lock()
	defer conn.DataRWMutex.Unlock()
	conn.Data[str] = v
}
