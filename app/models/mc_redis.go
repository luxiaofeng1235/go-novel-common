package models

// 创建结构体
type RedisKeyInfo struct {
	Key    string `form:"key"    json:"key"`
	Type   string `form:"type"  json:"type"`
	Number int64  `form:"number"   json:"number"`
	Expire int64  `form:"expire"   json:"expire"`
}
