package routers

import (
	"chatFileBackend/handlers/blogs"

	"github.com/gin-gonic/gin"
)

func CreateBlogsRouter(engine *gin.Engine) {
	blogs_router := engine.Group("blogs")
	{
		blogs_router.GET("", blogs.BlogPreviewHandler)
		blogs_router.GET("/:id", blogs.BlogGetHandler)
	}

}
