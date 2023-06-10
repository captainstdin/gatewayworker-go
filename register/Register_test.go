package register

import (
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"net"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

var coroutine sync.WaitGroup
var Conf = &workerman_go.ConfigGatewayWorker{
	RegisterListenAddr:             ":1238",
	RegisterListenPort:             ":1238",
	TLS:                            false,
	TlsKeyPath:                     "",
	TlsPemPath:                     "",
	RegisterPublicHostForComponent: "127.0.0.1:1237",
	GatewayPublicHostForClient:     "",
	GatewayListenAddr:              "",
	GatewayListenPort:              "",
	SkipVerify:                     false,
	SignKey:                        "da!!bskdhaskld#1238asjiocy89123",
}

func TestStartRegister(t *testing.T) {

	service := NewRegister("Business处理器", Conf)

	coroutine.Add(1)
	go func() {
		t.Logf("启动服务器")
		err := service.Run()
		if err != nil {
			t.Fatal(err)
			return
		}
	}()

	testRegisterBusiness(t)
	coroutine.Wait()

}

func testRegisterBusiness(t *testing.T) {

	t.Logf("开始模拟[Business]")
	// 设置WebSocket连接的地址和origin
	wsURL := &url.URL{
		Scheme: "ws",
		Path:   workerman_go.RegisterForBusniessWsPath,
		Host:   Conf.RegisterPublicHostForComponent,
	}

	// 创建WebSocket配置
	wsConfig := &websocket.Config{
		Location: wsURL,
		Dialer: &net.Dialer{
			Timeout: 3 * time.Second,
		},
		Version: websocket.ProtocolVersionHybi13,
		Origin: &url.URL{
			Scheme: "http",
			//Host: "chat.workerman.net",
		},
	}

	for i := range make([]struct{}, 1) {

		coroutine.Add(1)
		go func() {
			defer coroutine.Done()
			t.Log("尝试连接", i)
			// 连接WebSocket服务器
			wsConn, err := websocket.DialConfig(wsConfig)
			t.Log("已连接", i)
			if err != nil {
				t.Fatal(err)
			}
			defer wsConn.Close()

			var buff []byte
			for {
				_, readError := wsConn.Read(buff)
				if readError != nil {
					t.Fatal(readError)
					return
				}

				var jsonCmd workerman_go.ProtocolRegister

				json.Unmarshal(buff, &jsonCmd)

				if jsonCmd.Command == strconv.Itoa(workerman_go.CommandComponentAuthRequest) {

					responseJsno, _ := workerman_go.GenerateSignJsonTime(workerman_go.ProtocolRegister{
						Command:    strconv.Itoa(workerman_go.CommandComponentAuthResponse),
						IsBusiness: "1",
						IsGateway:  "0",
						Data:       "请求认证",
						Authed:     "0",
					}, Conf.SignKey, func() time.Duration {
						return time.Second * 10
					})

					wsConn.Write(responseJsno)
					return
				}

				if jsonCmd.Command == strconv.Itoa(workerman_go.CommandComponentAuthResponse) {

					t.Logf("TestRegisterBusiness[%d] 通过", i)
				}

			}

		}()

	}

}
