package search

import (
	"chatFileBackend/models"
	"chatFileBackend/utils/publish"
	publish_utils "chatFileBackend/utils/publish/utils"
	"strings"
)

func search(title string, tags []string) (res []SearchResult) {
	Db := publish.Reader_DB
	var metadatas []models.MetaData

	Db.AutoMigrate(&models.Tag{}, &models.MetaData{})

	query := Db.Preload("Tags").Model(&models.MetaData{})

	title = strings.TrimSpace(title)
	if title != "" { // 过滤标题
		query = query.Where("name LIKE ?", "%"+title+"%")
	}

	if len(tags) > 0 { // 如果有tag
		query = query.Joins("JOIN metadata_tags ON metadata_tags.meta_data_id = meta_data.id").
			Joins("JOIN tags ON tags.id = metadata_tags.tag_id").
			Where("tags.title IN (?)", tags).
			Group("meta_data.ID").
			Having("COUNT(DISTINCT tags.ID) = ?", len(tags)) //where中拣选出的tag个数和表中相等，
		// 意味 having的是原来数据库的子集
	}

	query.Find(&metadatas)

	for _, meta := range metadatas {
		ori_tags := []string{}
		for _, tag_str := range meta.Tags {
			ori_tags = append(ori_tags, tag_str.Title)
		}

		res = append(res, SearchResult{
			Title: meta.Name,
			URL:   publish_utils.GetURL(&meta),
			Tags:  ori_tags})
	}
	return
}
