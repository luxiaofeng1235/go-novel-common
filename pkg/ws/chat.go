/*
 * @Descripttion: WebSocket 群聊（无房间，广播给所有已连接用户）
 * @Author: red
 * @Date: 2026-01-13 11:40:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:55:00
 */
package ws

import (
	"strings"
	"time"
)

type chatRequest struct {
	client *Client
	text   string
}

func (c *Client) sendChatText(text string) {
	if c == nil || c.hub == nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	select {
	case c.hub.chat <- chatRequest{client: c, text: text}:
	default:
	}
}

func (h *Hub) broadcastChat(c *Client, text string) {
	if h == nil || c == nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	// 约定：群聊只对“带 token 建立连接”的用户开放
	if c.UserID <= 0 {
		select {
		case c.send <- mustJSON(outboundMessage{Type: "error", Data: map[string]interface{}{"msg": "unauthorized"}}):
		default:
		}
		return
	}

	msg := outboundMessage{
		Type: "chat",
		Data: map[string]interface{}{
			"text":     text,
			"userId":   c.UserID,
			"username": c.Username,
			"ts":       time.Now().Unix(),
		},
	}

	// 广播给所有“已鉴权连接”的用户
	b := mustJSON(msg)
	for client := range h.clients {
		if client.UserID <= 0 {
			continue
		}
		select {
		case client.send <- b:
		default:
		}
	}
}
