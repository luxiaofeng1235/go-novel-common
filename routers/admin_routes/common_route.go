package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

// 注册基础路由
func initCommonRoutes(r *gin.RouterGroup) gin.IRoutes {
	uploadAdmin := new(admin.Upload)
	commonApi := new(admin.Common)
	common := r.Group("/common")
	{
		common.POST("/upload", uploadAdmin.Upload)
		common.POST("/uploadClassPic", middleware.ApiJwt(), commonApi.UploadClassPic)
		common.POST("/uploadBookPic", middleware.ApiJwt(), commonApi.UploadBookPic)
		common.POST("/uploadAdverPic", middleware.ApiJwt(), commonApi.UploadAdverPic)
		common.POST("/uploadAPkFile", middleware.ApiJwt(), commonApi.UploadAPkFile)
	}
	return r
}
