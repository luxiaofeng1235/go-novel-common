/*
 * @Descripttion: 通用路由（上传等，默认不需要登录）
 * @Author: red
 * @Date: 2026-01-12 12:20:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 12:20:00
 */
package api_routes

import (
	"go-novel/app/controller/api"

	"github.com/gin-gonic/gin"
)

func initCommonRoutes(r *gin.RouterGroup) gin.IRoutes {
	commonApi := new(api.Common)
	common := r.Group("/common")
	{
		common.GET("/ping", commonApi.Ping)      //存活探针
		common.POST("/upload", commonApi.Upload) //上传
	}
	return r
}
