package db

import (
	"chatFileBackend/utils/global"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// type SubDB struct {
// 	db  *gorm.DB
// 	Cfg *SubUserConfig // 外部组件使用时先绑定，再link
// }

// var (
// 	Auth_db       SubDB
// 	FileReader_db SubDB
// 	FileWriter_db SubDB
// )

// 初始化连接
func Link(init_func func(rdb *gorm.DB)) {
	global.Init_count.Add(1)
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorln("Error in initiating link " + fmt.Sprint(err))
		}
	}()

	// once, 设置gorm.DB
	DBReadConfigWG.Wait()
	go root_once.Do(func() {
		StartRootDSN()
	})

	go func() {
		Root_link_wg.Wait()

		init_func(root_db)
		err := root_db.Exec("FLUSH PRIVILEGES").Error
		if err != nil {
			panic(err)
		}
		global.Init_count.Done()
	}()

}

// // Before get subordinated db, need to create `SubDB.db“ first
// func (s *SubUserConfig) CreateUser(init_func func(rdb *gorm.DB, scfg *SubUserConfig) error) {

// }

// // all entry
// func (s *SubDB) GetDB() *gorm.DB {
// 	return s.db
// }
