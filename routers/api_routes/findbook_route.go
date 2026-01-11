package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initFindBookRoutes(r *gin.RouterGroup) gin.IRoutes {
	feedbackApi := new(api.Findbook)
	feedback := r.Group("/findbook").Use(middleware.ApiReqDecrypt())
	{
		feedback.POST("/findbookList", middleware.ApiJwt(), feedbackApi.FindbookList)
		feedback.POST("/createFindBook", middleware.ApiJwt(), feedbackApi.CreateFindBook)

		//get请求
		feedback.GET("/findbookList", middleware.ApiJwt(), feedbackApi.FindbookList)
		feedback.GET("/createFindBook", middleware.ApiJwt(), feedbackApi.CreateFindBook)

	}
	return r
}
