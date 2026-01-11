package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initAdverRoutes(r *gin.RouterGroup) gin.IRoutes {
	adverApi := new(api.Adver)
	adver := r.Group("/adver").Use(middleware.ApiReqDecrypt())
	{
		adver.GET("/getAdverMap", adverApi.GetAdverMap)
		adver.POST("/getAdverMap", adverApi.GetAdverMap)                              //获取广告列表信息
		adver.POST("/updateClickCount", adverApi.UpdateClickCount)                    //获取第三方的广告统计
		adver.POST("/getAdvertInfoById", adverApi.GetAdvertInfoById)                  //获取广告的单挑配置信息
		adver.POST("/getNewbieStatus", middleware.ApiJwt(), adverApi.GetNewbieStatus) //获取新手保护期的生效时间
		adver.POST("/getAdverMapList", adverApi.GetAdverMapList)                      //获取所有的广告列表信息
		adver.POST("/getProjectInfo", adverApi.GetProjectInfo)                        //获取项目对应包的信息

		//get请求
		adver.GET("/updateClickCount", adverApi.UpdateClickCount)                    //获取第三方的广告统计
		adver.GET("/getAdvertInfoById", adverApi.GetAdvertInfoById)                  //获取广告的单挑配置信息
		adver.GET("/getNewbieStatus", middleware.ApiJwt(), adverApi.GetNewbieStatus) //获取新手保护期的生效时间
		adver.GET("/getAdverMapList", adverApi.GetAdverMapList)                      //获取所有的广告列表信息
		adver.GET("/getProjectInfo", adverApi.GetProjectInfo)                        //获取项目对应包的信息
	}
	return r
}
