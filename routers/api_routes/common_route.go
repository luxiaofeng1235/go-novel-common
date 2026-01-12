/*
 * @Descripttion: 通用路由（上传等，默认不需要登录）
 * @Author: red
 * @Date: 2026-01-12 12:20:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 12:20:00
 */
package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
)

func initCommonRoutes(r *gin.RouterGroup) gin.IRoutes {
	commonApi := new(api.Common)
	common := r.Group("/common")
	{
		common.POST("/upload", commonApi.Upload)
	}
	return r
}
