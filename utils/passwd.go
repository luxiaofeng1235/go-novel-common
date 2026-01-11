package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// md5加密
func Md5(str string) (md5str string) {
	data := []byte(str)
	has := md5.Sum(data)
	md5str = fmt.Sprintf("%x", has)
	return
}

// md5加盐加密 md5(md5(passwd) + salt)
func GetMd5(str string, salt string) (md5str string) {
	data := []byte(Md5(str) + salt)
	has := md5.Sum(data)
	md5str = fmt.Sprintf("%x", has)
	return
}

// 先base64，然后MD5
func Base64Md5(params string) string {
	return Md5(base64.StdEncoding.EncodeToString([]byte(params)))
}

// 密码加密 使用自适应hash算法, 不可逆
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 通过比较两个字符串hash判断是否出自同一个明文
// hashPasswd 需要对比的密文
// passwd 明文
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetHashingCost(hashedPassword []byte) int {
	cost, _ := bcrypt.Cost(hashedPassword) // 为了简单忽略错误处理
	return cost
}
