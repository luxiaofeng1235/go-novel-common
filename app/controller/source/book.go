package source

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/source/book_service"
	"go-novel/utils"
)

type Book struct{}

func (book *Book) GetBookInfo(c *gin.Context) {
	var req models.SourceGetBookInfoReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}

	res, err := book_service.GetBookInfo(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}

	utils.Success(c, res, "ok")
}

func (book *Book) UpdateBookInfo(c *gin.Context) {
	var req models.SourceUpdateBookInfoReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}

	err := book_service.UpdateBookChapter(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}

	utils.Success(c, "", "ok")
}
