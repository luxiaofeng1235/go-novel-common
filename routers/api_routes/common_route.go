package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initCommonRoutes(r *gin.RouterGroup) gin.IRoutes {
	commonApi := new(api.Common)
	common := r.Group("/common")
	{
		common.POST("/uploadImg", middleware.ApiJwt(), commonApi.UploadImg)      //上传图片
		common.GET("/uploadImg", middleware.ApiJwt(), commonApi.UploadImg)       //上传图片
		common.POST("/sendCode", middleware.ApiReqDecrypt(), commonApi.SendCode) //发送邮箱验证码+短信验证
		common.GET("/sendCode", middleware.ApiReqDecrypt(), commonApi.SendCode)  //发送邮箱验证码+短信验证
		common.GET("/baiduStatistics", commonApi.BaiduStatistics)                //百度统计
		common.GET("/chaojihuiStatistics", commonApi.ChaojihuiStatistics)        //神马平台统计
		common.POST("/vivoCallbackClk", commonApi.VivoCallBackClk)               //vivo 点击检测回传
		common.GET("/vivoCallbackClk", commonApi.VivoCallBackClk)                //vivo 点击检测回传
	}
	return r
}
