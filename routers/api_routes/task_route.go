package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initTaskRoutes(r *gin.RouterGroup) gin.IRoutes {
	taskApi := new(api.Task)
	task := r.Group("/task").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		task.POST("/list", taskApi.List)
		task.POST("/receive", taskApi.Receive)
		task.POST("/share", taskApi.Share)
		task.POST("/cionChangeList", taskApi.CionChangeList)
		task.POST("/completeMsgPush", taskApi.CompleteMsgPush)
		task.POST("/completeVideoReward", taskApi.CompleteVideoReward)

		//get请求
		task.GET("/list", taskApi.List)
		task.GET("/receive", taskApi.Receive)
		task.GET("/share", taskApi.Share)
		task.GET("/cionChangeList", taskApi.CionChangeList)
		task.GET("/completeMsgPush", taskApi.CompleteMsgPush)
		task.GET("/completeVideoReward", taskApi.CompleteVideoReward)

	}

	checkinApi := new(api.Checkin)
	checkin := r.Group("/checkin").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		checkin.POST("/list", checkinApi.CheckinList)
		checkin.POST("/index", checkinApi.Checkin)
		checkin.POST("/history", checkinApi.CheckinHistory)
		checkin.POST("/openCheckinRemind", checkinApi.OpenCheckinRemind)
	}
	return r
}
