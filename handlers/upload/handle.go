package upload

import (
	"chatFileBackend/handlers/upload/upload_utils"
	"chatFileBackend/models"
	"chatFileBackend/utils/publish"
	publish_utils "chatFileBackend/utils/publish/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	// 发现context中的username
	// username_i, ok := c.Get("username")
	uid := c.GetUint64("uid")
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "fail to read username in request"})
	// 	return
	// }
	// username, ok := username_i.(string)
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "fail in asserting username"})
	// 	return
	// }

	// debug
	// bb, _ := io.ReadAll(c.Request.Body)
	// fmt.Println(string(bb))

	// 获取表单中的字符串
	tags_str, ok := c.GetPostForm("tags")
	if !ok {
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

	tags_str_arr := publish_utils.Str2Tags(tags_str)
	var tags_arr []models.Tag
	for _, ts := range tags_str_arr {
		tags_arr = append(tags_arr, publish.NewTag(ts))
	}

	// 尝试上传到对象存储和将元数据传到db中
	t := time.Now()
	_, err = upload_utils.UploadFile(fileContent, &models.MetaData{
		Name:       file.Filename,
		Size:       file.Size,
		Tags:       tags_arr,
		UploadTime: &t,
		UserID:     uid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create record"})
		return
	}

	// 成功
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
