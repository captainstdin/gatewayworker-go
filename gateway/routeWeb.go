package gateway

import (
	"fmt"
	"gatewaywork-go/workerman_go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 对外提供web服务
func (s *Server) listenWeb() {

	const apiPath = "/api/sdk"

	s.gin.POST(apiPath, func(ctx *gin.Context) {

		cmdQuery := ctx.Query("cmd")

		cmd, _ := strconv.Atoi(cmdQuery)

		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errCode": http.StatusBadRequest,
				"errMsg":  err.Error(),
			})
			return
		}

		sdk := &gatewayApi{
			Server: s,
		}

		fmt.Println(data)

		Command := &workerman_go.GenerateComponentSign{
			PackageLen: 0,
			Sign:       [16]byte{},
			TimeStamp:  0,
			Cmd:        int32(cmd),
			Json:       data,
		}
		handleSdkCmd(Command, sdk)

	})

}
