package db

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Fallback 至此，在Link时没有命中，root创建相关数据库和用户

var (
	root_once    sync.Once
	root_link_wg sync.WaitGroup
	root_db      *gorm.DB
)

func init() {
	root_link_wg.Add(1)
}

// 保证存在对应子db
func grantSubDB(s *SubDB, grant_func func(*gorm.DB) error) {
	// 到达fallback最底层，创建所需数据库和用户
	root_once.Do(func() {
		startRootDSN()
	})
	root_link_wg.Wait()

	// 使用root账号检查
	exists, err := databaseExists(root_db, s.Cfg.DB_name)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("Database:" + s.Cfg.DB_name + "not exists, create first")
		err = databaseCreate(root_db, s.Cfg.DB_name)
		if err != nil {
			logrus.Errorln("Error create database "+
				s.Cfg.DB_name+"\t:", err)
			return
		}
	}

	user := GetDBUserName(s.Cfg.DB_name)
	pwd := s.Cfg.Password
	exists, err = userExists(root_db, user)
	if err != nil {
		logrus.Errorln("Error checking user existence:", err)
		return
	}
	if !exists {
		logrus.Warnln("user:" + user + "not exists, create first")
		// 创建账号
		err = userCreate(root_db, s.Cfg, grant_func)
		if err != nil {
			logrus.Errorln("Error create user "+
				user+"\t:", err)
			return
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, pwd, DBCfg.Addr)
	s.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Errorln("Can not create link to " + s.Cfg.DB_name)
		return
	}
}

func startRootDSN() {
	defer func() {
		root_link_wg.Done()
	}()

	rootDSN := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		DBCfg.RootUser, DBCfg.RootPassword, DBCfg.Addr)

	var err error
	root_db, err = gorm.Open(mysql.Open(rootDSN), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalln("Root connection failed:", err)
		return
	} else {
		logrus.Infoln("root db initiated successfully")
	}

}

func CloseRootDSN() {
	// 没有
}

func databaseExists(root_db *gorm.DB, dbname string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s')",
		dbname)
	err := root_db.Raw(query).Row().Scan(&exists)
	return exists, err
}

func databaseCreate(root_db *gorm.DB, dbname string) error {
	res := root_db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	return res.Error
}

func userExists(root_db *gorm.DB, username string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT User FROM mysql.user WHERE User = '%s')", username)
	err := root_db.Raw(query).Row().Scan(&exists)
	return exists, err
}

func userCreate(root_db *gorm.DB, s *SubDBConfig, grant_func func(*gorm.DB) error) error {
	res := root_db.Exec(fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'",
		GetDBUserName(s.DB_name), s.Password))
	err := res.Error
	if err != nil {
		return err
	}
	// 赋权账号
	err = grant_func(root_db)
	if err != nil {
		return err
	}
	res = root_db.Exec("FLUSH PRIVILEGES")
	return res.Error
}
