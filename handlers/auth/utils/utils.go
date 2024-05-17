package auth_utils

import (
	"crypto/sha256"
	"fmt"
)

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
