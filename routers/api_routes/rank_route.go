package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initRankRoutes(r *gin.RouterGroup) gin.IRoutes {
	rankApi := new(api.Rank)
	rank := r.Group("/rank").Use(middleware.ApiReqDecrypt())
	{
		rank.POST("/rankList", rankApi.RankList)

		//get请求
		rank.GET("/rankList", rankApi.RankList)
	}
	return r
}
