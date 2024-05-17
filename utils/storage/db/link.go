package db

import (
	"gorm.io/gorm"
)

type SubDB struct {
	db  *gorm.DB
	Cfg *SubDBConfig // 外部组件使用时先绑定，再link
}

type SubDBEntry interface {
	GetDB() *gorm.DB
}

var (
	Auth_db       SubDB
	FileReader_db SubDB
	FileWriter_db SubDB
)

// Before get subordinated db, need to create `SubDB.db“ first
func (s *SubDB) CreateSubDB(grant_func func(*gorm.DB) error) {
	// once, 设置gorm.DB
	DBReadConfigWG.Wait()
	grantSubDB(s, grant_func)
}

// all entry
func (s *SubDB) GetDB() *gorm.DB {
	return s.db
}
