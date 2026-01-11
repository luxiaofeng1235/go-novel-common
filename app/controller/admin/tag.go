package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/tag_service"
	"go-novel/utils"
	"strconv"
)

type Tag struct{}

func (tag *Tag) TagList(c *gin.Context) {
	var req models.TagListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := tag_service.TagListSearch(&req)
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

func (tag *Tag) CreateTag(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateTagReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := tag_service.CreateTag(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (tag *Tag) UpdateTag(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateTagReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := tag_service.UpdateTag(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	tagId, _ := strconv.Atoi(c.Query("id"))
	tagInfo, err := tag_service.GetTagById(int64(tagId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"tagInfo": tagInfo,
	}
	utils.Success(c, res, "ok")
}

func (tag *Tag) DelTag(c *gin.Context) {
	var req models.DeleteTagReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := tag_service.DeleteTag(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (tag *Tag) AssignTag(c *gin.Context) {
	var req models.AssignTagReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	res, err := tag_service.AssignTag(&req)
	if err != nil || res == false {
		utils.Fail(c, err, "归类信息失败")
		return
	}

	utils.Success(c, "", "归类信息成功")
	return
}
