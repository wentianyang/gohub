package bootstrap

import (
	"gohub/routes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 路由初始化
func SetupRoutes(router *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleWare(router)

	// 注册 API 路由
	routes.RegisterAPIRoutes(router)

	// 配置 404 路由
	setup404Handler(router)
}

// 注册全局中间件
func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(gin.Logger(), gin.Recovery())
}

// 处理 404 错误
func setup404Handler(router *gin.Engine) {
	router.NoRoute(func(ctx *gin.Context) {
		// 获取请求头信息 accept 信息
		acceptString := ctx.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话,返回页面 404
			ctx.HTML(http.StatusNotFound, "404.html", "页面返回 404")
		} else {
			// 默认返回 JSON
			ctx.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义, 请确认 url 是否正确",
			})
		}
	})
}
