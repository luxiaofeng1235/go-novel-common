package db

import (
	"go-novel/global"
	"go-novel/pkg/ws"
	"log"
	"runtime"
)

func InitWs() {
	// 根据 CPU 核心数自动计算分片数（最小 2，最大 16）
	shardCount := runtime.NumCPU()
	if shardCount < 2 {
		shardCount = 2
	}
	if shardCount > 16 {
		shardCount = 16
	}

	hubManager := ws.NewHubManager(shardCount)
	global.WsHubManager = hubManager
	hubManager.Run()

	log.Printf("WebSocket HubManager 已启动：%d 个分片（CPU 核心数：%d）", shardCount, runtime.NumCPU())
}
