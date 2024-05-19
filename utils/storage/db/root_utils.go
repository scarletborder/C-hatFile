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
	Root_link_wg sync.WaitGroup
	root_db      *gorm.DB
)

func init() {
	Root_link_wg.Add(1)
}

func StartRootDSN() {
	defer func() {
		Root_link_wg.Done()
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

func DatabaseExists(root_db *gorm.DB, dbname string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s')",
		dbname)
	err := root_db.Raw(query).Row().Scan(&exists)
	return exists, err
}

func DatabaseCreate(root_db *gorm.DB, dbname string) error {
	res := root_db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	return res.Error
}

func UserExists(root_db *gorm.DB, username string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT User FROM mysql.user WHERE User = '%s')", username)
	err := root_db.Raw(query).Row().Scan(&exists)
	return exists, err
}

func UserCreate(root_db *gorm.DB, s *SubUserConfig) error {
	res := root_db.Exec(fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'",
		s.Username, s.Password))

	return res.Error
}
