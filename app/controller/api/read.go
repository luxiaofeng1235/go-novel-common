package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/read_service"
	"go-novel/utils"
)

type Read struct{}

func (read *Read) ReadList(c *gin.Context) {
	var req models.BookReadListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	res, err := read_service.GetReadRes(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (read *Read) ReadAdd(c *gin.Context) {
	var req models.ReadAddReq
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
	err = read_service.ReadAdd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "记录成功")
}

func (read *Read) ReadInfo(c *gin.Context) {
	var req models.ReadInfoReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	readInfo, err := read_service.ReadInfo(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, readInfo, "ok")
}

func (read *Read) ReadDel(c *gin.Context) {
	var req models.ReadDelReq
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
	err = read_service.ReadDel(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "删除完成")
}

func (read *Read) BrowseList(c *gin.Context) {
	var req models.BrowseListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	res, err := read_service.GetBrowseRes(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (read *Read) BrowseDel(c *gin.Context) {
	var req models.BrowseDelReq
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
	err = read_service.BrowseDel(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "删除完成")
}
