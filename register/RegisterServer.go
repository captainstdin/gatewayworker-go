package register

import (
	"crypto/rand"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Register struct {
	ListenAddr    string
	OnWorkerStart func(Worker *Register)

	OnConnect func(conn workerman_go.TcpConnection)
	OnMessage func(Worker workerman_go.TcpConnection, msg []byte)
	OnClose   func(Worker workerman_go.TcpConnection)

	TLS    bool
	TlsKey string
	TlsPem string

	ConnectionList map[uint64]*workerman_go.TcpConnection
	//读写锁
	ConnectionListRWLock *sync.RWMutex

	Name string

	GatewayWorkerConfig *workerman_go.ConfigGatewayWorker
}

// 创建一个新的 WebSocket 升级器
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (register *Register) InnerOnWorkerStart(worker workerman_go.Worker) {

}

// InnerOnConnect 内部处理连接上来的 business或 gateway
func (register *Register) InnerOnConnect(TcpConnection workerman_go.TcpConnection) {
	//写锁
	register.ConnectionListRWLock.Lock()
	defer register.ConnectionListRWLock.Unlock()
	//测试可用uint64编号
	ok := false
	for ok == false {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := register.ConnectionList[num.Uint64()]; !exist {
			//设置ClientID信息
			TcpConnection.GetClientIdInfo().ClientGatewayNum = workerman_go.GatewayNum(num.Uint64())
			//设置列表实例
			register.ConnectionList[num.Uint64()] = &TcpConnection
			ok = true
		}
	}

	//发送认证请求等待认证,无论是business还是gateway
	TcpConnection.Send(workerman_go.ProtocolRegister{
		Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
		Data:    "workerman_go.CommandServiceAuthRequest.first.request",
		Authed:  "0", //告诉组件未授权
	})

	//开一个协程，用来倒计时30秒，如果没有认证
	go func(TcpConnect workerman_go.TcpConnection) {
		timer := time.NewTimer(30 * time.Second)
		<-timer.C

		item, exist := TcpConnect.Get(workerman_go.ComponentAuthed)
		if exist == false || item.(Authed) == false {
			TcpConnect.Send(workerman_go.ProtocolRegister{
				//请求授权标志
				Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
				Data:    "workerman_go.CommandServiceAuthRequest.timeout",
				Authed:  "0", //告诉组件未授权
			})
			//关闭
			TcpConnect.Close()
		}
	}(TcpConnection)

	//todo 30秒后踢掉未认证的service
}

func (register *Register) InnerOnMessage(TcpConnection workerman_go.TcpConnection, msg []byte) {

	//解析了一次json为map
	MapData, err := workerman_go.ParseAndVerifySignJsonTime(string(msg), register.GatewayWorkerConfig.SignKey)
	//非法协议
	if err != nil {
		TcpConnection.Send(workerman_go.ProtocolRegister{
			//请求授权标志
			Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
			Data:    "workerman_go.CommandServiceAuthRequest",
			Authed:  "0", //告诉组件未授权
		})
		return
	}

	//解析指令
	commandType, commandTypeOk := MapData[workerman_go.ProtocolCommandName]
	//非法指令
	if commandTypeOk == false {
		return
	}
	switch commandType {
	case strconv.Itoa(workerman_go.CommandComponentAuthResponse):

		var CommandMsg workerman_go.ProtocolRegister
		json.Unmarshal(msg, &CommandMsg)
		//上锁
		register.ConnectionListRWLock.Lock()

		//验证默认通过
		TcpConnection.Set(workerman_go.ComponentAuthed, true)

		if CommandMsg.IsBusiness == "1" {
			TcpConnection.Set(workerman_go.ComponentType, workerman_go.ComponentTypeBusiness)
			//todo 处理器则记录到MAP表，并且广播to Gateway
			//处理器则记录到MAP表，并且广播to Gateway
		}
		if CommandMsg.IsGateway == "1" {
			TcpConnection.Set(workerman_go.ComponentType, workerman_go.ComponentTypeGateway)
			//todo 广播则记录到MAP表（？真必要吗），广播 Business
			//广播则记录到MAP表（？真必要吗），广播 Business
		}

		//放锁
		register.ConnectionListRWLock.Unlock()
		//发信息，告诉组件认证通过
		TcpConnection.Send(workerman_go.ProtocolRegister{
			//请求授权标志
			Command: strconv.Itoa(workerman_go.CommandComponentAuthResponse),
			Data:    "workerman_go.CommandComponentAuthResponse.passed",
			Authed:  "1", //告诉组件已授权
		})

	case strconv.Itoa(workerman_go.CommandComponentHeartbeat):
		TcpConnection.Set(workerman_go.ComponentLastHeartbeat, strconv.Itoa(int(time.Now().Unix())))

	}

}

// InnerOnClose 当检测到离线时,启动内置回调，删除list中对应的Uint64 map
func (register *Register) InnerOnClose(conn workerman_go.TcpConnection) {
	register.ConnectionListRWLock.Lock()
	defer register.ConnectionListRWLock.Unlock()
	delete(register.ConnectionList, uint64(conn.GetClientIdInfo().ClientGatewayNum))
}

func (register *Register) Run() error {

	if register.OnWorkerStart != nil {
		register.OnWorkerStart(register)
	}

	handleServer := http.NewServeMux()
	handleServer.HandleFunc(workerman_go.RegisterForBusniessWsPath, func(response http.ResponseWriter, request *http.Request) {
		// 升级 HTTP 连接为 WebSocket 连接
		conn, err := upgrader.Upgrade(response, request, nil)
		if err != nil {
			log.Println("Upgrade Err:", err)
			return
		}
		defer conn.Close()

		//写入服务器，当前的wsConn
		registerClientConn := &RegisterClient{
			RegisterService: register,
			Address:         request.RemoteAddr,
			Port:            "",
			FdWs:            conn,
			DataRWMutex:     &sync.RWMutex{},
			Data:            nil,
			Request:         request,
		}

		register.InnerOnConnect(registerClientConn)

		if register.OnConnect != nil {
			register.OnConnect(registerClientConn)
		}

		// 处理 WebSocket 消息
		for {

			_, message, msgError := conn.ReadMessage()

			if msgError != nil {
				register.InnerOnClose(registerClientConn)
				if register.OnClose != nil {
					register.OnClose(registerClientConn)
				}
				break
			}

			register.InnerOnMessage(registerClientConn, message)
			if register.OnMessage != nil {
				register.OnMessage(registerClientConn, message)
			}
		}
	})

	// 启动 HTTP 服务器
	//addr := ":8080"
	log.Printf("[%s] Starting  server at -> %s ;Listening...", register.Name, register.ListenAddr)

	server := &http.Server{
		Addr:    register.ListenAddr,
		Handler: handleServer,
	}
	var startError error
	if register.TLS {
		startError = server.ListenAndServeTLS("server.crt", "server.key")
	} else {
		startError = server.ListenAndServe()
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

func NewRegister(name string, config *workerman_go.ConfigGatewayWorker) *Register {
	if name == "" {
		name = "Business"
	}

	return &Register{
		Name:                 name,
		ListenAddr:           config.RegisterListenAddr,
		TLS:                  false,
		ConnectionList:       make(map[uint64]*workerman_go.TcpConnection, 0),
		ConnectionListRWLock: &sync.RWMutex{},
		GatewayWorkerConfig:  config,
	}
}
