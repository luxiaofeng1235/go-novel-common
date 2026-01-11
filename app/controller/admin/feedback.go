package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/feedback_service"
	"go-novel/utils"
	"html"
	"strconv"
)

type Feedback struct{}

func (feedback *Feedback) HelpList(c *gin.Context) {
	var req models.FeedbackHelpListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := feedback_service.HelpListSearch(&req)
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

func (feedback *Feedback) CreateHelp(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateHelpReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := feedback_service.CreateHelp(&req)
		if err != nil {
			utils.Fail(c, err, "创建信息失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	helpId, _ := strconv.Atoi(c.Query("id"))
	helpInfo, err := feedback_service.GetHelpById(int64(helpId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	helpInfo.Content = html.UnescapeString(helpInfo.Content)

	res := gin.H{
		"helpInfo": helpInfo,
	}
	utils.Success(c, res, "ok")
}

func (feedback *Feedback) UpdateHelp(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateHelpReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := feedback_service.UpdateHelp(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	helpId, _ := strconv.Atoi(c.Query("id"))
	helpInfo, err := feedback_service.GetHelpById(int64(helpId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	helpInfo.Content = html.UnescapeString(helpInfo.Content)

	res := gin.H{
		"helpInfo": helpInfo,
	}
	utils.Success(c, res, "ok")
}

func (feedback *Feedback) DeleteHelp(c *gin.Context) {
	var req models.DelHelpReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := feedback_service.DelHelp(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (feedback *Feedback) FeedbackList(c *gin.Context) {
	var req models.FeedBackListSearchReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := feedback_service.FeedBackList(&req)
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

func (feedback *Feedback) FeedbackReply(c *gin.Context) {
	var req models.FeedBackReplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	if err := feedback_service.FeedBackReply(&req); err != nil {
		utils.Fail(c, err, "审核失败")
		return
	}

	utils.Success(c, "", "")
}

func (feedback *Feedback) FeedbackBookList(c *gin.Context) {
	var req models.FeedbackBookListSearchReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := feedback_service.FeedBackBookListSearch(&req)
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

func (feedback *Feedback) UpdateFeedbackBook(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateFeedbackBookReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := feedback_service.UpdateFeedBackBook(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	feedbackBookId, _ := strconv.Atoi(c.Query("id"))
	feedbackBookInfo, err := feedback_service.GetFeedbackBookById(int64(feedbackBookId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"feedbackBookInfo": feedbackBookInfo,
	}
	utils.Success(c, res, "ok")
}

func (feedback *Feedback) DelFeedbackBook(c *gin.Context) {
	var req models.DelFeedbackBookReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := feedback_service.DelFeedbackBook(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
