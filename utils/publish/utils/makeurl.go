package publish_utils

import "chatFileBackend/models"

// 制作返回值中的文件url
func GetURL(meta *models.MetaData) string {
	return "/api/download/file?name=" + meta.GenerateObjectName()
}
