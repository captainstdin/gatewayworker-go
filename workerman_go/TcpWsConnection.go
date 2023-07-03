package workerman_go

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	websocket2 "golang.org/x/net/websocket"
	"net"
	"strconv"
	"sync"
)

// TcpWsConnection 最基本的tcp接口结构体
type TcpWsConnection struct {
	RequestCtx *gin.Context
	worker     *Worker
	//包含一些基础连接内容：Ip地址和ip类型和fd序号
	GatewayIdInfo *GatewayIdInfo

	//组件名称
	Name string

	//用户的地址
	RemoteAddress string
	//连接地址
	Address string
	Port    uint64

	FdWs *websocket.Conn

	FdAsyncWs *websocket2.Conn

	OnConnect func(connection *TcpWsConnection)

	OnMessage func(connection *TcpWsConnection, buff []byte)

	OnClose func(connection *TcpWsConnection)

	Data map[string]string

	dataLock *sync.RWMutex

	Ctx  context.Context
	CtxF context.CancelFunc
}

// Close 人工关闭连接
func (t *TcpWsConnection) Close() {
	t.FdWs.Close()
}

func (t *TcpWsConnection) Send(data interface{}) error {

	switch data.(type) {
	case string:
		err := t.FdWs.WriteMessage(websocket.TextMessage, []byte(data.(string)))
		if err != nil {
			return err
		}
	case []byte:
		err := t.FdWs.WriteMessage(websocket.BinaryMessage, data.([]byte))
		if err != nil {
			return err
		}
	default:
		return errors.New("unkonw send Data type")
	}
	return nil
}

func parseIPPort(address string) (string, string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", "", err
	}

	return host, port, err

}

// GetRemoteIp 获取远程地址
func (t *TcpWsConnection) GetRemoteIp() (net.IP, error) {
	host, _, err := parseIPPort(t.RemoteAddress)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.New("无效的ip地址")
	}
	return ip, nil
}

// GetRemotePort 获取uint16端口
func (t *TcpWsConnection) GetRemotePort() (uint16, error) {
	_, port, err := parseIPPort(t.RemoteAddress)
	if err != nil {
		return 0, nil
	}
	portUint64, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, errors.New("无效的端口")
	}
	return uint16(portUint64), nil
}

func (t *TcpWsConnection) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (t *TcpWsConnection) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (t *TcpWsConnection) Pipe(connection *TcpWsConnection) {
	//TODO implement me

	t.OnMessage = func(connection *TcpWsConnection, buff []byte) {
		connection.Send(buff)
	}

	t.OnClose = func(connection *TcpWsConnection) {
		connection.Close()
	}

	//暂停和恢复尚未实现
}

func (t *TcpWsConnection) GetClientId() string {
	return t.GatewayIdInfo.GenerateGatewayClientId()
}

func (t *TcpWsConnection) GetClientIdInfo() *GatewayIdInfo {
	//TODO implement me
	return t.GatewayIdInfo
}

func (t *TcpWsConnection) GetRemoteAddress() string {
	return t.RemoteAddress
}

// Get 当心读锁,排写锁
func (t *TcpWsConnection) Get(str string) (string, bool) {
	t.dataLock.RLock()
	defer t.dataLock.RUnlock()
	i, ok := t.Data[str]
	return i, ok
}

func (t *TcpWsConnection) Set(str string, v string) {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()
	t.Data[str] = v
}

func (t *TcpWsConnection) Worker() *Worker {

	return t.worker
}

// TcpWsConnection
func (t *TcpWsConnection) TcpWsConnection() *TcpWsConnection {
	return t
}
