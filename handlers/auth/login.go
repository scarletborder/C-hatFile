package auth

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/models"
	cached "chatFileBackend/utils/storage/cache"
	"chatFileBackend/utils/storage/db"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoginHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(500, gin.H{
				"message": "unexpected error" + fmt.Sprint(err)})
		}
	}()
	// a, err := io.ReadAll(c.Request.Body)
	// _ = a
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg username"})
		return
	}
	enc2_pwd, ok := c.GetPostForm("encrypted_pwd")
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg encrypted_pwd"})
		return
	}
	time_stamp, ok := c.GetPostForm("time_stamp")
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg time_stamp"})
		return
	}
	timestamp, err := strconv.Atoi(time_stamp)
	if err != nil {
		c.JSON(501, gin.H{
			"message": "time_stamp must be integer"})
		return

	}
	if LoginVerify(username, enc2_pwd, timestamp) {
		// 发回token
		c.JSON(200, gin.H{
			"token":   "123abc",
			"message": "success"})
	} else {
		c.JSON(401, gin.H{
			"message": "password is incorrect, or username doesn't exist"})
	}
}

// @param enc2_pwd 传来的经过sha256后拼接时间戳盐再sha256的结果
func LoginVerify(username, enc2_pwd string, time_stamp int) bool {
	real_enc_pwd := getEncPwd(username)
	real_enc2_pwd := auth_utils.Sha256Str(real_enc_pwd + fmt.Sprint(time_stamp))
	return real_enc2_pwd == enc2_pwd
}

// 获得一次sha256
func getEncPwd(username string) (enc_pwd string) {
	var user models.User
	ok, err := cached.CacheGetByStr(username, &user)
	if err != nil {
		logrus.Warnln("无法从缓存中获取，Fallback至DB")
	}
	if !ok {
		adb := db.Auth_db.GetDB()
		adb.AutoMigrate(&models.User{})
		user.Username = username
		adb.Take(&user)
		return user.Enc_password

	} else {
		return user.Enc_password
	}
}
