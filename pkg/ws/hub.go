/*
 * @Descripttion: WebSocket Hub（gorilla/websocket，基础广播能力）
 * @Author: red
 * @Date: 2026-01-13 11:25:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:25:00
 */
package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultWriteWait      = 10 * time.Second
	defaultPongWait       = 60 * time.Second
	defaultPingPeriod     = (defaultPongWait * 9) / 10
	defaultMaxMessageSize = 64 * 1024
)

type Hub struct {
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client, 128),
		unregister: make(chan *Client, 128),
		broadcast:  make(chan []byte, 512),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				_ = c.conn.Close()
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
					_ = c.conn.Close()
				}
			}
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
	case h.broadcast <- msg:
	default:
		// 队列满了就丢弃，避免阻塞业务
	}
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
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
		// 最小实现：收到什么就广播什么
		c.hub.Broadcast(message)
	}
}

func (c *Client) writePump() {
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
