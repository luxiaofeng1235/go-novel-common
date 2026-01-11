package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initChapterRoutes(r *gin.RouterGroup) gin.IRoutes {
	chapterApi := new(api.Chapter)
	book := r.Group("/chapter").Use(middleware.ApiReqDecrypt())
	{
		book.POST("/feedbackAdd", middleware.ApiJwt(), chapterApi.FeedbackAdd)

		book.GET("/feedbackAdd", middleware.ApiJwt(), chapterApi.FeedbackAdd)
	}
	return r
}
