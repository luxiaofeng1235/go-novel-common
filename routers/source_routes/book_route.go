package source_routes

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/controller/source"
)

func initBookRoutes(r *gin.RouterGroup) gin.IRoutes {
	bookApi := new(source.Book)
	book := r.Group("/book")
	{
		book.POST("/getBookInfo", bookApi.GetBookInfo)
		book.POST("/updateBookInfo", bookApi.UpdateBookInfo)
	}
	return r
}
