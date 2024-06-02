package manager

import (
	"chatFileBackend/models"
	"chatFileBackend/utils/publish"
	chats3 "chatFileBackend/utils/publish/s3"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SelfViewHandler(c *gin.Context) {
	// 获取context中的uid
	uid := c.GetUint64("uid")

	var metas []models.MetaData
	publish.Reader_DB.AutoMigrate(&models.MetaData{})

	res := publish.Reader_DB.Preload("Tags").Model(&models.MetaData{}).
		Where("user_id = ?", uid).Find(&metas)
	num := res.RowsAffected

	c.JSON(200, UserResource{Files: metas, FileNum: num})
}

func DeleteHandler(c *gin.Context) {
	// 获取context中的uid
	uid := c.GetUint64("uid")

	fid, ok := c.GetQuery("fid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "require one file id"})
		return
	}

	var meta models.MetaData
	publish.Reader_DB.AutoMigrate(&models.MetaData{})

	res := publish.Reader_DB.First(&meta, fid)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("error in query fid: %s", res.Error.Error())})
		return
	}

	if meta.UserID != uid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "this file not belong to you"})
		return
	}

	// 验证成功
	res = publish.Writer_DB.Delete(&meta)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("error in delete: %s", res.Error.Error())})
		return
	}
	// 继续删除oss
	chats3.DeleteFile(&meta)

	c.JSON(200, gin.H{"message": "success"})
}
