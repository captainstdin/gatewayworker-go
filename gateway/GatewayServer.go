package gateway

import (
	"crypto/rand"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/big"
	"net/http"
	"sync"
	"time"
)

var upgraderWs = websocket.Upgrader{
	ReadBufferSize:  10240, //10kb
	WriteBufferSize: 10240, //10kb
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	gin           *gin.Engine
	Name          string
	ListenAddress string

	Config *workerman_go.ConfigGatewayWorker
	//已连接的register
	ConnectedRegisterLock *sync.RWMutex
	ConnectedRegisterMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

	//已连接的Gateway
	ConnectedBusinessLock *sync.RWMutex
	ConnectedBusinessMap  map[string]workerman_go.InterfaceConnection //key是remoteAddress

	//用户连接
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

	//step1 连接register
	go s.connectBusness()

	s.gin = gin.Default()

	//注册组件监听
	s.listenComponent()
	//监听用户
	s.listenUser()

	var err error
	if s.Config.TLS {
		err = s.gin.RunTLS(s.ListenAddress, s.Config.TlsPemPath, s.Config.TlsKeyPath)
	} else {
		err = s.gin.Run(s.ListenAddress)
	}

	if err != nil {
		return err
	}

	return nil
}

func NewGatewayServer(name string, conf *workerman_go.ConfigGatewayWorker) *Server {

	server := &Server{
		Name:                  name,
		Config:                conf,
		ConnectedBusinessLock: &sync.RWMutex{},
		ConnectedRegisterLock: &sync.RWMutex{},
		ConnectedRegisterMap:  make(map[string]workerman_go.InterfaceConnection),
		ConnectedBusinessMap:  make(map[string]workerman_go.InterfaceConnection),
	}
	return server
}
