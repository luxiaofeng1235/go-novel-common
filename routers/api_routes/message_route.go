package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initMessageRoutes(r *gin.RouterGroup) gin.IRoutes {
	messageApi := new(api.Message)
	wsApi := new(api.Ws)
	message := r.Group("/message")
	{
		message.GET("/ws", wsApi.HandleRequest)
		message.POST("/getLastNotice", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.GetLastNotice)
		message.POST("/messageList", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.MessageList)
		message.POST("/replyList", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.ReplyList)
		message.POST("/praiseList", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.PraiseList)
		message.POST("/updateIsRead", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.UpdateIsRead)

		//get请求
		message.GET("/getLastNotice", middleware.ApiReqDecrypt(), messageApi.GetLastNotice)
		message.GET("/messageList", middleware.ApiReqDecrypt(), messageApi.MessageList)
		message.GET("/replyList", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.ReplyList)
		message.GET("/praiseList", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.PraiseList)
		message.GET("/updateIsRead", middleware.ApiJwt(), middleware.ApiReqDecrypt(), messageApi.UpdateIsRead)

	}
	return r
}
