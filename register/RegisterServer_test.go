package register

import (
	"encoding/json"
	"fmt"
	"gatewaywork-go/workerman_go"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var coroutine sync.WaitGroup
var Conf = &workerman_go.ConfigGatewayWorker{
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

func TestStartRegister(t *testing.T) {
	t.Logf("启动【register/etcd】")
	//file, err := os.Create("output.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.SetOutput(file)
	log.SetOutput(nil)
	log.SetOutput(os.Stdout)

	service := NewRegister("Business处理器", Conf)

	coroutine.Add(1)
	go func() {
		err := service.Run()
		coroutine.Done()
		if err != nil {
			t.Fatal(err)
		}
	}()

	testRegisterBusiness(t)

	<-time.After(time.Second * 3)

	coroutine.Wait()
}

func testRegisterBusiness(t *testing.T) {

	t.Logf("开始模拟【Business】")
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
			Timeout: 10 * time.Second,
		},
		Version: websocket.ProtocolVersionHybi13,
		Origin: &url.URL{
			Scheme: "http",
			//Host: "chat.workerman.net",
		},
	}

	for i := range make([]struct{}, 3) {
		coroutine.Add(1)

		go func() {
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
				n, readError := wsConn.Read(buff)
				if readError != nil {
					t.Fatal(readError)
					return
				}

				jsonTime, jsonTimeErr := workerman_go.ParseAndVerifySignJsonTime(buff[:n], Conf.SignKey)
				if jsonTimeErr != nil {
					t.Fatal(jsonTimeErr)
					return
				}

				var dataRegister workerman_go.ProtocolRegister

				errUnmarshal := json.Unmarshal(jsonTime.Json, &dataRegister)
				if errUnmarshal != nil {
					t.Error("【register】发过来的协议错误，内容为:", string(buff))
					return
				}

				fmt.Println(jsonTime.Cmd)

				if jsonTime.Cmd == workerman_go.CommandComponentAuthRequest {

					responseJsonBin, responseJsonBinErr := workerman_go.GenerateSignTimeByte(workerman_go.CommandComponentAuthResponse, workerman_go.ProtocolRegister{
						ComponentType:                       0,
						Name:                                "register" + strconv.Itoa(i),
						ProtocolPublicGatewayConnectionInfo: workerman_go.ProtocolPublicGatewayConnectionInfo{},
						Data:                                "aaa",
						Authed:                              "1",
					}, Conf.SignKey, func() time.Duration {
						return time.Second * 60
					})
					if responseJsonBinErr != nil {
						t.Fatal(responseJsonBinErr)
					}

					//t.Log(responseJsonBin.ToString())

					_, writeErr := wsConn.Write(responseJsonBin.ToByte())
					if writeErr != nil {
						log.Fatal("发送到【Register】请求认证失败!", writeErr)
						return
					}
					return
				}

				if jsonTime.Cmd == workerman_go.CommandComponentGatewayListResponse {

					t.Logf("TestRegisterBusiness[%d] 通过，认证结果：%s", i, string(dataRegister.Authed))
				}

			}

		}()

	}

}
