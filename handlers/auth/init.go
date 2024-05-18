package auth

import (
	"chatFileBackend/models"
	cached "chatFileBackend/utils/storage/cache"
	"chatFileBackend/utils/storage/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	go func() {
		db.DBReadConfigWG.Wait()
		db.Auth_db.Cfg = &db.DBCfg.SubDBCfgs.AuthDB
		db.Auth_db.CreateSubDB(grantHandler)
		cached.StartDBSync(cached.TypeAuthCache, syncHandler, cached.Sync_interval)
	}()
}

func syncHandler(chunk []interface{}) error {
	if len(chunk) == 0 {
		return nil
	}
	var users []models.User
	for _, v := range chunk {
		if u, ok := v.(*models.User); ok {
			users = append(users, *u)
		}
	}

	if len(users) == 0 {
		return errors.New("something error happened in sync auth, for chunk is not empty first, but after asserts it becomes empty")
	}

	adb := db.Auth_db.GetDB()
	adb.AutoMigrate(&models.User{})
	err := adb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"enc_password"}),
	}).Create(&users).Error
	return err
}

func grantHandler(rdb *gorm.DB) error {
	// 只需要保证读写auth_db的user表权限
	res := rdb.Exec(fmt.Sprintf(
		"GRANT CREATE, SELECT, INSERT, UPDATE, DELETE, ALTER, INDEX, EXECUTE ON %s.* TO '%s'@'%%'",
		db.Auth_db.Cfg.DB_name, db.GetDBUserName(db.Auth_db.Cfg.DB_name)))
	return res.Error
}

// 初始化数据库和加入相关项目

// 缓存随用随初始，没有必要在这里初始化
