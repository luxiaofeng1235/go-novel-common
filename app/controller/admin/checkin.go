package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/checkin_service"
	"go-novel/utils"
	"strconv"
)

type Checkin struct{}

func (checkin *Checkin) RewardList(c *gin.Context) {
	var req models.CheckinRewardListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := checkin_service.RewardListSearch(&req)
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

func (checkin *Checkin) CreateReward(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateRewardReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := checkin_service.CreateReward(&req)
		if err != nil {
			utils.Fail(c, err, "创建奖励失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (checkin *Checkin) UpdateReward(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateRewardReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := checkin_service.UpdateReward(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	rewardId, _ := strconv.Atoi(c.Query("id"))
	rewardInfo, err := checkin_service.GetRewardById(int64(rewardId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"rewardInfo": rewardInfo,
	}
	utils.Success(c, res, "ok")
}

func (checkin *Checkin) DelReward(c *gin.Context) {
	var req models.DeleteRewardReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := checkin_service.DeleteReward(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (checkin *Checkin) CheckinList(c *gin.Context) {
	var req models.CheckinListSearchReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := checkin_service.CheckinListSearch(&req)
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
