package admin_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/admin"
	"go-novel/middleware"
)

func initBookRoutes(r *gin.RouterGroup) gin.IRoutes {
	bookAdmin := new(admin.Book)
	book := r.Group("/book").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		book.GET("/bookList", bookAdmin.BookList)
		book.GET("/bookPassList", bookAdmin.BookPassList)     //审核通过的书籍列表
		book.GET("/getBookRecList", bookAdmin.GetBookRecList) //获取推荐的热门、最新、完结等书籍ID
		book.POST("/setRecBookIndex", bookAdmin.SetRecBookIndex)
		book.GET("/setRecBookIndex", bookAdmin.SetRecBookIndex)
		book.GET("/addBook", bookAdmin.CreateBook)
		book.POST("/addBook", bookAdmin.CreateBook)
		book.GET("/editBook", bookAdmin.UpdateBook)
		book.POST("/updateBook", bookAdmin.UpdateBook)
		book.POST("/deleteBook", bookAdmin.DelBook)
		book.GET("/detailBook", bookAdmin.DetailBook)
		book.GET("/getRandHit", bookAdmin.GetRandHit)
	}

	chapterAdmin := new(admin.Chapter)
	chapter := r.Group("/chapter").Use(middleware.AdminJwt()).Use(middleware.Auth())
	{
		chapter.GET("/chapterList", chapterAdmin.ChapterList)
		chapter.GET("/addChapter", chapterAdmin.CreateChapter)
		chapter.POST("/addChapter", chapterAdmin.CreateChapter)
		chapter.GET("/editChapter", chapterAdmin.UpdateChapter)
		chapter.POST("/updateChapter", chapterAdmin.UpdateChapter)
		chapter.POST("/deleteChapter", chapterAdmin.DelChapter)
	}
	return r
}
