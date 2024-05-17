package main

import (
	"chatFileBackend/routers"
	"chatFileBackend/utils/storage/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化完后关闭root数据库,并同步配置文件
	db.CloseRootDSN()

	r := gin.Default()

	routers.CreateAuthRouter(r)

	r.Run("0.0.0.0:12800")
}
