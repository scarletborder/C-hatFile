package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
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

func LoginVerify(username, enc2_pwd string, time_stamp int) bool {
	real_enc_pwd := sha256Str("!!8964jss")
	// 涉及缓存 username

	real_enc2_pwd := sha256Str(real_enc_pwd + fmt.Sprint(time_stamp))
	return real_enc2_pwd == enc2_pwd
}
