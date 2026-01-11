package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/findbook_service"
	"go-novel/utils"
)

type Findbook struct{}

func (findbook *Findbook) FindbookList(c *gin.Context) {
	var req models.ApiFindbookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, err := findbook_service.List(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list":  list,
		"total": total,
	}

	utils.SuccessEncrypt(c, res, "获取列表成功")
}

func (findbook *Findbook) CreateFindBook(c *gin.Context) {
	var req models.CreateFindBookReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	bookId, err := findbook_service.CreateFindbook(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"book_id": bookId,
	}
	utils.SuccessEncrypt(c, res, "提交成功")
	return
}
