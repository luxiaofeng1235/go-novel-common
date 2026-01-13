/*
 * @Descripttion: WebSocket HubManager（分片 Hub 架构，提升并发能力）
 * @Author: red
 * @Date: 2026-01-13 12:30:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 12:30:00
 */
package ws

import (
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// HubManager 管理多个 Hub 分片，按 userID 哈希分配，实现水平扩展
type HubManager struct {
	shards      []*Hub
	shardCount  int
	totalConns  int64 // 原子计数器：总连接数
	totalAuthed int64 // 原子计数器：已鉴权连接数
}

// NewHubManager 创建指定数量分片的 HubManager（推荐 4-8 个分片）
func NewHubManager(shardCount int) *HubManager {
	if shardCount <= 0 {
		shardCount = 4 // 默认 4 个分片
	}

	shards := make([]*Hub, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = NewHub()
	}

	return &HubManager{
		shards:     shards,
		shardCount: shardCount,
	}
}

// Run 启动所有 Hub 分片（每个分片一个 goroutine）
func (hm *HubManager) Run() {
	if hm == nil {
		return
	}
	for _, hub := range hm.shards {
		go hub.Run()
	}
}

// GetHub 根据 userID 获取对应的 Hub 分片（相同 userID 始终路由到同一分片）
func (hm *HubManager) GetHub(userID int64) *Hub {
	if hm == nil || len(hm.shards) == 0 {
		return nil
	}
	// 游客（userID=0）随机分配，已登录用户按 userID 哈希
	if userID <= 0 {
		// 使用原子计数器轮询分配游客连接
		idx := atomic.AddInt64(&hm.totalConns, 1) % int64(hm.shardCount)
		return hm.shards[idx]
	}
	// 已登录用户按 userID 取模，保证同一用户的多端连接在同一分片
	return hm.shards[userID%int64(hm.shardCount)]
}

// RegisterClient 注册客户端到对应的 Hub 分片
func (hm *HubManager) RegisterClient(conn *websocket.Conn, userID int64, username string) *Client {
	if hm == nil {
		return nil
	}
	hub := hm.GetHub(userID)
	if hub == nil {
		return nil
	}
	client := NewClient(hub, conn, userID, username)
	client.Start()

	// 统计计数
	atomic.AddInt64(&hm.totalConns, 1)
	if userID > 0 {
		atomic.AddInt64(&hm.totalAuthed, 1)
	}

	return client
}

// Broadcast 向所有分片的所有连接广播消息
func (hm *HubManager) Broadcast(msg []byte) {
	if hm == nil {
		return
	}
	for _, hub := range hm.shards {
		hub.Broadcast(msg)
	}
}

// Stats 返回统计信息
func (hm *HubManager) Stats() map[string]interface{} {
	if hm == nil {
		return nil
	}
	return map[string]interface{}{
		"shardCount":  hm.shardCount,
		"totalConns":  atomic.LoadInt64(&hm.totalConns),
		"totalAuthed": atomic.LoadInt64(&hm.totalAuthed),
	}
}
