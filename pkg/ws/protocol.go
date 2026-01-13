/*
 * @Descripttion: WebSocket 协议定义与基础消息处理
 * @Author: red
 * @Date: 2026-01-13 11:40:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:55:00
 */
package ws

import (
	"encoding/json"
	"strings"
	"time"
)

type inboundMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type outboundMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func (c *Client) handleMessage(message []byte) {
	if c == nil || c.hub == nil {
		return
	}

	raw := strings.TrimSpace(string(message))
	if raw == "" {
		return
	}

	// 兼容：若不是 JSON，则按 chat 文本处理
	var in inboundMessage
	if err := json.Unmarshal([]byte(raw), &in); err != nil || strings.TrimSpace(in.Type) == "" {
		c.sendChatText(raw)
		return
	}

	switch strings.ToLower(strings.TrimSpace(in.Type)) {
	case "ping":
		c.replyPong()
	case "chat":
		var payload struct {
			Text string `json:"text"`
		}
		_ = json.Unmarshal(in.Data, &payload)
		c.sendChatText(payload.Text)
	case "dm":
		var payload struct {
			ToUserID int64  `json:"toUserId"`
			Text     string `json:"text"`
		}
		_ = json.Unmarshal(in.Data, &payload)
		c.sendDM(payload.ToUserID, payload.Text)
	default:
	}
}

func (c *Client) replyPong() {
	if c == nil {
		return
	}
	msg := outboundMessage{
		Type: "pong",
		Data: map[string]interface{}{
			"ts": time.Now().Unix(),
		},
	}
	select {
	case c.send <- mustJSON(msg):
	default:
	}
}
