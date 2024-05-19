package stat

import (
	"chatFileBackend/models"
	"chatFileBackend/utils/global"
	"chatFileBackend/utils/publish"
	"time"
)

func init() {
	publish.Init_handler_entry()
	ticker := time.NewTicker(12 * time.Hour)

	statRank()
	go func() {
		for range ticker.C {
			statRank()
		}
	}()
}

type result struct {
	Username string `json:"name"`
	Count    int64  `json:"number"`
}

var (
	rank_result  []result
	update_stamp int64
)

func statRank() {
	global.Init_count.Wait()
	Db := publish.Reader_DB
	Db.AutoMigrate(&models.MetaData{})

	Db.Model(&models.MetaData{}).Select("username, count(*) as count").
		Group("username").Order("count desc").Limit(6).Find(&rank_result)

	update_stamp = time.Now().UnixNano() / int64(time.Millisecond)
}
