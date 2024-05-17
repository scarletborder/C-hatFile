package routers

import (
	"chatFileBackend/handlers/upload"

	"github.com/gin-gonic/gin"
)

func CreateUploadRouter(engine *gin.Engine) {
	auth_router := engine.Group("upload")
	{
		auth_router.POST("/file", upload.UploadHandler)
	}
}
