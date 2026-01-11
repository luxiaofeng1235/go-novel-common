package book_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/service/common/redis_service"
	"log"
	"time"
)

type User struct {
	Id       int
	Username string
	Age      int
}

// 获取书籍排行列表
func SetBookRankList(key string, value interface{}, timeout time.Duration) (content string) {
	//特殊情况判断关联性
	if key == "" || value == nil {
		return
	}
	//设置获取的缓存信息
	err := redis_service.Set(key, value, timeout)
	if err != nil {
		err = fmt.Errorf("redis缓存失败 err=%v", err.Error())
		return
	} else {
		fmt.Println("gredis data ok")
	}
	fmt.Println(value)
	return "ok"
}

// 根据指定的key获取对应的排行列表数据
func GetBookListByCacheKey(key string) (contents string) {
	if key == "" {
		return ""
	}
	listBooKVal := redis_service.Get(key)
	if listBooKVal == "" {
		return ""
	}
	fmt.Println(listBooKVal)
	var userinfo User
	err := json.Unmarshal([]byte(listBooKVal), &userinfo)
	if err != nil {
		log.Println("解析失败 err=", err.Error())
		return ""
	}
	fmt.Println(userinfo)
	return userinfo.Username
}
