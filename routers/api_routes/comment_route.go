package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initCommentRoutes(r *gin.RouterGroup) gin.IRoutes {
	commentApi := new(api.Comment)
	comment := r.Group("/comment").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		comment.POST("/list", commentApi.List)
		comment.POST("/add", commentApi.Add)
		comment.POST("/reply", commentApi.Reply)
		comment.POST("/del", commentApi.Del)
		comment.POST("/praise", commentApi.Praise)
		comment.POST("/starGroup", commentApi.StarGroup)
		comment.POST("/report", commentApi.Report)

		//GET请求
		comment.GET("/list", commentApi.List)
		comment.GET("/add", commentApi.Add)
		comment.GET("/reply", commentApi.Reply)
		comment.GET("/del", commentApi.Del)
		comment.GET("/praise", commentApi.Praise)
		comment.GET("/starGroup", commentApi.StarGroup)
		comment.GET("/report", commentApi.Report)

	}
	return r
}
