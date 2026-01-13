/*
 * @Descripttion: WebSocket 单聊（dm）
 * @Author: red
 * @Date: 2026-01-13 11:40:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:40:00
 */
package ws

import (
	"strings"
	"time"
)

type dmRequest struct {
	client   *Client
	toUserID int64
	text     string
}

func (c *Client) sendDM(toUserID int64, text string) {
	if c == nil || c.hub == nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	select {
	case c.hub.dm <- dmRequest{client: c, toUserID: toUserID, text: text}:
	default:
	}
}

func (h *Hub) addUserLocked(c *Client) {
	if h == nil || c == nil {
		return
	}
	if c.UserID <= 0 {
		return
	}
	if _, ok := h.users[c.UserID]; !ok {
		h.users[c.UserID] = make(map[*Client]struct{})
	}
	h.users[c.UserID][c] = struct{}{}
}

func (h *Hub) removeUserLocked(c *Client) {
	if h == nil || c == nil {
		return
	}
	if c.UserID <= 0 {
		return
	}
	if set, ok := h.users[c.UserID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.users, c.UserID)
		}
	}
}

func (h *Hub) sendDirectMessage(from *Client, toUserID int64, text string) {
	if h == nil || from == nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	if from.UserID <= 0 {
		select {
		case from.send <- mustJSON(outboundMessage{Type: "error", Data: map[string]interface{}{"msg": "unauthorized"}}):
		default:
		}
		return
	}
	if toUserID <= 0 {
		select {
		case from.send <- mustJSON(outboundMessage{Type: "error", Data: map[string]interface{}{"msg": "toUserId required"}}):
		default:
		}
		return
	}

	msg := outboundMessage{
		Type: "dm",
		Data: map[string]interface{}{
			"text":         text,
			"fromUserId":   from.UserID,
			"fromUsername": from.Username,
			"toUserId":     toUserID,
			"ts":           time.Now().Unix(),
		},
	}
	b := mustJSON(msg)

	// send to all connections of the target user
	if targets, ok := h.users[toUserID]; ok {
		for c := range targets {
			select {
			case c.send <- b:
			default:
			}
		}
	}
	// echo back to sender (so client can render sent message)
	select {
	case from.send <- b:
	default:
	}
}
