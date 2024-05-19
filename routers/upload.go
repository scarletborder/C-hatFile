package routers

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/handlers/upload"

	"github.com/gin-gonic/gin"
)

func CreateUploadRouter(engine *gin.Engine) {
	auth_router := engine.Group("upload")
	auth_router.Use(auth_utils.JWT())
	{
		auth_router.POST("/file", upload.UploadHandler)
	}
}
