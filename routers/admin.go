package routers

import "github.com/gin-gonic/gin"

// 管理后台

func CreateAdminRouter(engine *gin.Engine) {
	engine.Group("admin")
}
