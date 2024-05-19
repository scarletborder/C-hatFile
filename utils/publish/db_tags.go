package publish

import "chatFileBackend/models"

func ReplaceExistTags(tags *[]models.Tag) {
	Db := Writer_DB
	Db.AutoMigrate(&models.Tag{})
	for i, tag := range *tags {
		var existingTag models.Tag
		if err := Db.Where("title = ?", tag.Title).First(&existingTag).Error; err == nil {
			// 标签已存在，使用现有的标签
			(*tags)[i] = existingTag
		}
	}
}

func NewTag(tag_title string) models.Tag {
	var existingTag models.Tag
	Db := Writer_DB

	Db.AutoMigrate(&models.Tag{})
	if err := Db.Model(&models.Tag{}).Where("title = ?", tag_title).
		First(&existingTag).Error; err == nil {
		// 标签已存在，使用现有的标签
		return existingTag
	}
	return models.Tag{Title: tag_title}
}

func ListSimilarTag(tag_title string) (ret []models.Tag) {
	Db := Writer_DB
	Db.AutoMigrate(&models.Tag{})
	Db.Model(&models.Tag{}).Where("title LIKE ?", "%"+tag_title+"%").
		Find(&ret)
	return
}
