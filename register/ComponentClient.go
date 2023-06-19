package register

import (
	"context"
	"errors"
	"gatewaywork-go/workerman_go"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ComponentType int

// ComponentClient 仅供 本Server使用
type ComponentClient struct {

	//root
	RegisterService *Register

	//组件实例名称
	Name string
	// Authed 是否通过认证
	Authed bool

	//可能是ipv4 或者Ipv6
	Address string
	//来源端口
	Port string

	//组件类型——内存
	ComponentType ComponentType

	//如果是gateway 填充公网连接地址
	PublicGatewayConnectionInfo workerman_go.ProtocolPublicGatewayConnectionInfo

	FdWs        *websocket.Conn
	DataRWMutex *sync.RWMutex
	Data        map[string]interface{}
	Request     *http.Request
	//TokenStructString
	ClientToken workerman_go.ClientToken

	Ctx       context.Context
	CtxCancel context.CancelFunc
}

// Close 主动关闭接口,会触发InnerOnClose()
func (rc *ComponentClient) Close() {
	//人工通知 关闭事件
	rc.RegisterService.InnerOnClose(rc)
}

// sendWithSignJsonString 内部方法
func (rc *ComponentClient) sendWithSignJsonString(cmd int, v any) error {

	timeByte, err := workerman_go.GenerateSignTimeByte(cmd, v, rc.RegisterService.GatewayWorkerConfig.SignKey, func() time.Duration {
		return time.Second * 10
	})
	if err != nil {
		return err
	}

	//_,ok:=rc.RegisterService.ConnectionListMap[rc.ClientToken.ClientGatewayNum]

	sendErr := rc.FdWs.WriteMessage(websocket.BinaryMessage, timeByte.ToByte())
	if sendErr != nil {
		rc.Close()
		return sendErr
	}
	return nil
}

// CommandToComponentForAllList Broadcast
func (rc *ComponentClient) CommandToComponentForAllList() {
	//rc.RegisterService.
}

// CommandToComponentForAuthRequire 要求发送身份验证
func (rc *ComponentClient) CommandToComponentForAuthRequire() {
	rc.sendWithSignJsonString(workerman_go.CommandComponentAuthRequest, workerman_go.ProtocolRegister{
		//请求授权标志
		Data:   "workerman_go.CommandServiceAuthRequest",
		Authed: strconv.Itoa(0), //告诉组件未授权
	})
}

// Send 发送json数据，但是带有签名校验和时间校验的
func (rc *ComponentClient) Send(data any) error {

	switch data.(type) {
	case workerman_go.ProtocolRegister:
		rc.sendWithSignJsonString(workerman_go.CommandComponentAuthRequest, data)
	case workerman_go.ProtocolRegisterBroadCastComponentGateway:
		rc.sendWithSignJsonString(workerman_go.CommandComponentGatewayList, data)
	case string:
	default:
		return errors.New("conn.Send(Unknown protocol message)")
	}

	return nil
}

func (rc *ComponentClient) GetRemoteIp() string {
	return ""
}

func (rc *ComponentClient) GetRemotePort() string {
	return ""
}

func (rc *ComponentClient) PauseRecv() {

}
func (rc *ComponentClient) ResumeRecv() {
}

func (rc *ComponentClient) Pipe(connection *workerman_go.TcpConnection) {

}

func (rc *ComponentClient) GetClientId() string {
	return rc.ClientToken.GenerateGatewayClientId()
}

func (rc *ComponentClient) GetClientIdInfo() *workerman_go.ClientToken {
	return &rc.ClientToken
}

func (rc *ComponentClient) Get(str string) (interface{}, bool) {
	//读锁，防止读的时候写
	rc.DataRWMutex.RLock()
	defer rc.DataRWMutex.RLock()
	item, ok := rc.Data[str]
	return item, ok
}

func (rc *ComponentClient) Set(str string, v interface{}) {
	//写锁，防止读
	rc.DataRWMutex.Lock()
	defer rc.DataRWMutex.Unlock()
	rc.Data[str] = v
}
