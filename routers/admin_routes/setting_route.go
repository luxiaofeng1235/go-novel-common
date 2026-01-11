package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

// 注册设置路由
func initSettingRoutes(r *gin.RouterGroup) gin.IRoutes {
	settingAdmin := new(admin.Setting)
	setting := r.Group("/setting").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		setting.GET("/getInfo", settingAdmin.GetInfo)             //获取设置列表的信息
		setting.GET("/getInfoByName", settingAdmin.GetInfoByName) //根据当前的value去获取相关的参数
		setting.POST("/updateInfo", settingAdmin.UpdateInfo)
		setting.POST("/updateInfoOne", settingAdmin.UpdateInfoOne)
		setting.POST("/updateAgreementOne", settingAdmin.UpdateAgreementOne) //更新每个项目里的隐私协议
		setting.GET("/getAgreementInfo", settingAdmin.GetAgreementInfo)      //获取单条保存的列表信息
		setting.POST("/upImg", settingAdmin.UploadSettingImg)
		setting.GET("/settingPackageList", settingAdmin.SettingPackageList) //获取
	}

	return r
}
