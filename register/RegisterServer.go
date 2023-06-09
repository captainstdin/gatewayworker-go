package register

import (
	"crypto/rand"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
	"workerman_go/workerman_go"
)

type Register struct {
	ListenAddr    string
	OnWorkerStart func(Worker *Register)

	OnConnect func(conn *RegisterClient)
	OnMessage func(Worker *RegisterClient, msg []byte)
	OnClose   func(Worker *RegisterClient)

	TLS    bool
	TlsKey string
	TlsPem string

	ConnectionList map[uint64]*RegisterClient
	//读写锁
	ConnectionListRWLock *sync.RWMutex
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
func (this *Register) _OnConnect(connection *RegisterClient) {
	//写锁
	this.ConnectionListRWLock.Lock()
	//放锁
	defer this.ConnectionListRWLock.Unlock()
	//编号最大

	ok := false
	for ok == false {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, ok = this.ConnectionList[num.Uint64()]; !ok {
			connection.ClientToken.ClientGatewayNum = workerman_go.GatewayNum(num.Uint64())
			this.ConnectionList[num.Uint64()] = connection
		}
	}

	//发送认证请求等待认证,无论是business还是gateway
	connection.SendCommand(workerman_go.ProtocolJsonRegister{Command: workerman_go.CommandServiceAuthRequest})

	go func(conn *RegisterClient) {
		timer := time.NewTimer(30 * time.Second)
		<-timer.C

		if conn.Authed == false {
			conn.Close()
		}
	}(connection)

	//todo 30秒后踢掉未认证的service
}

func (register *Register) _OnMessage(conn *RegisterClient, msg []byte) {

	var ResponseOfService workerman_go.ProtocolJsonRegister
	err := json.Unmarshal(msg, &ResponseOfService)
	if err != nil {
		return
	}

	register.ConnectionListRWLock.Lock()
	defer register.ConnectionListRWLock.Unlock()

	switch ResponseOfService.Command {
	case workerman_go.CommandServiceAuthResponse:
		if ResponseOfService.IsBusiness == 1 {
			register.ConnectionList[uint64(conn.ClientToken.ClientGatewayNum)].ServiceType = workerman_go.ServiceTypeBusiness
			//todo
			//处理器则记录到MAP表，并且广播to Gateway
		}

		if ResponseOfService.IsGateway == 0 {
			register.ConnectionList[uint64(conn.ClientToken.ClientGatewayNum)].ServiceType = workerman_go.ServiceTypeGateway
			//todo
			//广播则记录到MAP表（？真必要吗），广播 Business
		}
	}

}

// 当检测到离线时,启动内置回调，删除list中对应的Uint64
func (rc *Register) _OnClose(conn *RegisterClient) {
	rc.ConnectionListRWLock.Lock()
	defer rc.ConnectionListRWLock.Unlock()
	delete(rc.ConnectionList, uint64(conn.ClientToken.ClientGatewayNum))
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
		registerClientConn := &RegisterClient{
			RegisterService: this,
			Address:         request.RemoteAddr,
			FdWs:            conn,
			Data:            nil,
			Request:         request,
		}

		this._OnConnect(registerClientConn)
		if this.OnConnect != nil {
			this.OnConnect(registerClientConn)
		}
		// 处理 WebSocket 消息
		for {
			_, message, msgError := conn.ReadMessage()
			if msgError != nil {
				this._OnClose(registerClientConn)
				if this.OnClose != nil {
					this.OnClose(registerClientConn)
				}
				break
			}

			this._OnMessage(registerClientConn, message)
			if this.OnMessage != nil {
				this.OnMessage(registerClientConn, message)
			}
		}
	})

	// 启动 HTTP 服务器
	//addr := ":8080"
	log.Printf("[Register] Starting  server at %s ;Listening...", this.ListenAddr)

	var startError error
	if this.TLS {
		startError = http.ListenAndServeTLS(this.ListenAddr, "server.crt", "server.key", handleServer)
	} else {
		startError = http.ListenAndServe(this.ListenAddr, handleServer)
	}
	if startError != nil {
		return startError
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-make(chan struct{})
	}()
	wg.Wait()
	return nil
}

func NewRegister() *Register {
	return &Register{
		ListenAddr:           ":1237",
		TLS:                  false,
		ConnectionList:       make(map[uint64]*RegisterClient, 0),
		ConnectionListRWLock: &sync.RWMutex{},
	}
}
