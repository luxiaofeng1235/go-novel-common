package db

import (
	"go-novel/global"
	"go-novel/pkg/ws"
	"log"

	"github.com/spf13/viper"
)

func InitWs() {
	// 从配置读取分片数量（默认 4 个分片，可配置 1-16）
	shardCount := viper.GetInt("websocket.shardCount")
	if shardCount <= 0 {
		shardCount = 4 // 默认 4 个分片
	}
	if shardCount > 16 {
		shardCount = 16 // 最多 16 个分片（避免过度分片）
	}

	hubManager := ws.NewHubManager(shardCount)
	global.WsHubManager = hubManager
	hubManager.Run()

	log.Printf("WebSocket HubManager 已启动：%d 个分片", shardCount)
}
