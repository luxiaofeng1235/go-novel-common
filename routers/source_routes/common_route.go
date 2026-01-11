package source_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/source"
)

func initCommonRoutes(r *gin.RouterGroup) gin.IRoutes {
	commonApi := new(source.Upload)
	api := r.Group("/common")
	{
		api.POST("/uploadBookPic", commonApi.UploadBookPic)
	}
	return r
}
