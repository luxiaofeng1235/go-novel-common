/*
 * @Descripttion: 用户相关路由（脚手架：游客/注册/登录）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userApi := new(api.User)
	user := r.Group("/user").Use(middleware.ApiReqDecrypt())
	{
		user.POST("/guest", userApi.Guest)
		user.POST("/register", userApi.Register)
		user.POST("/login", userApi.Login)

		//get请求
		user.GET("/guest", userApi.Guest)
		user.GET("/register", userApi.Register)
		user.GET("/login", userApi.Login)

	}
	return r
}
