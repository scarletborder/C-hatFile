package auth

import (
	"chatFileBackend/models"
	cached "chatFileBackend/utils/storage/cache"
	"chatFileBackend/utils/storage/db"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	go db.Auth_db.CreateSubDB(grantHandler)
	cached.StartDBSync(cached.TypeAuthCache, syncHandler, 5*time.Minute)
}

func syncHandler(chunk []interface{}) error {
	var users []models.User
	for _, v := range chunk {
		if u, ok := v.(models.User); ok {
			users = append(users, u)
		}
	}
	adb := db.Auth_db.GetDB()
	adb.AutoMigrate(&models.User{})
	err := adb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "enc_password"}),
	}).Create(&users).Error
	return err
}

func grantHandler(rdb *gorm.DB) error {
	// 只需要保证读写auth_db的user表权限
	res := rdb.Exec(fmt.Sprintf(
		"GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%'",
		db.Auth_db.Cfg.DB_name, db.GetDBUserName(db.Auth_db.Cfg.DB_name)))
	return res.Error
}

// 初始化数据库和加入相关项目

// 缓存随用随初始，没有必要在这里初始化
