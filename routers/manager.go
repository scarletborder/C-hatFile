package routers

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/handlers/manager"

	"github.com/gin-gonic/gin"
)

func CreateManagerRouter(engine *gin.Engine) {
	manager_router := engine.Group("manager")
	manager_router.Use(auth_utils.JWT())
	{
		manager_router.GET("/self", manager.SelfViewHandler)
		manager_router.DELETE("/del", manager.DeleteHandler)
		// search_router.POST("/register", auth.RegisterHandler)
	}
}
