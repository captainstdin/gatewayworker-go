package sdk

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置 Gin 模式为线上模式
	gin.SetMode(gin.ReleaseMode)

	// 创建 Gin 路由
	router := gin.New()

	// 添加路由处理函数
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
		c.Request.Context()
	})

	// 启动服务器
	router.Run(":8080")
}
