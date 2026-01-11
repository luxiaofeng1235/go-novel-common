package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/checkin_service"
	"go-novel/utils"
)

type Checkin struct{}

func (checkin *Checkin) CheckinList(c *gin.Context) {
	var req models.CheckinListReq
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
	checkins, err := checkin_service.GetCheckinList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, checkins, "获取签到列表成功")
}

func (checkin *Checkin) Checkin(c *gin.Context) {
	var req models.CheckinReq
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
	cion, vip, err := checkin_service.GetCheckin(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"cion": cion,
		"vip":  vip,
	}
	utils.SuccessEncrypt(c, res, "签到成功")
}

func (checkin *Checkin) CheckinHistory(c *gin.Context) {
	var req models.CheckinHistoryReq
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
	historys, err := checkin_service.CheckinHistory(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"historys": historys,
	}
	utils.SuccessEncrypt(c, res, "获取数据成功")
}

func (checkin *Checkin) OpenCheckinRemind(c *gin.Context) {
	var req models.OpenCheckinRemindReq
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
	err = checkin_service.OpenCheckinRemind(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "设置成功")
}
