package auth_utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/e"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var claims *Claims
		var err error

		code = e.SUCCESS
		token := c.GetHeader("Authorization")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		token = parts[1]

		claims, err = ParseToken(token)
		if err != nil {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    data,
			})

			c.Abort()
			return
		}
		c.Set("level", claims.Level)
		c.Set("username", claims.Username)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}

func GetVerifyLevelWare(level uint8) (mid_ware func() gin.HandlerFunc) {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			// 看看context中有没有level
			var (
				lev interface{}
				ok  bool
			)
			lev, ok = c.Get("level")
			if !ok {
				// 没有
				c.Next()

				lev, ok = c.Get("level")
				if !ok {
					c.JSON(http.StatusUnauthorized, gin.H{"message": "fail to read level information"})
					c.Abort()
					return
				}
			}

			clev, ok := lev.(uint8)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "fail to read level information"})
				c.Abort()
				return
			}
			if clev < level {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "level not reach need"})
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
