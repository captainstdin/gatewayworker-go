package register

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"workerman_go/workerman_go"
)

type Register struct {
	ListenAddr    string
	OnWorkerStart func(Worker *Register)

	OnConnect func(conn *RegisterClientConnect)
	OnMessage func(Worker *RegisterClientConnect, msg []byte)
	OnClose   func(Worker *RegisterClientConnect)

	TLS    bool
	TlsKey string
	TlsPem string

	ConnectionList []*RegisterClientConnect

	//读写锁
	RWLock *sync.RWMutex
}

type RegisterClientConnect struct {
	Address string
	Port    string
	Fd      *websocket.Conn
	Data    map[string]string
	Request *http.Request
}

func (conn *RegisterClientConnect) SendCommand(v interface{}) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return
	}

	err = conn.Fd.WriteMessage(websocket.TextMessage, marshal)
	if err != nil {
		//close conn
		conn.Fd.Close()
		return
	}

}

// 创建一个新的 WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 内部处理连接上来的 business或 gateway
func (this *Register) _OnConnect(connection *RegisterClientConnect) {
	//写锁
	this.RWLock.Lock()
	this.ConnectionList = append(this.ConnectionList, connection)
	//放锁
	this.RWLock.Unlock()

	//发送认证请求等待认证,无论是business还是gateway
	connection.SendCommand(workerman_go.ProtocolJsonRegister{Command: workerman_go.CommandServiceAuthRequest})

}

func (receiver *RegisterClientConnect) _OnMessage(conn *RegisterClientConnect, msg []byte) {

	var ResponseOfService workerman_go.ProtocolJsonRegister
	json.Unmarshal(msg, &ResponseOfService)

	if ResponseOfService.IsBusiness == 1 {
		//处理器则记录到MAP表，并且广播to Gateway
	}

	if ResponseOfService.IsGateway == 0 {
		//广播则记录到MAP表（？真必要吗），广播 Business
	}
}

func (receiver *RegisterClientConnect) _OnClose(conn *RegisterClientConnect) {

}

func (this *Register) Run() error {

	if this.OnWorkerStart != nil {
		this.OnWorkerStart(this)
	}

	handleServer := http.NewServeMux()
	handleServer.HandleFunc(RegisterBusniessWsPath, func(response http.ResponseWriter, request *http.Request) {
		// 升级 HTTP 连接为 WebSocket 连接
		conn, err := upgrader.Upgrade(response, request, nil)
		if err != nil {
			log.Println("Upgrade Err:", err)
			return
		}
		defer conn.Close()
		//写入服务器，当前的wsConn
		registerClientConnection := &RegisterClientConnect{
			Address: request.RemoteAddr,
			Fd:      conn,
			Data:    nil,
			Request: request,
		}

		this._OnConnect(registerClientConnection)
		if this.OnConnect != nil {
			this.OnConnect(registerClientConnection)
		}
		// 处理 WebSocket 消息
		for {
			_, message, msgError := conn.ReadMessage()
			if msgError != nil {
				registerClientConnection._OnClose(registerClientConnection)
				if this.OnClose != nil {
					this.OnClose(registerClientConnection)
				}
				break
			}

			registerClientConnection._OnMessage(registerClientConnection, message)
			if this.OnMessage != nil {
				this.OnMessage(registerClientConnection, message)
			}
			// 回复 WebSocket 消息
			//err = conn.WriteMessage(messageType, message)
			//if err != nil {
			//	log.Println("WriteMessage:", err)
			//	break
			//}
		}
	})

	// 启动 HTTP 服务器
	//addr := ":8080"
	log.Printf("Starting server at %s", this.ListenAddr)

	var startError error
	if this.TLS {
		startError = http.ListenAndServeTLS(this.ListenAddr, "server.crt", "server.key", handleServer)
	} else {
		startError = http.ListenAndServe(this.ListenAddr, handleServer)

	}
	if startError != nil {
		return startError
		log.Fatal("ListenAndServeTLS: ", startError)
	}
	return nil

}
