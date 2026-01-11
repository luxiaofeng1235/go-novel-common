package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initSearchRoutes(r *gin.RouterGroup) gin.IRoutes {
	searchApi := new(api.Search)
	search := r.Group("/search").Use(middleware.ApiReqDecrypt())
	{
		search.POST("/searchHistory", middleware.ApiJwt(), searchApi.SearchHistory)
		search.POST("/hotSearchRank", searchApi.HotSearchRank)

		//GET请求
		search.GET("/searchHistory", middleware.ApiJwt(), searchApi.SearchHistory)
		search.GET("/hotSearchRank", searchApi.HotSearchRank)
	}
	return r
}
