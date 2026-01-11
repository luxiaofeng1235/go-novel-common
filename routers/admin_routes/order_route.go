package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initOrderRoutes(r *gin.RouterGroup) gin.IRoutes {
	orderAdmin := new(admin.Order)
	feedback := r.Group("/order").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		feedback.GET("/orderList", orderAdmin.OrderList)
		feedback.GET("/editOrder", orderAdmin.UpdateOrder)
		feedback.POST("/updateOrder", orderAdmin.UpdateOrder)
	}

	return r
}
