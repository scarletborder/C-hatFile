package auth

import (
	auth_utils "chatFileBackend/handlers/auth/utils"
	"chatFileBackend/models"
	cached "chatFileBackend/utils/storage/cache"
	"fmt"
	"strconv"
	"time"

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
	// 毫秒
	time_stamp_str, ok := c.GetPostForm("time_stamp")
	if !ok {
		c.JSON(501, gin.H{
			"message": "need arg time_stamp"})
		return
	}
	timestamp, err := strconv.Atoi(time_stamp_str)
	if err != nil {
		c.JSON(501, gin.H{
			"message": "time_stamp must be integer"})
		return

	}

	if int64(timestamp) < (time.Now().UnixNano()/int64(time.Millisecond))-10000 {
		c.JSON(401, gin.H{"message": "timestamp submitted is too later than server > 10s!"})
		return
	}

	if ok, lev := LoginVerify(username, enc2_pwd, timestamp); ok {
		// 登录成功
		token, exp_time, err := auth_utils.GenerateToken(username, enc2_pwd, lev)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Login succeed, but crash in generating token"})
			return
		}

		c.JSON(200, gin.H{
			"token":        token,
			"message":      "success",
			"expire_stamp": exp_time.UnixNano() / int64(time.Millisecond)})
	} else {
		c.JSON(401, gin.H{
			"message": "password is incorrect, or username doesn't exist"})
	}
}

// @param enc2_pwd 传来的经过sha256后拼接时间戳盐再sha256的结果
func LoginVerify(username, enc2_pwd string, time_stamp int) (bool, uint8) {
	real_enc_pwd, level := getEncPwd(username)
	real_enc2_pwd := auth_utils.Sha256Str(real_enc_pwd + fmt.Sprint(time_stamp))
	return real_enc2_pwd == enc2_pwd, level
}

// 获得一次sha256
func getEncPwd(username string) (enc_pwd string, level uint8) {
	var user models.User
	user.Username = username
	user.FlushDirty()

	ok, err := cached.CacheGetByStr(cached.TypeAuthCache, &user)
	if err != nil {
		logrus.Warnln("无法从缓存中获取，Fallback至DB")
	}
	if !ok {
		adb := auth_utils.Auth_DB
		adb.AutoMigrate(&models.User{})
		adb.Where("username = ?", username).Take(&user)
		cached.CacheSetByStr(cached.TypeAuthCache, &user)
		return user.Enc_password, user.Level

	} else {
		return user.Enc_password, user.Level
	}
}
