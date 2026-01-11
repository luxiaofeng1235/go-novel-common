package redis_service

import (
	"context"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"log"
)

func LoadRedisKeysNumber(keyword string) (keys []string, keysInfo []models.RedisKeyInfo, err error) {
	var redis = global.Redis
	ctx := context.Background()
	var cursor uint64
	//扫描所有key 每次100条
	var limit int64 = 100
	var num int

	if keyword != "" {
		keyword = fmt.Sprintf("*%s*", keyword)
	} else {
		keyword = "*"
	}
	for {
		var res []string
		//scan 0(cursor) match *ll* count 2
		res, cursor, err = redis.Scan(ctx, cursor, keyword, limit).Result()
		if err != nil {
			log.Println("err:", err)
			return
		}
		keys = append(keys, res...)

		num += 1
		//log.Println("查询次数：", num)

		if cursor == 0 {
			break
		}
	}

	var number int64
	keysInfo = []models.RedisKeyInfo{}
	for _, hashKey := range keys {
		number = redis.HLen(ctx, hashKey).Val()
		if number == 0 {
			number = redis.SCard(ctx, hashKey).Val()
		}
		if number == 0 {
			number = redis.ZCard(ctx, hashKey).Val()
		}
		if redis.Type(ctx, hashKey).Val() == "string" {
			number = 1
		}
		key := models.RedisKeyInfo{
			Key:    hashKey,
			Type:   redis.Type(ctx, hashKey).Val(),
			Number: number,
			Expire: int64(redis.PTTL(ctx, hashKey).Val().Seconds()),
		}
		keysInfo = append(keysInfo, key)
	}
	return
}
