package routers

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/handlers/search"

	"github.com/gin-gonic/gin"
)

func CreateSearchRouter(engine *gin.Engine) {
	search_router := engine.Group("search")
	search_router.Use(auth_utils.JWT())
	{
		search_router.GET("/search", search.SearchHandler)
		// search_router.POST("/register", auth.RegisterHandler)
	}
}
