package utils

import (
	"log"
	"strconv"
	"strings"
)

/*
* @note 数值切割成字符串
* @param nums int64 输入字符
* @return object
 */
func JoinInt64ToString(nums []int64) string {
	var stringSlice []string
	for _, num := range nums {
		stringSlice = append(stringSlice, strconv.FormatInt(num, 10))
	}
	return strings.Join(stringSlice, ", ")
}

/*
* @note 判断最后两位字符是否为==或者最后一位是=
* @param s string string 输入字符
* @return bool
 */
func RemoveEqualSigns(s string) string {
	if s == "" {
		return ""
	}
	n := len(s)
	if n == 0 {
		return s
	}
	// 检查最后两个字符是否为 ==
	if n > 1 && s[n-2] == '=' && s[n-1] == '=' {
		newStr := s[:n-2]
		log.Printf("查找到【%s】包含末尾有两个【=】 ,新的字符为：%s", s, newStr)
		return newStr // 去掉最后两个字符
	}
	// 检查最后一个字符是否为 =
	if s[n-1] == '=' {
		secondStr := s[:n-1]
		log.Printf("查找到【%s】包含末尾有一个【=】 ,新的字符为：%s", s, secondStr)
		return secondStr // 去掉最后一个字符
	}
	log.Printf("未发现特殊字符【=】，不需要替换，原样返回字符串 【%s】", s)
	return s // 如果都不是，返回原字符串
}
