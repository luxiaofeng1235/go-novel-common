package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

// 注册权限路由
func initAuthRoutes(r *gin.RouterGroup) gin.IRoutes {
	adminAdmin := new(admin.Admin)
	roleAdmin := new(admin.Role)
	menuAdmin := new(admin.Menu)

	auth := r.Group("/auth").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		//管理员管理
		auth.GET("/adminList", adminAdmin.AdminList)
		auth.GET("/addAdmin", adminAdmin.CreateAdmin)
		auth.POST("/addAdmin", adminAdmin.CreateAdmin)
		auth.GET("/editAdmin", adminAdmin.UpdateAdmin)
		auth.POST("/editAdmin", adminAdmin.UpdateAdmin)
		auth.POST("/deleteAdmin", adminAdmin.DeleteAdmin)
		auth.POST("/resetAdminPwd", adminAdmin.ResetAdminPwd)
		auth.POST("/changeAdminStatus", adminAdmin.ChangeAdminStatus)

		//角色管理
		auth.GET("/roleList", roleAdmin.RoleList)
		auth.GET("/addRole", roleAdmin.CreateRole)
		auth.POST("/addRole", roleAdmin.CreateRole)
		auth.GET("/editRole", roleAdmin.UpdateRole)
		auth.POST("/editRole", roleAdmin.UpdateRole)
		auth.POST("/deleteRole", roleAdmin.DeleteRole)
		auth.POST("/changeRoleStatus", roleAdmin.ChangeRoleStatus)

		//菜单管理
		auth.GET("/menuList", menuAdmin.MenuList)
		auth.GET("/addMenu", menuAdmin.CreateMenu)
		auth.POST("/addMenu", menuAdmin.CreateMenu)
		auth.GET("/editMenu", menuAdmin.UpdateMenu)
		auth.POST("/editMenu", menuAdmin.UpdateMenu)
		auth.POST("/deleteMenu", menuAdmin.DeleteMenu)
	}

	return r
}
