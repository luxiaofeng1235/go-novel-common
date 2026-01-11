package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initReadRoutes(r *gin.RouterGroup) gin.IRoutes {
	readApi := new(api.Read)
	read := r.Group("/read").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		read.POST("/readList", readApi.ReadList)
		read.POST("/readAdd", readApi.ReadAdd)
		read.POST("/readInfo", readApi.ReadInfo)
		read.POST("/readDel", readApi.ReadDel)
		read.POST("/browseList", readApi.BrowseList)
		read.POST("/browseDel", readApi.BrowseDel)

		//get请求
		read.GET("/readList", readApi.ReadList)
		read.GET("/readAdd", readApi.ReadAdd)
		read.GET("/readInfo", readApi.ReadInfo)
		read.GET("/readDel", readApi.ReadDel)
		read.GET("/browseList", readApi.BrowseList)
		read.GET("/browseDel", readApi.BrowseDel)
	}
	return r
}
