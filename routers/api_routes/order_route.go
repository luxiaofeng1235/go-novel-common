package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initOrderRoutes(r *gin.RouterGroup) gin.IRoutes {
	orderApi := new(api.Order)
	order := r.Group("/order")
	{
		order.POST("/createOrder", middleware.ApiJwt(), middleware.ApiReqDecrypt(), orderApi.CreateOrder)
		order.POST("/returnUrl", middleware.Logger(), orderApi.ReturnUrl)
		order.POST("/notifyUrlTest", middleware.Logger(), orderApi.NotifyUrlTest)
		order.POST("/notifyUrl", middleware.Logger(), orderApi.NotifyUrl)
		order.POST("/queryOrder", middleware.Logger(), middleware.ApiReqDecrypt(), orderApi.QueryOrder)

		//get请求
		order.GET("/createOrder", middleware.ApiJwt(), middleware.ApiReqDecrypt(), orderApi.CreateOrder)
		order.GET("/returnUrl", middleware.Logger(), orderApi.ReturnUrl)
		order.GET("/notifyUrlTest", middleware.Logger(), orderApi.NotifyUrlTest)
		order.GET("/notifyUrl", middleware.Logger(), orderApi.NotifyUrl)
		order.GET("/queryOrder", middleware.Logger(), middleware.ApiReqDecrypt(), orderApi.QueryOrder)

	}
	return r
}
