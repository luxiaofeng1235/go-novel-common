package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initVersionRoutes(r *gin.RouterGroup) gin.IRoutes {
	versionAdmin := new(admin.Version)
	book := r.Group("/version").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		book.GET("/versionList", versionAdmin.AppVersionList)
		book.GET("/editVersion", versionAdmin.UpdateAppVersion)
		book.POST("/updateVersion", versionAdmin.UpdateAppVersion)
		book.GET("/addVersion", versionAdmin.CreateAppVersion)     //添加应用-GET方式
		book.POST("/createVersion", versionAdmin.CreateAppVersion) //创建应用配置
		book.POST("/deleteVersion", versionAdmin.DeleteAppVersion) //删除引用-POST
	}
	return r
}
