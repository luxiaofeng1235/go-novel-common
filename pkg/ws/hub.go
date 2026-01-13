/*
 * @Descripttion: WebSocket Hub（gorilla/websocket，连接与主循环）
 * @Author: red
 * @Date: 2026-01-13 11:25:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:55:00
 */
package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 关键超时/限制：用于保活与断链检测（ping/pong）、避免超大消息占用内存。
	defaultWriteWait      = 10 * time.Second
	defaultPongWait       = 60 * time.Second
	defaultPingPeriod     = (defaultPongWait * 9) / 10
	defaultMaxMessageSize = 64 * 1024
)

type Hub struct {
	// Hub 主循环（Run）是单 goroutine：集中处理状态变更，避免 map 并发读写。
	register     chan *Client
	unregister   chan *Client
	broadcastAll chan []byte
	chat         chan chatRequest
	dm           chan dmRequest

	// clients：所有在线连接（包含未登录连接）；users：按 userID 归档（同账号多端在线）。
	clients map[*Client]struct{}
	users   map[int64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		register:     make(chan *Client, 128),
		unregister:   make(chan *Client, 128),
		broadcastAll: make(chan []byte, 512),
		chat:         make(chan chatRequest, 512),
		dm:           make(chan dmRequest, 512),
		clients:      make(map[*Client]struct{}),
		users:        make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
			h.addUserLocked(c)
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				h.removeUserLocked(c)
				close(c.send)
				_ = c.conn.Close()
			}
		case msg := <-h.broadcastAll:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
					_ = c.conn.Close()
				}
			}
		case cr := <-h.chat:
			h.broadcastChat(cr.client, cr.text)
		case dr := <-h.dm:
			h.sendDirectMessage(dr.client, dr.toUserID, dr.text)
		}
	}
}

func (h *Hub) Broadcast(msg []byte) {
	if h == nil {
		return
	}
	if len(msg) == 0 {
		return
	}
	select {
	case h.broadcastAll <- msg:
	default:
		// 队列满则丢弃，避免阻塞调用方。
	}
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	// send：Hub -> Client 的写队列；writePump 负责真正写入 WebSocket。
	send     chan []byte
	UserID   int64
	Username string
}

func NewClient(h *Hub, conn *websocket.Conn, userID int64, username string) *Client {
	return &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
	}
}

func (c *Client) Start() {
	if c == nil || c.hub == nil || c.conn == nil {
		return
	}
	c.hub.register <- c
	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
	}()

	// 读侧：限制消息大小 + 读超时；收到 pong 时续期 deadline。
	c.conn.SetReadLimit(defaultMaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(defaultPongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.handleMessage(message)
	}
}

func (c *Client) writePump() {
	// 写侧：发送业务消息 + 定时 ping 保活（触发客户端返回 pong）。
	ticker := time.NewTicker(defaultPingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(defaultWriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(defaultWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
