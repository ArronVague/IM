package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

func ValidatePasswd(plainPwd, salt, passwd string) bool {
	return Md5Encode(plainPwd+salt) == passwd
}

// MakePasswd 生成用户密码
func MakePasswd(plainPwd, salt string) string {
	return Md5Encode(plainPwd + salt)
}
