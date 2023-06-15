package register

import (
	"bytes"
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

	ConnectionListMap map[uint64]*ComponentClient

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

func (register *Register) InnerOnWorkerStart(worker *Register) {

}

// InnerOnConnect 内部处理连接上来的 business或 gateway
func (register *Register) InnerOnConnect(ComponentConn *ComponentClient) {
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
		if _, exist := register.ConnectionListMap[num.Uint64()]; !exist {
			//设置ClientID信息
			ComponentConn.ClientToken.ClientGatewayNum = workerman_go.GatewayNum(num.Uint64())
			//设置列表实例
			register.ConnectionListMap[num.Uint64()] = ComponentConn
			ok = true
		}
	}

	//发送认证请求等待认证,无论是business还是gateway
	ComponentConn.Send(workerman_go.ProtocolRegister{
		Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
		Data:    "workerman_go.CommandServiceAuthRequest.first.request",
		Authed:  "0", //告诉组件未授权
	})

	//开一个协程，用来倒计时30秒，如果没有认证
	go func(ComponentConn *ComponentClient) {
		timer := time.NewTimer(30 * time.Second)
		<-timer.C

		if ComponentConn.Authed {
			ComponentConn.Send(workerman_go.ProtocolRegister{
				//请求授权标志
				Command: strconv.Itoa(workerman_go.CommandComponentAuthRequest),
				Data:    "workerman_go.CommandServiceAuthRequest.timeout",
				Authed:  "0", //告诉组件未授权
			})
			//关闭
			ComponentConn.Close()
		}
	}(ComponentConn)

	//todo 30秒后踢掉未认证的service
}

func (register *Register) InnerOnMessage(ComponentConn *ComponentClient, msg []byte) {

	//解析了一次json为map
	MapData, err := workerman_go.ParseAndVerifySignJsonTime(string(msg), register.GatewayWorkerConfig.SignKey)
	//非法协议
	if err != nil {
		ComponentConn.Send(workerman_go.ProtocolRegister{
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
	//认证回应指令
	case strconv.Itoa(workerman_go.CommandComponentAuthResponse):
		var CommandMsg workerman_go.ProtocolRegister
		json.Unmarshal(msg, &CommandMsg)
		//上锁
		register.ConnectionListRWLock.Lock()

		//此处的的 CommmandMessage已经通过签名校验可信
		//ComponentConn.Set(workerman_go.ComponentIdentifiersAuthed, true)
		ComponentConn.Authed = true
		//设置名字
		ComponentConn.Name = CommandMsg.Name
		switch CommandMsg.ComponentType {

		case workerman_go.ComponentIdentifiersTypeGateway:
			//设置内存中的类型
			ComponentConn.ComponentType = workerman_go.ComponentIdentifiersTypeGateway
			//gateway 记录公网连接信息
			ComponentConn.PublicGatewayConnectionInfo = CommandMsg.ProtocolPublicGatewayConnectionInfo
		case workerman_go.ComponentIdentifiersTypeBusiness:
			//设置内存中的类型
			ComponentConn.ComponentType = workerman_go.ComponentIdentifiersTypeBusiness
			//business 触发广播
			register.BroadcastOnBusinessConnected()
		}

		//放锁
		register.ConnectionListRWLock.Unlock()
		//发信息，告诉组件认证通过
		ComponentConn.Send(workerman_go.ProtocolRegister{
			//请求授权标志
			Command: strconv.Itoa(workerman_go.CommandComponentAuthResponse),
			Data:    "workerman_go.CommandComponentAuthResponse.passed",
			Authed:  "1", //告诉组件已授权
		})

	case strconv.Itoa(workerman_go.CommandComponentHeartbeat):
		ComponentConn.Set(workerman_go.ComponentLastHeartbeat, strconv.Itoa(int(time.Now().Unix())))
	}

}

// InnerOnClose 当检测到离线时,启动内置回调，删除list中对应的Uint64 map
func (register *Register) InnerOnClose(conn *ComponentClient) {
	register.ConnectionListRWLock.Lock()
	defer register.ConnectionListRWLock.Unlock()
	delete(register.ConnectionListMap, uint64(conn.ClientToken.ClientGatewayNum))
}

// BroadcastOnBusinessConnected 每当新的Business连接：广播给处理器，有关gateway的信息，
func (register *Register) BroadcastOnBusinessConnected() {

	GatewayList := make([]workerman_go.ProtocolPublicGatewayConnectionInfo, 0)

	BusinessList := make([]*ComponentClient, 0)
	for _, ComponentItem := range register.ConnectionListMap {
		//只筛选校验通过的
		if !ComponentItem.Authed {
			continue
		}
		//开始筛选组件类型
		switch ComponentItem.ComponentType {
		case workerman_go.ComponentIdentifiersTypeGateway:
			GatewayList = append(GatewayList, ComponentItem.PublicGatewayConnectionInfo)
		case workerman_go.ComponentIdentifiersTypeBusiness:
			BusinessList = append(BusinessList, ComponentItem)
		}
	}

	//channel阻塞式发送给business广播
	for _, BusinessConn := range BusinessList {
		BusinessConn.Send(workerman_go.ProtocolRegisterBroadCastComponentGateway{
			Command:     strconv.Itoa(workerman_go.CommandComponentGatewayListResponse),
			Data:        "workerman_go.CommandComponentGatewayListResponse",
			GatewayList: GatewayList,
		})
	}

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
		registerClientConn := &ComponentClient{
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

	startInfo := bytes.Buffer{}
	startInfo.WriteByte('[')
	startInfo.WriteString(register.Name)
	startInfo.WriteString("] Starting  server at  ->【")
	startInfo.WriteString(register.ListenAddr)
	startInfo.WriteString("】 Listening...")

	//log.Println(startInfo.Bytes())
	log.Println(strconv.Quote(startInfo.String()))

	//addr := ":8080"

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
	//正常exit
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
		ConnectionListMap:    make(map[uint64]*ComponentClient, 0),
		ConnectionListRWLock: &sync.RWMutex{},
		GatewayWorkerConfig:  config,
	}
}
