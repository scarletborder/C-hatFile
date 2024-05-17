package auth

import (
	"chatFileBackend/utils/storage/db"

	"gorm.io/gorm"
)

func init() {
	db.DBReadConfigWG.Wait()
	err := db.Auth_db.CreateSubDB(grantHandler)
}

func grantHandler(rdb *gorm.DB) error {
	// 只需要保证读写auth_db的user表权限
	return nil
}

// 初始化数据库和加入相关项目

// 缓存随用随初始，没有必要在这里初始化
