package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initRankRoutes(r *gin.RouterGroup) gin.IRoutes {
	rankAdmin := new(admin.Rank)
	tag := r.Group("/rank").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		tag.GET("/rankList", rankAdmin.RankList)
		tag.GET("/addRank", rankAdmin.CreateRank)
		tag.POST("/addRank", rankAdmin.CreateRank)
		tag.GET("/editRank", rankAdmin.UpdateRank)
		tag.POST("/updateRank", rankAdmin.UpdateRank)
		tag.POST("/deleteRank", rankAdmin.DelRank)
	}
	return r
}
