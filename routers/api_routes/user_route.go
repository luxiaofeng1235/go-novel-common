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
		user.POST("/login", userApi.Login)
		user.POST("/logoff", userApi.Logoff)
		user.POST("/info", middleware.ApiJwt(), userApi.Info)
		user.POST("/edit", middleware.ApiJwt(), userApi.Edit)
		user.POST("/follow", middleware.ApiJwt(), userApi.Follow)
		user.POST("/followList", middleware.ApiJwt(), userApi.FollowList)
		user.POST("/bindRegistId", middleware.ApiJwt(), userApi.BindRegistId)
		user.POST("/myInvitRewards", middleware.ApiJwt(), userApi.MyInvitRewards)

		//get请求
		user.GET("/guest", userApi.Guest)
		user.GET("/login", userApi.Login)
		user.GET("/logoff", userApi.Logoff)
		user.GET("/info", middleware.ApiJwt(), userApi.Info)
		user.GET("/edit", middleware.ApiJwt(), userApi.Edit)
		user.GET("/follow", middleware.ApiJwt(), userApi.Follow)
		user.GET("/followList", middleware.ApiJwt(), userApi.FollowList)
		user.GET("/bindRegistId", middleware.ApiJwt(), userApi.BindRegistId)
		user.GET("/myInvitRewards", middleware.ApiJwt(), userApi.MyInvitRewards)

	}
	return r
}
