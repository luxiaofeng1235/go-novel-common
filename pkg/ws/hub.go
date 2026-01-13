/*
 * @Descripttion: WebSocket Hub（gorilla/websocket，基础广播能力）
 * @Author: red
 * @Date: 2026-01-13 11:25:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 11:25:00
 */
package ws

import (
	"encoding/json"
	"strings"
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
	register      chan *Client
	unregister    chan *Client
	broadcastAll  chan []byte
	joinRoom      chan joinRequest
	broadcastRoom chan roomBroadcast
	chat          chan chatRequest
	dm            chan dmRequest

	clients map[*Client]struct{}
	rooms   map[string]map[*Client]struct{}
	users   map[int64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		register:      make(chan *Client, 128),
		unregister:    make(chan *Client, 128),
		broadcastAll:  make(chan []byte, 512),
		joinRoom:      make(chan joinRequest, 128),
		broadcastRoom: make(chan roomBroadcast, 512),
		chat:          make(chan chatRequest, 512),
		dm:            make(chan dmRequest, 512),
		clients:       make(map[*Client]struct{}),
		rooms:         make(map[string]map[*Client]struct{}),
		users:         make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
			h.addUserLocked(c)
			h.joinRoomLocked(c, defaultRoom())
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				h.removeUserLocked(c)
				h.leaveAllRoomsLocked(c)
				close(c.send)
				_ = c.conn.Close()
			}
		case msg := <-h.broadcastAll:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					h.leaveAllRoomsLocked(c)
					close(c.send)
					_ = c.conn.Close()
				}
			}
		case jr := <-h.joinRoom:
			h.joinRoomLocked(jr.client, jr.room)
			h.sendJoinOK(jr.client, jr.room)
			h.broadcastJoinEvent(jr.client, jr.room)
		case rb := <-h.broadcastRoom:
			room := normalizeRoom(rb.room)
			if room == "" {
				// fallback broadcast all
				for c := range h.clients {
					select {
					case c.send <- rb.msg:
					default:
						delete(h.clients, c)
						h.leaveAllRoomsLocked(c)
						close(c.send)
						_ = c.conn.Close()
					}
				}
				continue
			}
			for c := range h.rooms[room] {
				select {
				case c.send <- rb.msg:
				default:
					delete(h.clients, c)
					h.leaveAllRoomsLocked(c)
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
		// 队列满了就丢弃，避免阻塞业务
	}
}

func (h *Hub) BroadcastToRoom(room string, msg []byte) {
	if h == nil {
		return
	}
	room = normalizeRoom(room)
	if room == "" {
		h.Broadcast(msg)
		return
	}
	if len(msg) == 0 {
		return
	}
	select {
	case h.broadcastRoom <- roomBroadcast{room: room, msg: msg}:
	default:
	}
}

func (h *Hub) joinRoomLocked(c *Client, room string) {
	if h == nil || c == nil {
		return
	}
	room = normalizeRoom(room)
	if room == "" {
		room = defaultRoom()
	}
	h.leaveAllRoomsLocked(c)
	if _, ok := h.rooms[room]; !ok {
		h.rooms[room] = make(map[*Client]struct{})
	}
	h.rooms[room][c] = struct{}{}
}

func (h *Hub) leaveAllRoomsLocked(c *Client) {
	if h == nil || c == nil {
		return
	}
	for room, rm := range h.rooms {
		delete(rm, c)
		if len(rm) == 0 {
			delete(h.rooms, room)
		}
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
		c.handleMessage(message)
	}
}

type inboundMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type outboundMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
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
	case "join":
		var payload struct {
			Room string `json:"room"`
		}
		_ = json.Unmarshal(in.Data, &payload)
		c.joinRoom(payload.Room)
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
		// unknown -> ignore
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
	b, _ := json.Marshal(msg)
	select {
	case c.send <- b:
	default:
	}
}

func (c *Client) joinRoom(room string) {
	if c == nil || c.hub == nil {
		return
	}
	room = normalizeRoom(room)
	if room == "" {
		room = defaultRoom()
	}
	select {
	case c.hub.joinRoom <- joinRequest{client: c, room: room}:
	default:
	}
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

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
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

func defaultRoom() string { return "lobby" }

func normalizeRoom(room string) string {
	room = strings.TrimSpace(room)
	if room == "" {
		return ""
	}
	room = strings.TrimPrefix(room, "/")
	// 只保留简单安全字符
	room = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '_' || r == '-' || r == '.':
			return r
		default:
			return -1
		}
	}, room)
	if room == "" {
		return ""
	}
	if len(room) > 64 {
		room = room[:64]
	}
	return room
}

type joinRequest struct {
	client *Client
	room   string
}

type roomBroadcast struct {
	room string
	msg  []byte
}

type chatRequest struct {
	client *Client
	text   string
}

type dmRequest struct {
	client   *Client
	toUserID int64
	text     string
}

func (h *Hub) sendJoinOK(c *Client, room string) {
	if h == nil || c == nil {
		return
	}
	room = normalizeRoom(room)
	if room == "" {
		room = defaultRoom()
	}
	select {
	case c.send <- mustJSON(outboundMessage{Type: "join_ok", Data: map[string]interface{}{"room": room}}):
	default:
	}
}

func (h *Hub) broadcastJoinEvent(c *Client, room string) {
	if h == nil || c == nil {
		return
	}
	room = normalizeRoom(room)
	if room == "" {
		room = defaultRoom()
	}
	msg := outboundMessage{
		Type: "join",
		Data: map[string]interface{}{
			"room":     room,
			"userId":   c.UserID,
			"username": c.Username,
			"ts":       time.Now().Unix(),
		},
	}
	h.BroadcastToRoom(room, mustJSON(msg))
}

func (h *Hub) currentRoomLocked(c *Client) string {
	if h == nil || c == nil {
		return defaultRoom()
	}
	for room, rm := range h.rooms {
		if _, ok := rm[c]; ok {
			return room
		}
	}
	return defaultRoom()
}

func (h *Hub) broadcastChat(c *Client, text string) {
	if h == nil || c == nil {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	room := h.currentRoomLocked(c)
	msg := outboundMessage{
		Type: "chat",
		Data: map[string]interface{}{
			"text":     text,
			"room":     room,
			"userId":   c.UserID,
			"username": c.Username,
			"ts":       time.Now().Unix(),
		},
	}
	h.BroadcastToRoom(room, mustJSON(msg))
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
