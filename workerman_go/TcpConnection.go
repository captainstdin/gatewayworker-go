package workerman_go

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"

	"net"
	"strconv"
	"sync"
)

// TcpConnection 最基本的tcp接口结构体
type TcpConnection struct {

	//包含一些基础连接内容：Ip地址和ip类型和fd序号
	ClientToken ClientToken

	//组件名称
	Name string

	//用户的地址
	remoteAddress string
	//连接地址
	Address string
	Port    uint64

	FdWs *websocket.Conn

	OnConnect func(connection *TcpConnection)

	OnMessage func(connection *TcpConnection, buff []byte)

	OnClose func(connection *TcpConnection)

	data map[string]interface{}

	dataLock *sync.RWMutex

	Ctx  context.Context
	CtxF context.CancelFunc
}

// Close 人工关闭连接
func (t *TcpConnection) Close() {
	t.FdWs.Close()
}

func (t *TcpConnection) Send(data interface{}) error {

	switch data.(type) {
	case string:
		err := t.FdWs.WriteMessage(websocket.TextMessage, []byte(data.(string)))
		if err != nil {
			return err
		}
	case byte:
		err := t.FdWs.WriteMessage(websocket.BinaryMessage, data.([]byte))
		if err != nil {
			return err
		}

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
func (t *TcpConnection) GetRemoteIp() (net.IP, error) {
	host, _, err := parseIPPort(t.remoteAddress)
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
func (t *TcpConnection) GetRemotePort() (uint16, error) {
	_, port, err := parseIPPort(t.remoteAddress)
	if err != nil {
		return 0, nil
	}
	portUint64, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return 0, errors.New("无效的端口")
	}
	return uint16(portUint64), nil
}

func (t *TcpConnection) PauseRecv() {
	//TODO implement me
	panic("implement me")
}

func (t *TcpConnection) ResumeRecv() {
	//TODO implement me
	panic("implement me")
}

func (t *TcpConnection) Pipe(connection *TcpConnection) {
	//TODO implement me

	t.OnMessage = func(connection *TcpConnection, buff []byte) {
		connection.Send(buff)
	}

	t.OnClose = func(connection *TcpConnection) {
		connection.Close()
	}

	//暂停和恢复尚未实现
}

func (t *TcpConnection) GetClientId() string {

	return t.ClientToken.GenerateGatewayClientId()
}

func (t *TcpConnection) GetClientIdInfo() *ClientToken {
	//TODO implement me
	panic("implement me")
}

// Get 当心读锁,排写锁
func (t *TcpConnection) Get(str string) (interface{}, bool) {
	t.dataLock.RLock()
	defer t.dataLock.RUnlock()
	i, ok := t.data[str]
	return i, ok
}

func (t *TcpConnection) Set(str string, v interface{}) {
	t.dataLock.Lock()
	defer t.dataLock.Unlock()
	t.data[str] = v
}

func (t *TcpConnection) GotCtxWithF() (context.Context, context.CancelFunc) {
	return t.Ctx, t.CtxF
}

func (t *TcpConnection) GotFd() *websocket.Conn {

	return t.FdWs
}
