package auth

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/models"
	"chatFileBackend/utils/storage/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(500, gin.H{
				"message": "unexpected error" + fmt.Sprint(err)})
		}
	}()
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg username"})
		return
	}

	// 检测是否有注册过的用户
	adb := db.Auth_db.GetDB()

	adb.AutoMigrate(&models.User{})

	// 检查数据库中是否存在具有指定名称的用户
	var count int64
	nameToCheck := "exampleName"
	adb.Model(&models.User{}).Where("name = ?", nameToCheck).Count(&count)

	if count > 0 {
		c.JSON(406, gin.H{
			"message": "same username has been registered"})
		return
	}

	enc_pwd, ok := c.GetPostForm("encrypted_pwd") // 一次sha256
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg encrypted_pwd"})
		return
	}

	auth_utils.UpdateAccount(models.User{Username: username,
		Enc_password: enc_pwd})

	c.JSON(200, gin.H{
		"message": "successfully register your account"})
}
