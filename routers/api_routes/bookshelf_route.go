package api_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/api"
	"go-novel/middleware"
)

func initBookShelfRoutes(r *gin.RouterGroup) gin.IRoutes {
	bookshelfApi := new(api.Bookshelf)
	bookshelf := r.Group("/bookshelf").Use(middleware.ApiJwt(), middleware.ApiReqDecrypt())
	{
		bookshelf.POST("/book", bookshelfApi.Book)
		bookshelf.POST("/add", bookshelfApi.Add)
		bookshelf.POST("/del", bookshelfApi.Del)
		bookshelf.POST("/top", bookshelfApi.Top)
		bookshelf.POST("/isBookShelf", bookshelfApi.IsBookShelf)

		//get请求
		bookshelf.GET("/book", bookshelfApi.Book)
		bookshelf.GET("/add", bookshelfApi.Add)
		bookshelf.GET("/del", bookshelfApi.Del)
		bookshelf.GET("/top", bookshelfApi.Top)
		bookshelf.GET("/isBookShelf", bookshelfApi.IsBookShelf)

	}
	return r
}
