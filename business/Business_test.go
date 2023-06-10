package business

import (
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/url"
	"testing"
	"time"
)

func TestBusiness_InnerOnWorkerStart(t *testing.T) {
	// 设置WebSocket连接的地址和origin
	wsURL := &url.URL{
		Scheme: "ws",
		Host:   "chat.workerman.net:7272",
	}
	// 创建WebSocket配置
	wsConfig := &websocket.Config{
		Location: wsURL,
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		Origin: &url.URL{Scheme: "http", Host: "chat.workerman.net"},
	}

	// 连接WebSocket服务器
	wsConn, err := websocket.DialConfig(wsConfig)
	if err != nil {
		log.Fatalln(err)
	}
	// 发送和接收数据
	// ...

	// 关闭WebSocket连接
	wsConn.Close()
}
