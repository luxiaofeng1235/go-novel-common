package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initCollectRoutes(r *gin.RouterGroup) gin.IRoutes {
	collectAdmin := new(admin.Collect)
	collect := r.Group("/collect").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		collect.GET("/collectList", collectAdmin.CollectList)
		collect.GET("/addCollect", collectAdmin.CreateCollect)
		collect.POST("/addCollect", collectAdmin.CreateCollect)
		collect.GET("/editCollect", collectAdmin.UpdateCollect)
		collect.POST("/updateCollect", collectAdmin.UpdateCollect)
		collect.POST("/deleteCollect", collectAdmin.DelCollect)
	}
	return r
}
