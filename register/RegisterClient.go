package register

import (
	"errors"
	"gatewaywork-go/workerman_go"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Authed 是否通过认证
type Authed bool

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

// Close 主动关闭接口,会触发InnerOnClose()
func (rc *RegisterClient) Close() {
	rc.FdWs.Close()
	rc.RegisterService.InnerOnClose(rc)
}

// sendWithSignJsonString 内部方法
func (rc *RegisterClient) sendWithSignJsonString(v any) error {
	jsonString, err := workerman_go.GenerateSignJsonTime(v, rc.RegisterService.GatewayWorkerConfig.SignKey, func() time.Duration {
		return time.Second * 10
	})

	if err != nil {
		return err
	}

	sendErr := rc.FdWs.WriteMessage(websocket.TextMessage, []byte(jsonString))
	if sendErr != nil {
		rc.Close()
		return sendErr
	}
	return nil
}

// CommandToComponentForAllList Broadcast
func (rc *RegisterClient) CommandToComponentForAllList() {
	//rc.RegisterService.
}

// CommandToComponentForAuthRequire 要求发送身份验证
func (rc *RegisterClient) CommandToComponentForAuthRequire() {
	rc.sendWithSignJsonString(workerman_go.ProtocolRegister{
		//请求授权标志
		Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
		Data:    "workerman_go.CommandServiceAuthRequest",
		Authed:  strconv.Itoa(0), //告诉组件未授权
	})
}

// Send 发送json数据，但是带有签名校验和时间校验的
func (rc *RegisterClient) Send(data any) error {

	if register, registerOk := data.(workerman_go.ProtocolRegister); registerOk {
		rc.sendWithSignJsonString(register)
	}
	return errors.New("conn.Send(Unknown protocol message)")
}

func (rc *RegisterClient) GetRemoteIp() string {
	return ""
}

func (rc *RegisterClient) GetRemotePort() string {
	return ""
}

func (rc *RegisterClient) PauseRecv() {

}
func (rc *RegisterClient) ResumeRecv() {
}

func (rc *RegisterClient) Pipe(connection *workerman_go.TcpConnection) {

}

func (rc *RegisterClient) GetClientId() string {
	return rc.ClientToken.GenerateGatewayClientId()
}

func (rc *RegisterClient) GetClientIdInfo() *workerman_go.ClientToken {
	return &rc.ClientToken
}

func (rc *RegisterClient) Get(str string) (interface{}, bool) {
	//读锁，防止读的时候写
	rc.DataRWMutex.RLock()
	defer rc.DataRWMutex.RLock()
	item, ok := rc.Data[str]
	return item, ok
}

func (rc *RegisterClient) Set(str string, v interface{}) {
	//写锁，防止读
	rc.DataRWMutex.Lock()
	defer rc.DataRWMutex.Unlock()
	rc.Data[str] = v
}
