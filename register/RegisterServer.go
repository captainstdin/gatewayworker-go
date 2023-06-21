package register

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
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
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		w.WriteHeader(http.StatusBadRequest)
		marshal, marshalErr := json.MarshalIndent(map[string]any{"ErrorCode": http.StatusBadRequest, "ErrorMsg": "请升级为websocket协议"}, "", "    ")
		if marshalErr != nil {
			return
		}
		w.Write(marshal)
	},
}

func (register *Register) InnerOnWorkerStart(worker *Register) {

}

// InnerOnConnect 内部处理连接上来的 business或 gateway
func (register *Register) InnerOnConnect(ComponentConn *ComponentClient) {
	//写锁
	register.ConnectionListRWLock.Lock()

	//测试可用uint64编号
	ok := false
	for ok == false {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := register.ConnectionListMap[num.Uint64()]; !exist {
			//设置ClientID信息
			ComponentConn.ClientToken.ClientGatewayNum = num.Uint64()
			//设置列表实例
			register.ConnectionListMap[num.Uint64()] = ComponentConn
			ok = true
		}
	}

	//发送认证请求等待认证,无论是business还是gateway
	ComponentConn.Send(workerman_go.ProtocolRegister{
		Data:   "workerman_go.CommandServiceAuthRequest.first.request",
		Authed: "0", //告诉组件未授权
	})

	//开一个协程，用来倒计时30秒，如果没有认证
	go func(ComponentConn *ComponentClient) {
		timer := time.NewTimer(30 * time.Second)

		select {
		case <-ComponentConn.Ctx.Done():
			//关闭当前协程
			return
		case <-timer.C:
			//关闭连接
			if ComponentConn.Authed == false {
				ComponentConn.Send(workerman_go.ProtocolRegister{
					//请求授权标志
					Data:   "workerman_go.CommandServiceAuthRequest.timeout",
					Authed: "0", //告诉组件未授权
				})
				//人工关闭，OnClose()会处理这一切的
				ComponentConn.Close()
			}
		}

	}(ComponentConn)

	register.ConnectionListRWLock.Unlock()
}

func (register *Register) InnerOnMessage(ComponentConn *ComponentClient, msg []byte) {

	//解析了一次json为map
	CmdData, err := workerman_go.ParseAndVerifySignJsonTime(msg, register.GatewayWorkerConfig.SignKey)
	//不是组件的签名json协议
	if err != nil {
		fmt.Println(err)
		//发送警告日志
		ComponentConn.Send(workerman_go.ProtocolRegister{
			//请求授权标志
			Data:   "workerman_go.CommandServiceAuthRequest.error",
			Authed: "0", //告诉组件未授权
		})
		return
	}

	//解析指令
	switch CmdData.Cmd {
	//认证回应指令
	case workerman_go.CommandComponentAuthRequest:

		var ProtocolRegister workerman_go.ProtocolRegister
		json.Unmarshal(CmdData.Json, &ProtocolRegister)
		//上锁
		ComponentConn.Authed = true
		//设置名字
		ComponentConn.Name = ProtocolRegister.Name
		switch ProtocolRegister.ComponentType {

		case workerman_go.ComponentIdentifiersTypeGateway:
			//设置内存中的类型
			ComponentConn.ComponentType = workerman_go.ComponentIdentifiersTypeGateway
			//gateway 记录公网连接信息
			ComponentConn.PublicGatewayConnectionInfo = ProtocolRegister.ProtocolPublicGatewayConnectionInfo
		case workerman_go.ComponentIdentifiersTypeBusiness:
			//设置内存中的类型
			ComponentConn.ComponentType = workerman_go.ComponentIdentifiersTypeBusiness
			//business 触发广播
			register.BroadcastOnBusinessConnected()
		}

		log.Println("新组件连接：", ComponentConn.ClientToken.ClientGatewayNum)

	case workerman_go.CommandComponentHeartbeat:
		ComponentConn.Set(workerman_go.ComponentLastHeartbeat, strconv.Itoa(int(time.Now().Unix())))
	}

}

// InnerOnClose 当检测到离线时,启动内置回调，删除list中对应的Uint64 map
func (register *Register) InnerOnClose(conn *ComponentClient) {
	register.ConnectionListRWLock.Lock()

	//通知当前用户协程关闭
	conn.CtxCancel()

	//关闭尝试再次关闭conn
	conn.FdWs.Close()

	//删除delete,先判断下，有没有被其他二次删除
	v, ok := register.ConnectionListMap[uint64(conn.ClientToken.ClientGatewayNum)]
	if ok {
		delete(register.ConnectionListMap, uint64(v.ClientToken.ClientGatewayNum))
	}
	register.ConnectionListRWLock.Unlock()

}

// BroadcastOnBusinessConnected 每当新的Business连接：广播给处理器，有关gateway的信息，
func (register *Register) BroadcastOnBusinessConnected() {

	register.ConnectionListRWLock.RLock()
	defer register.ConnectionListRWLock.RUnlock()

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
			Msg:         "authed ! give you gatewayList[]",
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
			//http访问或者非ws
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		//写入服务器，当前的wsConn
		registerClientConn := &ComponentClient{
			RegisterService: register,
			Address:         request.RemoteAddr,
			Port:            "",
			FdWs:            conn,
			DataRWMutex:     &sync.RWMutex{},
			Data:            nil,
			Request:         request,
			Ctx:             ctx,
			CtxCancel:       cancel,
		}

		register.InnerOnConnect(registerClientConn)

		if register.OnConnect != nil {
			register.OnConnect(registerClientConn)
		}

		// 处理 WebSocket 消息
		for {
			_, message, msgError := conn.ReadMessage()

			if msgError != nil {
				fmt.Println("msgError:", msgError)
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
