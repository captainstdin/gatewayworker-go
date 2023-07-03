package workerman_go

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// TimeOutSecond 默认tcp超时时间
var TimeOutSecond = 30

type Worker struct {
	Connections     map[uint64]InterfaceConnection
	ConnectionsLock *sync.RWMutex
	//only Websocket
	ListenAddress string
	//
	ListenPath string

	Name string

	Tls bool

	TlsPem string
	TlsKey string

	OnWorkerStart func(worker *Worker)
	OnConnect     func(conntion InterfaceConnection)
	OnMessage     func(connection InterfaceConnection, buff []byte)

	OnClose func(connection InterfaceConnection)

	Ctx  context.Context
	CtxF context.CancelFunc

	Config *ConfigGatewayWorker
}

func (w *Worker) Run() error {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  10240, //10kb
		WriteBufferSize: 10240, //10kb
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	route := gin.Default()

	if w.ListenPath == "" {
		w.ListenPath = "/"
	}

	route.GET(w.ListenPath, func(ctx *gin.Context) {

		clientConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("【client】connect   Failed to upgrade to WebSocket:", err)
			return
		}

		w.ConnectionsLock.Lock()

		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())
		uint64Value := genPrimaryKeyUint64(w.Connections)
		Connection := &TcpWsConnection{
			RequestCtx: ctx,
			worker:     w,
			Ctx:        TcpWsCtx,
			CtxF:       TcpWsCancel,
			GatewayIdInfo: &GatewayIdInfo{
				//IPType:            0,
				//ClientGatewayIpv4: nil,
				//ClientGatewayIpv6: nil,
				//ClientGatewayPort: 0,
				ClientGatewayAddr: "",
				ClientGatewayNum:  uint64Value,
			},
			Name:          "default",
			RemoteAddress: ctx.RemoteIP(),
			Address:       "",
			Port:          0,
			FdWs:          clientConn,
			OnConnect:     nil,
			OnMessage:     nil,
			OnClose:       nil,
			Data:          map[string]string{},
			dataLock:      &sync.RWMutex{},
		}

		w.Connections[uint64Value] = Connection
		w.ConnectionsLock.Unlock()

		//注意这里应该是个阻塞函数,不然当前连接就defer了
		w.onConnect(Connection)

		//关闭client ，可能是   business ，当对方主动断开
		defer clientConn.Close()

	})

	w.onWorkerStart(w)

	var err error
	if w.Tls {
		err = route.RunTLS(w.ListenAddress, w.TlsPem, w.TlsKey)
	} else {
		err = route.Run(w.ListenAddress)
	}

	return err
}

func (w *Worker) onWorkerStart(worker InterfaceWorker) {
	if w.OnWorkerStart != nil {
		w.OnWorkerStart(w)
	}
}

func (w *Worker) onConnect(connection InterfaceConnection) {
	if w.OnConnect != nil {
		w.OnConnect(connection)
	}
	//这里是一个block函数，

	go func(connID uint64) {
		timeTick := time.NewTicker(time.Duration(TimeOutSecond) * time.Second)
		for {
			select {
			case <-timeTick.C:
				//踢人
				w.ConnectionsLock.RLock()
				defer w.ConnectionsLock.RUnlock()
				conn, ok := w.Connections[connID]
				if !ok {
					return
				}
				conn.TcpWsConnection().CtxF()
			}
		}
	}(connection.GetClientIdInfo().ClientGatewayNum)

	fd := connection.TcpWsConnection().FdWs
	for {
		select {
		case <-connection.TcpWsConnection().Ctx.Done():
			connection.Close()
			//协程被关闭
			return
		default:
			_, msg, err := fd.ReadMessage()
			if err != nil {
				connection.TcpWsConnection().CtxF()
				continue
			}

			w.onMessage(connection, msg)
		}

	}

}

func (w *Worker) onMessage(connection InterfaceConnection, msg []byte) {

	if w.OnMessage != nil {
		w.OnMessage(connection, msg)
	}
}

func (w *Worker) onClose(connection InterfaceConnection) {

	if w.OnClose != nil {
		w.OnClose(connection)
	}

	connection.Close()

	w.ConnectionsLock.Lock()
	defer w.ConnectionsLock.Unlock()

	//删除
	delete(w.Connections, connection.GetClientIdInfo().ClientGatewayNum)

}
