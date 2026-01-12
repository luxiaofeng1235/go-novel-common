/*
 * @Descripttion: WebSocket 控制器（脚手架占位）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/global"
)

type Ws struct{}

func (ws *Ws) HandleRequest(c *gin.Context) {
	_ = global.Ws.HandleRequest(c.Writer, c.Request)
}
