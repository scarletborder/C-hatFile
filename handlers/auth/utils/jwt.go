package auth_utils

import (
	"time"

	Jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtSecretKeyStr = `d5vg,t^kt6,#gN_Ue7bT#P*whDzrVQirG]4!^tP6ZfB8VCGJD0nAyrD6@b7-!96ZQ_3Ao,NWUWRfhhK4K*bCBU)11k_QJ*DJGt0i}d5@YyPrZUb7@KtD#6vKY=-myE>g`
)

var jwtSecret = []byte(jwtSecretKeyStr)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Level    uint8  `json:"level"` // 等级
	Jwt.StandardClaims
}

func GenerateToken(username, password string, level uint8) (string, time.Time, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		level,
		Jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := Jwt.NewWithClaims(Jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, expireTime, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := Jwt.ParseWithClaims(token, &Claims{},
		func(token *Jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
