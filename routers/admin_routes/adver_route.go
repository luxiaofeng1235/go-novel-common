package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initAdverRoutes(r *gin.RouterGroup) gin.IRoutes {
	adverAdmin := new(admin.Adver)
	adver := r.Group("/adver").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		adver.GET("/adverList", adverAdmin.AdverList)      //广告列表
		adver.GET("/addAdver", adverAdmin.CreateAdver)     //获取广告的基础信息
		adver.POST("/addAdver", adverAdmin.CreateAdver)    //获取提交广告的设置信息
		adver.GET("/editAdver", adverAdmin.UpdateAdver)    //获取更新广告的每个对应的信息
		adver.POST("/updateAdver", adverAdmin.UpdateAdver) //获取更新广告的设置
		adver.POST("/deleteAdver", adverAdmin.DelAdver)    //删除当前广告
		//广告包相关接口****************************************
		adver.GET("/adverPackageList", adverAdmin.AdverPackageList)      //广告包管理列表
		adver.GET("/addAdverPackage", adverAdmin.CreateAdverPackage)     //添加广告包-GET获取
		adver.POST("/addAdverPackage", adverAdmin.CreateAdverPackage)    //添加广告包-POST获取
		adver.GET("/editAdverPackage", adverAdmin.UpdateAdverPackage)    //获取某个包信息
		adver.POST("/updateAdverPackage", adverAdmin.UpdateAdverPackage) //更新当前广告包信息
		adver.POST("/deleteAdverPackage", adverAdmin.DeleteAdverPackage) //删除当前广告包信息
		adver.GET("/getALlPackgeList", adverAdmin.GetAllPackageList)     //获取所有广告包信息

		//广告项目包管理
		adver.GET("/adverProjectList", adverAdmin.AdverProjectList)   //项目包列表
		adver.GET("/editAdverProject", adverAdmin.AddAdverProject)    //保存项目信息-GET方式
		adver.POST("/updateAdverProject", adverAdmin.AddAdverProject) //保存项目信息设置
		adver.GET("/projectTypeList", adverAdmin.ProjectTypeList)     //获取列表数据信息
	}
	return r
}
