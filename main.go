package main

import (
	"flag"
	"fmt"
	"gohub/bootstrap"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"

	btsConfig "gohub/config"
)

func init() {
	// 加载 config 包下的配置信息
	btsConfig.Initialize()
}

func main() {

	// 初始化配置信息,依赖命令行 --env 参数
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件, 如 --env=testing 记载的是 .env.testing 文件")
	flag.Parse()
	config.InitConfig(env)

	// 初始化 Logger
	bootstrap.SetupLogger()

	// 初始化 DB
	bootstrap.SetupDB()

	// 初始化 Gin 实例
	router := gin.New()

	// 路由初始化
	bootstrap.SetupRoutes(router)

	// 运行服务
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
