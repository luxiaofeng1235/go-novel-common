package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initVipRoutes(r *gin.RouterGroup) gin.IRoutes {
	vipApi := new(api.Vip)
	vip := r.Group("/vip").Use(middleware.ApiReqDecrypt())
	{
		vip.POST("/getVipMessage", middleware.ApiJwt(), vipApi.GetVipMessage)
		vip.POST("/getVipCardRmb", vipApi.GetVipCardRmb)
		vip.POST("/getVipCardCion", vipApi.GetVipCardCion)
		vip.POST("/vipBookStore", middleware.ApiJwt(), vipApi.VipBookStore)
		vip.POST("/vipBooks", vipApi.VipBooks)
		vip.POST("/exchangeVip", middleware.ApiJwt(), vipApi.ExchangeVip)

		//get请求
		vip.GET("/getVipMessage", middleware.ApiJwt(), vipApi.GetVipMessage)
		vip.GET("/getVipCardRmb", vipApi.GetVipCardRmb)
		vip.GET("/getVipCardCion", vipApi.GetVipCardCion)
		vip.GET("/vipBookStore", middleware.ApiJwt(), vipApi.VipBookStore)
		vip.GET("/vipBooks", vipApi.VipBooks)
		vip.GET("/exchangeVip", middleware.ApiJwt(), vipApi.ExchangeVip)

	}
	return r
}
