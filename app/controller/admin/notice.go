package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/notice_service"
	"go-novel/utils"
	"html"
	"strconv"
)

type Notice struct{}

func (notice *Notice) NoticeList(c *gin.Context) {
	var req models.NoticeListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := notice_service.NoticeListSearch(&req)
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

func (notice *Notice) CreateNotice(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateNoticeReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := notice_service.CreateNotice(&req)
		if err != nil {
			utils.Fail(c, err, "创建信息失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	noticeId, _ := strconv.Atoi(c.Query("id"))
	noticeInfo, err := notice_service.GetNoticeById(int64(noticeId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	noticeInfo.Content = html.UnescapeString(noticeInfo.Content)

	res := gin.H{
		"noticeInfo": noticeInfo,
	}
	utils.Success(c, res, "ok")
}

func (notice *Notice) UpdateNotice(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateNoticeReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := notice_service.UpdateNotice(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	noticeId, _ := strconv.Atoi(c.Query("id"))
	noticeInfo, err := notice_service.GetNoticeById(int64(noticeId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	noticeInfo.Content = html.UnescapeString(noticeInfo.Content)

	res := gin.H{
		"noticeInfo": noticeInfo,
	}
	utils.Success(c, res, "ok")
}

func (notice *Notice) DeleteNotice(c *gin.Context) {
	var req models.DelNoticeReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := notice_service.DelNotice(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
