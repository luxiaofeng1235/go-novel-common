package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/shelf_service"
	"go-novel/utils"
)

type Bookshelf struct{}

func (bookshelf *Bookshelf) Book(c *gin.Context) {
	var req models.BookShelfListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, seconds, err := shelf_service.GetShelfList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list":    list,
		"total":   total,
		"seconds": seconds,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (bookshelf *Bookshelf) Add(c *gin.Context) {
	var req models.BookShelfAddReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	err = shelf_service.BookShelfAdd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "加入书架成功")
}

func (bookshelf *Bookshelf) Del(c *gin.Context) {
	var req models.BookShelfDelReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	err = shelf_service.BookFavDel(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "删除完成")
}

func (bookshelf *Bookshelf) Top(c *gin.Context) {
	var req models.BookShelfTopReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := shelf_service.BookShelfTop(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "ok")
	return
}

func (bookshelf *Bookshelf) IsBookShelf(c *gin.Context) {
	var req models.IsBookShelfReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	isExist, err := shelf_service.IsBookShelf(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"isExist": isExist,
	}
	utils.SuccessEncrypt(c, res, "ok")
	return
}
