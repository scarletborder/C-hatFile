package upload

import (
	"chatFileBackend/handlers/upload/utils"
	"chatFileBackend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// 获取header中的authentication
	authHeader := c.GetHeader("authentication")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "authentication header is required"})
		return
	}

	// 获取表单中的字符串
	tags_str := c.PostForm("tags")
	if tags_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "tags field is required"})
		return
	}

	// 获取表单中的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "file field is required"})
		return
	}

	// 打开文件并得到io.Reader
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not open file"})
		return
	}
	defer fileContent.Close()

	tags_str_arr := utils.Str2Tags(tags_str)
	var tags_arr []models.Tag
	for _, ts := range tags_str_arr {
		tags_arr = append(tags_arr, models.Tag{Title: ts})
	}

	// 尝试上传到对象存储和将元数据传到db中
	_, err = utils.UploadFile(fileContent, models.MetaData{
		Name:   file.Filename,
		Size:   file.Size,
		Tags:   tags_arr,
		UserID: 1223, // 使用token找到目前的userid
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create record"})
		return
	}

	// 成功
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
