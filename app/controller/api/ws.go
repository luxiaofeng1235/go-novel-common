/*
 * @Descripttion: WebSocket 控制器（脚手架占位）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-novel/global"
	"go-novel/utils"
)

type Ws struct{}

func (wsApi *Ws) HandleRequest(c *gin.Context) {
	if global.WsHubManager == nil {
		utils.FailEncrypt(c, nil, "ws未初始化")
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token = strings.TrimSpace(auth[len("bearer "):])
		}
	}

	var userID int64
	var username string
	if token != "" {
		if claims, err := utils.ParseToken(token); err == nil && claims != nil {
			userID = claims.UserID
			username = strings.TrimSpace(claims.Username)
		}
	}

	// 使用 HubManager 自动路由到对应分片
	_ = global.WsHubManager.RegisterClient(conn, userID, username)
}
