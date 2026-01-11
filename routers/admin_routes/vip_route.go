package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initVipRoutes(r *gin.RouterGroup) gin.IRoutes {
	vipAdmin := new(admin.Vip)
	vip := r.Group("/vip").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		vip.GET("/cardList", vipAdmin.CardList)
		vip.GET("/addCard", vipAdmin.CreateCard)
		vip.POST("/addCard", vipAdmin.CreateCard)
		vip.GET("/editCard", vipAdmin.UpdateCard)
		vip.POST("/updateCard", vipAdmin.UpdateCard)
		vip.POST("/deleteCard", vipAdmin.DelVipCard)
	}
	return r
}
