package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initVersionRoutes(r *gin.RouterGroup) gin.IRoutes {
	versionApi := new(api.Version)
	version := r.Group("/version").Use(middleware.ApiReqDecrypt())
	{
		version.POST("/getVersionInfo", middleware.ApiJwt(), versionApi.GetVersionInfo)
		version.POST("/getVersionNewInfo", middleware.ApiJwt(), versionApi.GetVersionNewInfo)

		//get请求
		version.GET("/getVersionInfo", middleware.ApiJwt(), versionApi.GetVersionInfo)
		version.GET("/getVersionNewInfo", middleware.ApiJwt(), versionApi.GetVersionNewInfo)

	}
	return r
}
