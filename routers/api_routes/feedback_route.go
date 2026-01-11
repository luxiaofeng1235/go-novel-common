package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initFeedbackRoutes(r *gin.RouterGroup) gin.IRoutes {
	feedbackApi := new(api.Feedback)
	feedback := r.Group("/feedback").Use(middleware.ApiReqDecrypt())
	{
		feedback.POST("/helpList", feedbackApi.HelpList)
		feedback.POST("/helpDetail", feedbackApi.HelpDetail)
		feedback.POST("/list", middleware.ApiJwt(), feedbackApi.List)
		feedback.POST("/add", middleware.ApiJwt(), feedbackApi.Add)

		//get请求
		feedback.GET("/helpList", feedbackApi.HelpList)
		feedback.GET("/helpDetail", feedbackApi.HelpDetail)
		feedback.GET("/list", middleware.ApiJwt(), feedbackApi.List)
		feedback.GET("/add", middleware.ApiJwt(), feedbackApi.Add)

	}
	return r
}
