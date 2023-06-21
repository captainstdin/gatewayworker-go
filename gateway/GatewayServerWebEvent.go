package gateway

import (
	"fmt"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	// 对于SDK 或者组件连接上来的地址
	r.GET(workerman_go.GatewayForBusinessWsPath, func(c *gin.Context) {
		clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("【component】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}
		//关闭client ，可能是   business ，当对方主动断开
		defer clientConn.Close()

		//就是business
		ComClient := &ComponentClient{
			root: g,
			ClientId: &workerman_go.ClientToken{
				IPType:            0,
				ClientGatewayIpv4: nil,
				ClientGatewayIpv6: nil,
				ClientGatewayPort: 0,
				ClientGatewayNum:  getUniqueKey(g.ComponentsMap).Uint64(),
			},
			Name:          "sdk|business",
			ComponentType: 0,
			Address:       c.Request.RemoteAddr,
			Port:          0,
			FdWs:          nil,
		}

		for true {
			// 读取客户端发送过来的消息
			_, message, errMsg := clientConn.ReadMessage()

			if errMsg != nil {
				ComClient.onClose(ComClient)
				//退出携程
				break
			}

			g.InnerOnMessage(ComClient, message)

		}

	})

	r.GET("/", func(c *gin.Context) {
		UserClientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"ErrorCode": http.StatusBadRequest,
				"ErrorMsg":  "请升级为websocket协议",
			})
			log.Println("【component】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}
		//关闭client ，可能是   business ，当对方主动断开
		defer UserClientConn.Close()

		//就是business
		userConn := &UserClient{
			root: g,
			ClientId: &workerman_go.ClientToken{
				IPType:            0,
				ClientGatewayIpv4: nil,
				ClientGatewayIpv6: nil,
				ClientGatewayPort: 0,
				ClientGatewayNum:  getUniqueKeyByUserClient(g.ConnectionMap).Uint64(),
			},
			Name:          "sdk|business",
			ComponentType: 0,
			Address:       c.Request.RemoteAddr,
			Port:          0,
			FdWs:          nil,
		}

		g.onConnectForward(userConn)
		for true {
			// 读取客户端发送过来的消息
			_, message, errMsg := UserClientConn.ReadMessage()

			if errMsg != nil {
				g.onCloseForward(userConn)
				//退出携程
				return
			}
			g.onMessageForward(userConn, message)
		}

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
