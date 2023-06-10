package gateway

import (
	"fmt"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10240, //10kb
	WriteBufferSize: 10240,
}

func (g *GatewayServer) RunGinServer(addr string, port string) {
	// 关闭 Gin 的日志输出
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	// 创建一个 Gin 引擎实例
	r := gin.New()
	// 注册一个路由处理函数
	r.GET(workerman_go.GatewayExportForBusinessWsPath, func(c *gin.Context) {
		clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {

			log.Println("Failed to upgrade to WebSocket:", err)
			return
		}
		//关闭client ，可能是   business ，当对方主动断开
		defer clientConn.Close()

		tcpConnect := &GatewayClient{}

		//
		go func() {

		}()
		//todo 写一个定时器30  秒后验证关闭未验证的business
		go func() {
			timer := time.NewTimer(30 * time.Second)
			for true {

				select {
				case <-timer.C:

				}
			}
		}()
		//todo 写一个定时器

		for true {
			// 读取客户端发送过来的消息
			_, message, errMsg := clientConn.ReadMessage()

			if errMsg != nil {
				//退出携程
				break
			}

			g.InnerOnMessage(tcpConnect, message)

		}
		g.InnerOnClose(tcpConnect)

	})
	//
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	if addr == "" {
		addr = ""
	}

	if port == "" {
		port = "8080"
	}
	// 启动服务器
	r.Run(fmt.Sprintf("%s:%s", addr, port))
}
