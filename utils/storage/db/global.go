package db

import "gorm.io/gorm"

type db_user struct {
	username string
	password string
}

type db_config struct {
	db2user map[string]db_user // db名-验证用户，仅管理非客户

	addr string
}

var ()

func (d db_config) LinkToDB(db_name string) (*gorm.DB, error) {
	return nil, nil
}

// Fallback用，如果环境中缺少指定用户名的数据库
func (d db_config) CreateDB(db_name string) (*gorm.DB, error) {

}
