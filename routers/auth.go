package routers

import (
	"chatFileBackend/handlers/auth" // 引入Gin框架

	"github.com/gin-gonic/gin"
)

func CreateAuthRouter(engine *gin.Engine) {
	auth_router := engine.Group("auth")
	{
		auth_router.POST("/login", auth.LoginHandler)
	}
}
