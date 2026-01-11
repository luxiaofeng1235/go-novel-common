package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initTagRoutes(r *gin.RouterGroup) gin.IRoutes {
	tagAdmin := new(admin.Tag)
	tag := r.Group("/tag").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		tag.GET("/tagList", tagAdmin.TagList)
		tag.GET("/addTag", tagAdmin.CreateTag)
		tag.POST("/addTag", tagAdmin.CreateTag)
		tag.GET("/editTag", tagAdmin.UpdateTag)
		tag.POST("/updateTag", tagAdmin.UpdateTag)
		tag.POST("/deleteTag", tagAdmin.DelTag)
		tag.POST("/assignTag", tagAdmin.AssignTag)
	}
	return r
}
