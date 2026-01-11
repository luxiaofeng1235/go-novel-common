package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

// 注册商户路由
func initUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userAdmin := new(admin.User)
	auth := r.Group("/user").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		auth.GET("/userList", userAdmin.UserList)
		auth.GET("/editUser", userAdmin.UpdateUser)
		auth.POST("/editUser", userAdmin.UpdateUser)
		auth.GET("/detailUser", userAdmin.DetailUser)
		auth.POST("/deleteUser", userAdmin.DelUser)
		auth.GET("/cionChangeList", userAdmin.CionChangeList)
	}

	withdraweAdmin := new(admin.Withdraw)
	withdraw := r.Group("/withdraw").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		withdraw.GET("/withdrawAccountList", withdraweAdmin.WithdrawAccountList)
		withdraw.GET("/limitList", withdraweAdmin.LimitList)
		withdraw.GET("/addLimit", withdraweAdmin.CreateLimit)
		withdraw.POST("/addLimit", withdraweAdmin.CreateLimit)
		withdraw.GET("/editLimit", withdraweAdmin.UpdateLimit)
		withdraw.POST("/updateLimit", withdraweAdmin.UpdateLimit)
		withdraw.POST("/delLimit", withdraweAdmin.DelLimit)
		withdraw.GET("/withdrawList", withdraweAdmin.WithdrawList)
		withdraw.POST("/withdrawCheck", withdraweAdmin.WithdrawCheck)
	}

	return r
}
