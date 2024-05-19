package publish

import (
	"chatFileBackend/utils/global"
	"chatFileBackend/utils/storage/db"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	once               sync.Once
	fs_storage_db_name = "fs_storage"
	Reader_DB          *gorm.DB
	Writer_DB          *gorm.DB
	// wg                 sync.WaitGroup
)

// upload, search共用初始数据库方法
func Init_handler_entry() {
	once.Do(
		func() {
			global.Init_count.Add(1)
			db.Link(init_handler)
			global.Init_count.Done()
		})
}

func init_handler(root_db *gorm.DB) {
	var writer_s *db.SubUserConfig = &db.DBCfg.SubDBCfgs.FileWriterDB
	var reader_s *db.SubUserConfig = &db.DBCfg.SubDBCfgs.FileReaderDB

	// 检查db是否存在
	exists, err := db.DatabaseExists(root_db, fs_storage_db_name)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("Database:" + fs_storage_db_name + " not exists, create first")
		err = db.DatabaseCreate(root_db, fs_storage_db_name)
		if err != nil {
			logrus.Errorln("Error create database "+
				fs_storage_db_name+"\t:", err)
			return
		}
	}

	// 两个账号
	exists, err = db.UserExists(root_db, writer_s.Username)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("user:" + writer_s.Username + " not exists, create first")
		// 创建账号
		err = db.UserCreate(root_db, writer_s)
		if err != nil {
			logrus.Errorln("Error create user "+
				writer_s.Username+"\t:", err)
			return
		}
	}

	exists, err = db.UserExists(root_db, reader_s.Username)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("user:" + reader_s.Username + " not exists, create first")
		// 创建账号
		err = db.UserCreate(root_db, reader_s)
		if err != nil {
			logrus.Errorln("Error create user "+
				reader_s.Username+"\t:", err)
			return
		}
	}

	grantHandler(root_db)
	if err != nil {
		logrus.Errorln("Error in grant permissions" + fs_storage_db_name)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		reader_s.Username, reader_s.Password, db.DBCfg.Addr, fs_storage_db_name)
	Reader_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Errorln("Can not create link to " + fs_storage_db_name)
		return
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		writer_s.Username, writer_s.Password, db.DBCfg.Addr, fs_storage_db_name)
	Writer_DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Errorln("Can not create link to " + fs_storage_db_name)
		return
	}
}

func grantHandler(rdb *gorm.DB) error {
	// writer
	res := rdb.Exec(fmt.Sprintf(
		"GRANT CREATE, SELECT, INSERT, UPDATE, DELETE, ALTER, INDEX, EXECUTE ON %s.* TO '%s'@'%%'",
		fs_storage_db_name, db.DBCfg.SubDBCfgs.FileWriterDB.Username))

	if res.Error != nil {
		return res.Error
	}

	// reader
	res = rdb.Exec(fmt.Sprintf(
		"GRANT CREATE, SELECT, ALTER, INDEX, EXECUTE ON %s.* TO '%s'@'%%'",
		fs_storage_db_name, db.DBCfg.SubDBCfgs.FileReaderDB.Username))
	return res.Error
}
