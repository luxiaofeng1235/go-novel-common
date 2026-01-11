package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

// 注册基础路由
func initIndexRoutes(r *gin.RouterGroup) gin.IRoutes {
	indexAdmin := new(admin.Index)

	r.GET("/", indexAdmin.IndexGet)
	r.GET("/captcha", indexAdmin.Captcha)
	//登录
	r.POST("/login", indexAdmin.Login)
	r.GET("/token", indexAdmin.GetToken)
	r.POST("/refreshToken", indexAdmin.RefreshToken)

	index := r.Group("/index").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		index.POST("/avatar", indexAdmin.Avatar)           //修改头像
		index.GET("/profile", indexAdmin.Profile)          //查看管理员基本资料
		index.POST("/updatePwd", indexAdmin.UpdatePwd)     //修改管理员登录密码
		index.POST("/editProfile", indexAdmin.EditProfile) //修改管理员信息

		index.GET("/getUserInfo", indexAdmin.GetUserInfo) //获取管理员详细信息
		index.GET("/getRouters", indexAdmin.GetRouters)   //获取路由信息
		//验证token测试
		index.GET("/check", indexAdmin.CheckToken) //校验测试
		index.GET("/retoken", indexAdmin.ReToken)  //刷新token验证信息
	}

	return r
}
