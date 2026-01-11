package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/class_service"
	"go-novel/app/service/admin/collect_service"
	"go-novel/utils"
	"strconv"
)

type Collect struct{}

func (collect *Collect) CollectList(c *gin.Context) {
	var req models.CollectListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := collect_service.CollectListSearch(&req)
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

func (collect *Collect) CreateCollect(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateCollectReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := collect_service.CreateCollect(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	classList, err := class_service.GetClassBySex()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}
	res := gin.H{
		"classList": classList,
	}
	utils.Success(c, res, "ok")
}

func (collect *Collect) UpdateCollect(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateCollectReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := collect_service.UpdateCollect(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	collectId, _ := strconv.Atoi(c.Query("id"))
	collectInfo, err := collect_service.GetCollectResById(int64(collectId))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}

	classList, err := class_service.GetClassBySex()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}

	res := gin.H{
		"collectInfo": collectInfo,
		"classList":   classList,
	}
	utils.Success(c, res, "ok")
}

func (collect *Collect) DelCollect(c *gin.Context) {
	var req models.DeleteCollectReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := collect_service.DeleteCollect(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
