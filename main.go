package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Gin 实例
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 注册一个路由
	r.GET("", func(c *gin.Context) {
		// 以 JSON 的格式响应
		c.JSON(http.StatusOK, gin.H{
			"Hello": "Wolrd!!!",
			"Name":  "杨天文",
		})
	})

	// 处理 404 请求
	r.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 accept 信息
		accept := c.Request.Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			// 如果是 HTML 的话,返回页面 404
			c.HTML(http.StatusNotFound, "404.html", "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义, 请确认 url 是否正确",
			})
		}
	})

	// 运行服务
	r.Run(":8000")
}
