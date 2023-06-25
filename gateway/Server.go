package gateway

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	workerman_go.Worker

	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

	//已连接的Gateway
	ConnectedBusinessLock *sync.RWMutex
	ConnectedBusinessMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

	//连接
	Connections     map[uint64]workerman_go.InterfaceConnection
	ConnectionsLock *sync.RWMutex
}

func genPrimaryKeyUint64(mapData map[uint64]workerman_go.InterfaceConnection) uint64 {

	for {
		num, err := rand.Int(rand.Reader, big.NewInt(1<<63-1))
		if err != nil {
			panic(err)
		}
		if _, exist := mapData[num.Uint64()]; !exist {
			mapData[num.Uint64()] = nil
			return num.Uint64()
		}
	}
	return 0
}

func (s *Server) sendSignData(data any, conn workerman_go.InterfaceConnection) {

	var CMDInt int
	switch data.(type) {
	case workerman_go.ProtocolRegister:
		CMDInt = workerman_go.CommandComponentAuthRequest
	case workerman_go.ProtocolPublicGatewayConnectionInfo:

	}

	timeByte, err := workerman_go.GenerateSignTimeByte(CMDInt, data, s.Config.SignKey, func() time.Duration {
		return time.Second * 120
	})
	if err != nil {
		return
	}
	conn.Send(timeByte.ToByte())
}

func (s *Server) Run() error {

	business := workerman_go.NewAsyncTcpWsConnection(s.Config.RegisterPublicHostForComponent)

	business.OnConnect = func(connection *workerman_go.AsyncTcpWsConnection) {
		s.ConnectedBusinessLock.Lock()
		s.ConnectedBusinessMap[connection.GetRemoteAddress()] = connection
		s.ConnectedBusinessLock.Unlock()
	}

	business.OnMessage = func(connection *workerman_go.TcpWsConnection, buff []byte) {
		//business的指令
		Data, parse := workerman_go.ParseAndVerifySignJsonTime(buff, s.Config.SignKey)
		if parse != nil {
			return
		}
		switch Data.Cmd {
		case workerman_go.CommandComponentAuthRequest:
			var reg workerman_go.ProtocolRegister
			err := json.Unmarshal(buff, &reg)
			if err != nil {
				return
			}
			if reg.Authed == "1" {
				s.ConnectedRegisterLock.RLock()
				s.ConnectedRegisterMap[connection.GetRemoteAddress()] = connection
				s.ConnectedRegisterLock.Unlock()
				return
			}
			s.sendSignData(workerman_go.ProtocolRegister{}, connection)
		}
	}

	business.OnClose = func(connection *workerman_go.TcpWsConnection) {
		//todo 重连business
	}

	go func() {
		err := business.Connect()
		if err != nil {
			business.Connect()
		}
	}()

	var upgraderWs = websocket.Upgrader{
		ReadBufferSize:  10240, //10kb
		WriteBufferSize: 10240, //10kb
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	r := gin.Default()

	//SDk或者Business处理器
	r.GET(workerman_go.GatewayForBusinessWsPath, func(ctx *gin.Context) {
		clientConn, err := upgraderWs.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("【business】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}

		s.ConnectedBusinessLock.Lock()
		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())

		ConnectionBusiness := &workerman_go.TcpWsConnection{
			RequestHttp: ctx,
			Ctx:         TcpWsCtx,
			CtxF:        TcpWsCancel,
			ClientToken: &workerman_go.ClientToken{},
			Name:        "default",
			Address:     "",
			Port:        0,
			FdWs:        clientConn,
			OnConnect:   nil,
			OnMessage:   nil,
			OnClose:     nil,
		}

		s.ConnectedBusinessMap[ctx.Request.RemoteAddr] = ConnectionBusiness
		s.ConnectedBusinessLock.Unlock()
	})

	//普通用户
	r.GET(workerman_go.GatewayForUserPath, func(ctx *gin.Context) {
		clientConn, err := upgraderWs.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Println("【clientUser】connect gateway Failed to upgrade to WebSocket:", err)
			return
		}

		s.ConnectionsLock.Lock()
		TcpWsCtx, TcpWsCancel := context.WithCancel(context.Background())

		ConnectionUser := &workerman_go.TcpWsConnection{
			RequestHttp: ctx,
			Ctx:         TcpWsCtx,
			CtxF:        TcpWsCancel,
			ClientToken: &workerman_go.ClientToken{},
			Name:        "default",
			Address:     "",
			Port:        0,
			FdWs:        clientConn,
			OnConnect:   nil,
			OnMessage:   nil,
			OnClose:     nil,
		}

		s.Connections[genPrimaryKeyUint64(s.Connections)] = ConnectionUser
		s.ConnectionsLock.Unlock()

	})

	var err error
	if s.Config.TLS {
		err = r.RunTLS(s.ListenAddress, s.Config.TlsPemPath, s.Config.TlsKeyPath)
	} else {
		err = r.Run(s.ListenAddress)
	}

	if err != nil {
		return err
	}

	return nil
}

func NewGatewayServer(name string, conf *workerman_go.ConfigGatewayWorker) *Server {

	ctx, cf := context.WithCancel(context.Background())
	Worker := workerman_go.Worker{
		Connections:     map[uint64]workerman_go.InterfaceConnection{},
		ConnectionsLock: &sync.RWMutex{},
		ListenAddress:   conf.RegisterListenAddr,
		ListenPath:      "/",
		Name:            name,
		Tls:             false,
		TlsPem:          "",
		TlsKey:          "",
		OnWorkerStart:   nil,
		OnConnect:       nil,
		OnMessage:       nil,
		OnClose:         nil,
		Ctx:             ctx,
		CtxF:            cf,
		Config:          conf,
	}

	server := &Server{
		Worker:                Worker,
		ConnectedBusinessLock: &sync.RWMutex{},
		ConnectedRegisterLock: &sync.RWMutex{},
		ConnectedRegisterMap:  make(map[string]workerman_go.InterfaceConnection),
		ConnectedBusinessMap:  make(map[string]workerman_go.InterfaceConnection),
	}
	server.Worker.OnWorkerStart = server.OnWorkerStart
	server.Worker.OnConnect = server.OnConnect
	server.Worker.OnMessage = server.OnMessage
	server.Worker.OnClose = server.OnClose

	server.Run()

	return server
}
