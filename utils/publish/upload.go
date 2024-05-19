package publish

import (
	"chatFileBackend/models"
	chats3 "chatFileBackend/utils/publish/s3"
	"io"
)

func UploadDocument(file io.Reader, meta *models.MetaData) (msg string, err error) {
	// 首先尝试塞入db
	Writer_DB.AutoMigrate(&models.MetaData{}, &models.Tag{})
	res := Writer_DB.Create(meta)

	if res.Error != nil {
		return "fail in insert to DB", res.Error
	}
	// 其次尝试对象存储
	metaID := meta.ID
	msg, err = chats3.Upload_file(file, meta)
	if err != nil {
		// 如果对象存储不行，删除db相关line
		Writer_DB.Delete(&models.MetaData{}, metaID)
		// 删除连接表
		Writer_DB.Table("metadata_tags").Where("meta_data_id = ?", metaID).Delete(nil)
		return
	}

	return "success", nil
}
