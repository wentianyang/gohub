package main

import (
	"fmt"
	"gohub/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 实例
	router := gin.New()

	// 路由初始化
	bootstrap.SetupRoutes(router)

	// 运行服务
	err := router.Run(":8000")
	if err != nil {
		fmt.Println(err.Error())
	}
}
