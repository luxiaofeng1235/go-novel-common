/*
 * @Descripttion: Redis 初始化（从 config.yml 读取）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-novel/config"
	"go-novel/global"
	"log"
	"strings"
)

func GetRedis() (addr string, passwd string, defaultdb int) {
	addr = strings.TrimSpace(config.GetString("redis.addr"))
	if addr == "" {
		host := strings.TrimSpace(config.GetString("redis.host"))
		port := config.GetInt("redis.port")
		if port == 0 {
			port = 6379
		}
		if host == "" {
			host = "127.0.0.1"
		}
		if strings.Contains(host, ":") {
			addr = host
		} else {
			addr = fmt.Sprintf("%s:%d", host, port)
		}
	}
	passwd = config.GetString("redis.password")
	defaultdb = config.GetInt("redis.db")
	return
}

// Redis 初始化Redis
func InitRedis(addr string, passwd string, defaultdb int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		//Addr:     GetRedisConnString(),
		//Password: config.GetString("gredis.password"),
		//DB:       0,
		Addr:     addr,
		Password: passwd,
		DB:       defaultdb,
	})

	var ctx = context.Background()
	err := redisClient.Ping(ctx).Err()
	if err != nil || redisClient == nil {
		log.Fatalln(fmt.Sprintf("初始化Redis异常：%v", err))
	} else {
		log.Println("gredis connect success")
	}

	global.Redis = redisClient
	return redisClient
}
