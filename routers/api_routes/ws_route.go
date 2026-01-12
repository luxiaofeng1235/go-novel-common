/*
 * @Descripttion: WebSocket 路由（脚手架占位）
 * @Author: red
 * @Date: 2026-01-12 11:10:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
)

func initWsRoutes(r *gin.RouterGroup) gin.IRoutes {
	wsApi := new(api.Ws)
	ws := r.Group("/ws")
	{
		ws.GET("", wsApi.HandleRequest)
	}
	return r
}
