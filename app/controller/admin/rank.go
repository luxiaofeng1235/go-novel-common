package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/rank_service"
	"go-novel/utils"
	"strconv"
)

type Rank struct{}

func (rank *Rank) RankList(c *gin.Context) {
	var req models.RankListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := rank_service.RankListSearch(&req)
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

func (rank *Rank) CreateRank(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateRankReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := rank_service.CreateRank(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (rank *Rank) UpdateRank(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateRankReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := rank_service.UpdateRank(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	rankId, _ := strconv.Atoi(c.Query("id"))
	rankInfo, err := rank_service.GetRankById(int64(rankId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"rankInfo": rankInfo,
	}
	utils.Success(c, res, "ok")
}

func (rank *Rank) DelRank(c *gin.Context) {
	var req models.DeleteRankReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := rank_service.DeleteRank(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
