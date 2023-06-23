package register

import (
	"encoding/json"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

var coroutine sync.WaitGroup
var Conf = workerman_go.ConfigGatewayWorker{
	RegisterListenAddr:             "0.0.0.0:1238",
	RegisterListenPort:             ":1238",
	TLS:                            false,
	TlsKeyPath:                     "",
	TlsPemPath:                     "",
	RegisterPublicHostForComponent: "127.0.0.1:1238",
	GatewayPublicHostForClient:     "",
	GatewayListenAddr:              "",
	GatewayListenPort:              "",
	SkipVerify:                     false,
	SignKey:                        "da!!bskdhaskld#1238asjiocy89123",
}

func TestServer_Run(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)
	S := NewServer("Test-Register", &Conf)

	go func() {
		defer coroutine.Done()
		err := S.Run()
		if err != nil {
			t.Log(err)
			return
		}
	}()
	//启动测试business
	business(t, 10, false)

	timer := time.NewTimer(1 * time.Second)
	<-timer.C

	coroutine.Wait()
}

func business(t *testing.T, num int, gateway bool) {
	t.Log("开始模拟【Business】")
	wsURL := &url.URL{
		Scheme: "ws",
		Path:   workerman_go.RegisterForComponent,
		Host:   Conf.RegisterPublicHostForComponent,
	}

	// 创建WebSocket配置
	wsConfig := &websocket.Config{
		Location: wsURL,
		Dialer: &net.Dialer{
			Timeout: 10 * time.Second,
		},
		Version: websocket.ProtocolVersionHybi13,
		Origin: &url.URL{
			Scheme: "http",
			//Host: "chat.workerman.net",
		},
	}

	for i := 0; i <= num; i++ {
		coroutine.Add(1)

		go func(i int) {
			defer coroutine.Done()

			// 连接WebSocket服务器
			wsConn, err := websocket.DialConfig(wsConfig)
			if err != nil {
				t.Fatal(err)
			}
			defer wsConn.Close()

			//1kb缓冲区
			buff := make([]byte, 10240)
			for {
				//read
				n, readError := wsConn.Read(buff)
				if readError != nil {
					t.Fatal(readError)
					return
				}

				//解析收到指令
				ReciveData, jsonTimeErr := workerman_go.ParseAndVerifySignJsonTime(buff[:n], Conf.SignKey)
				if jsonTimeErr != nil {
					t.Fatal(jsonTimeErr)
					return
				}

				switch ReciveData.Cmd {
				case workerman_go.CommandComponentAuthRequest:
					var dataRegister workerman_go.ProtocolRegister
					errUnmarshal := json.Unmarshal(ReciveData.Json, &dataRegister)
					if errUnmarshal != nil {
						t.Error("【register】发过来的协议错误，内容为:", string(buff))
						return
					}

					//如果是认证结果
					if dataRegister.Authed == "1" {
						t.Logf("TestRegisterBusiness[%d] 通过，认证结果：%s", i, dataRegister.Data)
						//如果测试没有gateway就不要等了
						if !gateway {
							return
						}
					}
					//如果是请求认证，那就发送
					responseJsonBin, responseJsonBinErr := workerman_go.GenerateSignTimeByte(workerman_go.CommandComponentAuthRequest, workerman_go.ProtocolRegister{
						ComponentType:                       workerman_go.ComponentIdentifiersTypeBusiness,
						Name:                                "register" + strconv.Itoa(i),
						ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
						Data:                                "aaa",
						Authed:                              "1",
					}, Conf.SignKey, func() time.Duration {
						return time.Second * 60
					})
					//生成签名消息失败
					if responseJsonBinErr != nil {
						t.Fatal(responseJsonBinErr)
					}

					_, writeErr := wsConn.Write(responseJsonBin.ToByte())
					if writeErr != nil {
						log.Fatal("发送到【Register】请求认证失败!", writeErr)
						return
					}
				case workerman_go.CommandComponentGatewayList:
					var dataGatewayList workerman_go.ProtocolRegisterBroadCastComponentGateway
					//fmt.Println(string(ReciveData.Json))
					errUnmarshal := json.Unmarshal(ReciveData.Json, &dataGatewayList)
					if errUnmarshal != nil {
						t.Error("【register】发过来的协议错 :workerman_go.CommandComponentGatewayList")
						return
					}
					return
				}

			}

		}(i)

	}
}
