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

type id_result struct {
	UserID uint64
	Count  int64
}

var (
	rank_result  []result
	update_stamp int64
)

func statRank() {
	var tmp_id_result []id_result
	global.Init_count.Wait()
	Db := publish.Reader_DB
	Db.AutoMigrate(&models.MetaData{})

	Db.Model(&models.MetaData{}).Select("user_id, count(*) as count").
		Group("user_id").Order("count desc").Limit(6).Find(&tmp_id_result)

	rank_result = rank_result[:0]
	Db.AutoMigrate(&models.User{})
	for _, val := range tmp_id_result {
		var tmp_user models.User
		Db.First(&tmp_user, val.UserID)
		rank_result = append(rank_result, result{Count: val.Count, Username: tmp_user.Username})
	}
	update_stamp = time.Now().UnixNano() / int64(time.Millisecond)
}
