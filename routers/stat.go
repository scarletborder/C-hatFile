package routers

import (
	"chatFileBackend/handlers/stat"

	"github.com/gin-gonic/gin"
)

func CreateStatRouter(engine *gin.Engine) {
	stat_router := engine.Group("stat")
	{
		stat_router.GET("/rank", stat.StatRankHandler)
	}
}
