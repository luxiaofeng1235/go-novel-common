package bookcache

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/pkg/gredis"
	"log"
	"time"
)

// 设置缓存的配置信息
func GetTest() (str string) {
	return "111222"
}

// 这个里面主要涉及一些列表的设置更新，可以指定key进行更新
func SetBookListCache(key string, total int64, data interface{}) (err error) {
	if key == "" {
		return
	}
	setData := make(map[string]interface{})
	setData["total"] = total
	setData["data"] = data
	//转换成json数据
	jsonData, err := json.Marshal(setData)
	if err != nil {
		return
	}
	cacheBookList := string(jsonData)
	////默认设置一个小时时间
	err = gredis.Set(key, cacheBookList, time.Hour*1)
	if err != nil {
		return
	}
	return nil
}

// 解析小说的具体结构体
type RedisBookCacheList struct {
	Total int64            `json:"total"`
	Data  []*models.McBook `json:"data"`
}

// 获取今日的缓存书籍列表
func GetBookCacheList(redisKey string) (bookList []*models.McBook, total int64, err error) {
	if redisKey == "" {
		return
	}
	jsonData := gredis.Get(redisKey)
	log.Printf("key = %v value = %v", redisKey, jsonData)
	if jsonData == "" {
		return
	}
	var response RedisBookCacheList
	////解析缓存的结果信息
	err = json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		log.Printf("Error parsing JSON results = %+v", err)
		return
	}
	for _, val := range response.Data {
		fmt.Println(val.BookName)
	}
	if len(response.Data) <= 0 {
		return
	}
	total = response.Total
	bookItem := response.Data
	fmt.Println(bookItem, total)
	return bookItem, total, nil
}

// 获取书籍单个书的获取
func GetBookInfo(bookId int64) {

}
