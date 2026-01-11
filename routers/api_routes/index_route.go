package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
)

// 注册基础路由
func initIndexRoutes(r *gin.RouterGroup) gin.IRoutes {
	indexController := new(api.Index)

	r.GET("/", indexController.IndexGet)
	return r
}
