package routes

import (
	"gohub/app/http/controllers/api/v1/auth"

	"github.com/gin-gonic/gin"
)

// 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {
	// 测试一个 v1 的路由组, 所有的 v1 版本的路由都存放到这里
	v1 := r.Group("v1")
	authGroup := v1.Group("/auth")
	{
		suc := new(auth.SignupController)
		// 判断手机号是否被注册
		authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
		// 判断邮箱是否被注册
		authGroup.POST("/signup/email/exist", suc.IsEmailExist)
	}
}
