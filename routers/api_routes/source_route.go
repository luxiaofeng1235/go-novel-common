package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initSourceRoutes(r *gin.RouterGroup) gin.IRoutes {
	sourceApi := new(api.Source)
	source := r.Group("/source").Use(middleware.ApiReqDecrypt())
	{
		source.POST("/sourceList", sourceApi.SourceList)

		//get请求
		source.GET("/sourceList", sourceApi.SourceList)
	}
	return r
}
