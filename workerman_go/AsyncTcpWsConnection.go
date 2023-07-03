package workerman_go

import (
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/url"
	"time"
)

type AsyncTcpWsConnection struct {
	TcpWsConnection

	OnConnect func(connection *AsyncTcpWsConnection)
}

func (a *AsyncTcpWsConnection) Connect() error {

	parse, err := url.Parse(a.RemoteAddress)
	if err != nil {
		return err
	}

	wsURL := &url.URL{
		Scheme: parse.Scheme,
		Path:   parse.Path,
		Host:   parse.Path,
	}

	// 创建WebSocket配置
	wsConfig := &websocket.Config{
		Location: wsURL,
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		Version: websocket.ProtocolVersionHybi13,
		Origin: &url.URL{
			Scheme: "http",
			//Host: "chat.workerman.net",
		},
	}
	// 连接WebSocket服务器
	wsConn, connErr := websocket.DialConfig(wsConfig)
	if connErr != nil {
		for connErr != nil {
			wsConn, connErr = websocket.DialConfig(wsConfig)
			log.Printf("[%s]无法连接 %s   10秒后重连.. ", a.Name, wsURL.String())
			t := time.NewTicker(time.Second * 10)
			<-t.C
		}
	}

	a.FdAsyncWs = wsConn

	//如果链接成功，记得结束的时候关闭
	defer wsConn.Close()
	//内部的事件通知，阻塞函数
	a.onConnect()

	//100kb缓冲区
	buff := make([]byte, 1024*100)

	//阻塞函数
	for {
		//read
		n, readError := wsConn.Read(buff)
		//onClose
		if readError != nil {
			if a.OnClose != nil {
				a.OnClose(&a.TcpWsConnection)
			}
			a.Close()
			return nil
		}
		//onMessage
		if a.OnMessage != nil {
			a.OnMessage(&a.TcpWsConnection, buff[:n])
		}
	}

}

func (a *AsyncTcpWsConnection) ReConnect() {

}
func (a *AsyncTcpWsConnection) cancelReconnect() {

}

func (a *AsyncTcpWsConnection) onConnect() {
	if a.OnConnect != nil {
		a.OnConnect(a)
	}
}

func NewAsyncTcpWsConnection(remoteAddress string) *AsyncTcpWsConnection {

	ws := TcpWsConnection{
		worker:        nil,
		GatewayIdInfo: nil,
		Name:          "AsyncTcpWsConnection",
		RemoteAddress: remoteAddress,
		Address:       "",
		Port:          0,
		FdWs:          nil,
		OnConnect:     nil,
		OnMessage:     nil,
		OnClose:       nil,
		data:          nil,
		dataLock:      nil,
		Ctx:           nil,
		CtxF:          nil,
	}

	return &AsyncTcpWsConnection{
		TcpWsConnection: ws,
		OnConnect:       nil,
	}

}
