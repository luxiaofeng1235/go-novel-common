package main

import (
	"context"
	"go-novel/db"
	"go-novel/global"
	"log"
	"time"
)

func main() {
	addr, passwd, _ := db.GetRedis()
	db.InitRedis(addr, passwd, 10)
	var redis = global.Redis
	ctx := context.Background()
	start := time.Now()
	log.Println(redis.Get(ctx, "novel_info_key:101225"))
	log.Println(time.Since(start))
}
