package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initClassRoutes(r *gin.RouterGroup) gin.IRoutes {
	classAdmin := new(admin.Class)
	classType := r.Group("/class").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		classType.GET("/typeList", classAdmin.TypeList)
		classType.GET("/addType", classAdmin.CreateType)
		classType.POST("/addType", classAdmin.CreateType)
		classType.GET("/editType", classAdmin.UpdateType)
		classType.POST("/updateType", classAdmin.UpdateType)
		classType.POST("/deleteType", classAdmin.DelType)

		classType.GET("/classList", classAdmin.ClassList)
		classType.GET("/bookList", classAdmin.BookList)
		classType.GET("/addClass", classAdmin.CreateClass)
		classType.POST("/addClass", classAdmin.CreateClass)
		classType.GET("/editClass", classAdmin.UpdateClass)
		classType.POST("/updateClass", classAdmin.UpdateClass)
		classType.POST("/deleteClass", classAdmin.DelClass)
		classType.POST("/assignClass", classAdmin.AssignClass)
	}

	return r
}
