package routers

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/handlers/download"

	"github.com/gin-gonic/gin"
)

func CreateDownloadRouter(engine *gin.Engine) {
	search_router := engine.Group("download")
	search_router.Use(auth_utils.JWT())
	{
		search_router.GET("/file", download.DownloadHandler)
		// search_router.POST("/register", auth.RegisterHandler)
	}
}
