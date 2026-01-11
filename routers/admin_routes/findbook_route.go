package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initFindBookRoutes(r *gin.RouterGroup) gin.IRoutes {
	findbookAdmin := new(admin.FindBook)
	findbook := r.Group("/findbook").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		findbook.GET("/findbookList", findbookAdmin.FindBookList)
		findbook.GET("/editFindbook", findbookAdmin.UpdateFindBook)
		findbook.POST("/updateFindbook", findbookAdmin.UpdateFindBook)
		findbook.POST("/deleteFindbook", findbookAdmin.DeleteFindBook)
	}
	return r
}
