package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initWithdrawRoutes(r *gin.RouterGroup) gin.IRoutes {
	withdrawApi := new(api.Withdraw)
	withdraw := r.Group("/withdraw").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		withdraw.POST("/getWithdrawLimit", withdrawApi.GetWithdrawLimit)
		withdraw.POST("/AccountDetail", withdrawApi.AccountDetail)
		withdraw.POST("/accountSave", withdrawApi.AccountSave)
		withdraw.POST("/accountDel", withdrawApi.AccountDel)
		withdraw.POST("/apply", withdrawApi.Apply)
		withdraw.POST("/withdrawList", withdrawApi.WithdrawList)

		//get请求
		withdraw.GET("/getWithdrawLimit", withdrawApi.GetWithdrawLimit)
		withdraw.GET("/AccountDetail", withdrawApi.AccountDetail)
		withdraw.GET("/accountSave", withdrawApi.AccountSave)
		withdraw.GET("/accountDel", withdrawApi.AccountDel)
		withdraw.GET("/apply", withdrawApi.Apply)
		withdraw.GET("/withdrawList", withdrawApi.WithdrawList)
	}
	return r
}
