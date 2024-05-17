package auth

import (
	"time"

	Jwt "github.com/dgrijalva/jwt-go"
)

// 生成JWT Token
func generateToken(username string) (string, error) {
	token := Jwt.NewWithClaims(Jwt.SigningMethodHS256, Jwt.MapClaims{
		"username": username,                              // 将用户名编入Token
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token有效期
	})

	// 用一个安全的密钥签名Token
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	return tokenString, err
}
