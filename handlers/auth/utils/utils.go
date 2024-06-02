package auth_utils

import (
	"chatFileBackend/models"
	cached "chatFileBackend/utils/storage/cache"
	"crypto/sha256"
	"fmt"

	"gorm.io/gorm"
)

var Auth_DB *gorm.DB

func Sha256Str(tok string) string {
	// 创建一个新的SHA256哈希对象
	hash := sha256.New()
	// 写入数据到哈希对象
	hash.Write([]byte(tok))
	// 计算哈希值
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	hashString := fmt.Sprintf("%x", hashBytes)
	return hashString
}

// 由修改密码或注册账号引发
func UpdateAccount(user *models.User) {
	cached.CacheSetByStr(cached.TypeAuthCache, user)
}

// 获得context中的用户id
// 通常是在用户登录后，每次请求带有token
