package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initSettingRoutes(r *gin.RouterGroup) gin.IRoutes {
	settingApi := new(api.Setting)
	setting := r.Group("/setting").Use(middleware.ApiReqDecrypt())
	{
		setting.POST("/getValue", settingApi.GetValue)
		setting.POST("/getAppConfigInfo", settingApi.GetAppConfigInfo) //获取配置信息
		//setting.POST("/getRanksName", settingApi.GetRanksName)

		//get请求
		setting.GET("/getValue", settingApi.GetValue)
		setting.GET("/getAppConfigInfo", settingApi.GetAppConfigInfo) //获取配置

	}
	return r
}
