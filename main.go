package main

import (
	"chatFileBackend/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routers.CreateAuthRouter(r)

	r.Run("0.0.0.0:12800")
}
