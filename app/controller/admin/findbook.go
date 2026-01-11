package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/findbook_service"
	"go-novel/utils"
	"strconv"
)

type FindBook struct{}

func (findbook *FindBook) FindBookList(c *gin.Context) {
	var req models.FindBookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := findbook_service.FindBookListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}
	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (findbook *FindBook) UpdateFindBook(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateFindBookReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := findbook_service.UpdateFindBook(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	findBookId, _ := strconv.Atoi(c.Query("id"))
	findBookInfo, err := findbook_service.GetFindBookById(int64(findBookId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"findbookInfo": findBookInfo,
	}
	utils.Success(c, res, "ok")
}

func (findbook *FindBook) DeleteFindBook(c *gin.Context) {
	var req models.DelFindBookReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := findbook_service.DelFindBook(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
