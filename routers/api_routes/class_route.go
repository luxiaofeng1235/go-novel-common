package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initClassRoutes(r *gin.RouterGroup) gin.IRoutes {
	classApi := new(api.Class)
	class := r.Group("/class").Use(middleware.ApiReqDecrypt())
	{
		class.POST("/list", middleware.ApiJwt(), classApi.List)

		//get请求
		class.GET("/list", middleware.ApiJwt(), classApi.List)
	}
	return r
}
