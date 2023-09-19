package pkg

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + " " + password)
	return md == password
}

// token生成
func GenerateLongHash() string {
	// 获取当前时间戳
	str := fmt.Sprintf("%d", RandomString(10))
	// 使用MD5加密
	md5Hash := md5.Sum([]byte(str))
	md5Str := hex.EncodeToString(md5Hash[:])
	// 使用SHA-256加密
	sha256Hash := sha256.Sum256([]byte(str))
	sha256Str := hex.EncodeToString(sha256Hash[:])
	// 组合两个哈希值
	longHash := md5Str + sha256Str
	return longHash
}
