package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-novel/config"
	"go-novel/global"
	"go-novel/utils"
	"log"
)

func GetRedis() (addr string, passwd string, defaultdb int) {
	//初始化数据库
	env := config.GetString("server.env")
	addr = "127.0.0.1:6379"
	passwd = ""
	defaultdb = 0
	if env == utils.Dev {
		passwd = "DmDw8vmGGe"
	} else if env == utils.Prod {
		addr = "192.168.10.17:6379" //使用96.35的redis进行链接
		passwd = "cCoF3Yrqd9"
	} else if env == utils.Local {

	}
	//passwd = "o9kHvO95bP"
	return
}

// Redis 初始化Redis
func InitRedis(addr string, passwd string, defaultdb int) *redis.Client {
	redisClinet := redis.NewClient(&redis.Options{
		//Addr:     GetRedisConnString(),
		//Password: config.GetString("gredis.password"),
		//DB:       0,
		Addr:     addr,
		Password: passwd,
		DB:       defaultdb,
	})

	var ctx = context.Background()
	err := redisClinet.Ping(ctx).Err()
	if err != nil || redisClinet == nil {
		log.Fatalln(fmt.Sprintf("初始化Redis异常：%v", err))
	} else {
		log.Println("gredis connect success")
	}

	global.Redis = redisClinet
	return redisClinet
}
