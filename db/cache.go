package db

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"go-novel/global"
	"log"
	"time"
)

func InitBigcache() {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Hour))
	if err != nil {
		log.Fatal("err", err.Error())
		return
	}
	global.Bigcache = cache
}
