package auth_utils

import (
	"net/http"
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
			code = e.INVALID_PARAMS
		} else {
			claims, err = ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
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
		c.Next()
	}
}
