package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
)

func initBiqugeRoutes(r *gin.RouterGroup) gin.IRoutes {
	biqugeApi := new(api.Biquge)
	biquge := r.Group("/biquge")
	{
		//biquge.GET("/generalDecrypt", biqugeApi.GeneralDecrypt) //通用解密，包含搜索和列表返回的解密
		biquge.POST("/generalDecrypt", biqugeApi.GeneralDecrypt)
		biquge.POST("/chapterDecrypt", biqugeApi.ChapterDecrypt) //章节内容解密，单独对章节内容做解密，这个地方可以直接用API来实现 --请求章节目录专用
		biquge.POST("/contentDecrypt", biqugeApi.ContentDecrypt) //解析内容，只对内容做单独的解析，没有其他用途

		//get请求
		biquge.GET("/generalDecrypt", biqugeApi.GeneralDecrypt)
		biquge.GET("/chapterDecrypt", biqugeApi.ChapterDecrypt) //章节内容解密，单独对章节内容做解密，这个地方可以直接用API来实现 --请求章节目录专用
		biquge.GET("/contentDecrypt", biqugeApi.ContentDecrypt) //解析内容，只对内容做单独的解析，没有其他用途
	}
	return r
}
