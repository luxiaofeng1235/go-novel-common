/*
 * @Descripttion: 用户相关路由（脚手架：游客/注册/登录）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api_routes

import (
	"go-novel/app/controller/api"
	"go-novel/middleware"

	"github.com/gin-gonic/gin"
)

func initUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userApi := new(api.User)
	user := r.Group("/user").Use(middleware.ApiReqDecrypt())
	{
		user.POST("/guest", userApi.Guest)                   //访客注册
		user.POST("/register", userApi.Register)             //注册
		user.POST("/login", userApi.Login)                   //登录
		user.GET("/info", middleware.ApiJwt(), userApi.Info) //获取用户信息

	}
	return r
}
