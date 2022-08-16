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

	// 设置 gin 的运行模式,支持 debug, release, test
	// release 会屏蔽调试信息, 官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息,干扰到我们程序中的 log
	// 故此设置为 release, 有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// 初始化 DB
	bootstrap.SetupDB()

	// 初始化 Redis
	bootstrap.SetupRedis()

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
