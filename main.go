package main

import (
	"chatFileBackend/routers"
	"chatFileBackend/utils/global"
	"chatFileBackend/utils/storage/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化完后关闭root数据库,并同步配置文件
	global.Init_count.Wait()
	logrus.Infoln("Init successfully")
	db.CloseRootDSN()

	r := gin.Default()

	routers.CreateAuthRouter(r)
	routers.CreateSearchRouter(r)
	routers.CreateUploadRouter(r)
	routers.CreateDownloadRouter(r)

	r.Run("0.0.0.0:12800")
}
