package auth

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/models"
	"chatFileBackend/utils/global"
	cached "chatFileBackend/utils/storage/cache"
	"chatFileBackend/utils/storage/db"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const auth_db_name = "auth_db"

func init() {
	global.Init_count.Add(1)
	go func() {
		db.Link(init_handler)
		cached.StartDBSync(cached.TypeAuthCache, syncHandler, cached.Sync_interval)
		global.Init_count.Done()
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

	adb := auth_utils.Auth_DB
	adb.AutoMigrate(&models.User{})
	err := adb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoUpdates: clause.AssignmentColumns([]string{"enc_password"}),
	}).Create(&users).Error
	return err
}

func init_handler(root_db *gorm.DB) {
	var s *db.SubUserConfig = &db.DBCfg.SubDBCfgs.AuthDB
	// 使用root账号检查
	exists, err := db.DatabaseExists(root_db, auth_db_name)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("Database:" + auth_db_name + " not exists, create first")
		err = db.DatabaseCreate(root_db, auth_db_name)
		if err != nil {
			logrus.Errorln("Error create database "+
				auth_db_name+"\t:", err)
			return
		}
	}

	user := s.Username
	pwd := s.Password
	exists, err = db.UserExists(root_db, user)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("user:" + user + " not exists, create first")
		// 创建账号
		err = db.UserCreate(root_db, s)
		if err != nil {
			logrus.Errorln("Error create user "+
				user+"\t:", err)
			return
		}
	}

	err = grantHandler(root_db)
	if err != nil {
		logrus.Errorln("Error in grant permissions")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pwd, db.DBCfg.Addr, auth_db_name)
	auth_utils.Auth_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Errorln("Can not create link to " + auth_db_name)
		return
	}
}

// 赋予部分账号规定的权限
func grantHandler(rdb *gorm.DB) error {
	// 只需要保证读写auth_db的user表权限
	res := rdb.Exec(fmt.Sprintf(
		"GRANT CREATE, SELECT, INSERT, UPDATE, DELETE, ALTER, INDEX, EXECUTE ON %s.* TO '%s'@'%%'",
		auth_db_name, db.DBCfg.SubDBCfgs.AuthDB.Username))
	return res.Error
}

// 初始化数据库和加入相关项目

// 缓存随用随初始，没有必要在这里初始化
