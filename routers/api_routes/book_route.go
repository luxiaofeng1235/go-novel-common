package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initBookRoutes(r *gin.RouterGroup) gin.IRoutes {
	bookApi := new(api.Book)
	book := r.Group("/book").Use(middleware.ApiReqDecrypt())
	{
		book.POST("/bookCommentUserList", bookApi.BookCommentUserList) //用户列表
		book.POST("/BookCommentRankList", bookApi.BookCommentRankList) //书评的排行列表
		book.POST("/BookWonderfulList", bookApi.BookWonderfulList)     //精彩书评列表
		book.POST("/BookCommentByUserId", bookApi.BookCommentByUserId) //根据用户ID获取书评信息
		book.POST("/bookCommentInfo", bookApi.BookCommentInfo)         //获取书评基础信息

		book.POST("/list", middleware.ApiJwt(), bookApi.List)
		book.POST("/info", middleware.ApiJwt(), bookApi.Info)
		book.POST("/getHighScoreBook", middleware.ApiJwt(), bookApi.GetHighScoreBook)
		book.POST("/chapter", bookApi.Chapter)
		book.POST("/read", middleware.ApiJwt(), bookApi.Read)
		book.POST("/rankList", middleware.ApiJwt(), bookApi.RankList)
		book.POST("/getSectionForYouRec", middleware.ApiJwt(), bookApi.GetSectionForYouRec)
		book.POST("/getSectionHighScore", middleware.ApiJwt(), bookApi.GetSectionHighScore)
		book.POST("/getSectionEnd", bookApi.GetSectionEnd) //查询完结的书
		book.POST("/getSectionNew", bookApi.GetSectionNew) //查询新书的书
		book.POST("/getTags", middleware.ApiJwt(), bookApi.GetTags)
		book.POST("/teenZoneList", bookApi.TeenZoneList)
		book.POST("/getNewBookRec", middleware.ApiJwt(), bookApi.GetNewBookRec)
		book.POST("/getNewBookList", bookApi.GetNewBookList)
		book.POST("/todayUpdateBooks", bookApi.TodayUpdateBooks)
		book.POST("/getHotCount", bookApi.GetHotCount)
		book.POST("/test", bookApi.Test)
		book.POST("/getCateBookList", bookApi.GetCateBookList) //根据分类获取书籍

		//get请求
		book.GET("/bookCommentUserList", bookApi.BookCommentUserList) //用户列表
		book.GET("/BookCommentRankList", bookApi.BookCommentRankList) //书评的排行列表
		book.GET("/BookWonderfulList", bookApi.BookWonderfulList)     //精彩书评列表
		book.GET("/BookCommentByUserId", bookApi.BookCommentByUserId) //根据用户ID获取书评信息
		book.GET("/bookCommentInfo", bookApi.BookCommentInfo)         //获取书评基础信息

		book.GET("/list", middleware.ApiJwt(), bookApi.List)
		book.GET("/info", middleware.ApiJwt(), bookApi.Info)
		book.GET("/getHighScoreBook", middleware.ApiJwt(), bookApi.GetHighScoreBook)
		book.GET("/chapter", bookApi.Chapter)
		book.GET("/read", middleware.ApiJwt(), bookApi.Read)
		book.GET("/rankList", middleware.ApiJwt(), bookApi.RankList)
		book.GET("/getSectionForYouRec", middleware.ApiJwt(), bookApi.GetSectionForYouRec)
		book.GET("/getSectionHighScore", middleware.ApiJwt(), bookApi.GetSectionHighScore)
		book.GET("/getSectionEnd", bookApi.GetSectionEnd) //查询完结的书
		book.GET("/getSectionNew", bookApi.GetSectionNew) //查询新书的书
		book.GET("/getTags", middleware.ApiJwt(), bookApi.GetTags)
		book.GET("/teenZoneList", bookApi.TeenZoneList)
		book.GET("/getNewBookRec", middleware.ApiJwt(), bookApi.GetNewBookRec)
		book.GET("/getNewBookList", bookApi.GetNewBookList)
		book.GET("/todayUpdateBooks", bookApi.TodayUpdateBooks)
		book.GET("/getHotCount", bookApi.GetHotCount)
		book.GET("/test", bookApi.Test)
		book.GET("/getCateBookList", bookApi.GetCateBookList) //根据分类获取书籍
	}
	return r
}
