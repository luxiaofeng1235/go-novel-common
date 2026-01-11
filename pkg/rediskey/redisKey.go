package rediskey

import (
	"fmt"
	"math/rand"
	"time"
)

// 基础缓存时间 1 小时
const baseTTL = time.Hour

// 随机增加 0 到 10 分钟的过期时间
var randomTTL = time.Duration(rand.Intn(600)) * time.Second

// 最终缓存过期时间
var cacheTTL = baseTTL + randomTTL

const VivoUserIsRetention = "vivo_user_is_retention_%d"
const prefix = "api:"
const BookInfo = prefix + "book:info:book_%d"        //图片详情
const WhiteList = "yq_white_list"                    //白名单
const TodayBookList = prefix + "book:today_bok_list" //今日更新

func GetVivoUserIsRetention(user_id int64) string {
	return fmt.Sprintf(VivoUserIsRetention, user_id)
}

// GetBookInfoKey 获取书籍详情信息
func GetBookInfoKey(book_id int) string {
	return fmt.Sprintf(BookInfo, book_id)
}

// GetWhiteListKey 获取白名单key
func GetWhiteListKey() string {
	return WhiteList
}

// 获取今日更新的书籍key
func GetTodayBookKey() string {
	return TodayBookList
}

// 获取排行榜单数据的key
func GetBookRakList() {

}
