package feedback_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func FeedBackListSearch(req *models.FeedBackListSearchReq) (list []*models.McFeedback, total int64, err error) {
	db := global.DB.Model(&models.McFeedback{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	text := strings.TrimSpace(req.Text)
	if text != "" {
		db = db.Where("text like ?", "%"+text+"%")
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}
	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return
}

func FeedBackList(req *models.FeedBackListSearchReq) (feeds []*models.FeedBackListSearchRes, total int64, err error) {
	var list []*models.McFeedback
	list, total, err = FeedBackListSearch(req)
	if len(list) <= 0 {
		return
	}
	for _, val := range list {
		var pics []string
		if val.Pics != "" {
			pics = utils.GetAdminPic(val.Pics)
		}
		feed := &models.FeedBackListSearchRes{
			Id:        val.Id,
			Text:      val.Text,
			Pics:      pics,
			Status:    val.Status,
			Reply:     val.Reply,
			Ip:        val.Ip,
			UserId:    val.Uid,
			Phone:     val.Phone,
			Email:     val.Email,
			Addtime:   val.Addtime,
			Replytime: val.Replytime,
		}
		feeds = append(feeds, feed)
	}
	return
}

func FeedBackReply(req *models.FeedBackReplyReq) (err error) {
	if req.Status == 2 && req.Reply == "" {
		err = fmt.Errorf("%v", "处理内容不能为空")
		return
	}
	mapData := make(map[string]interface{})
	mapData["status"] = req.Status
	mapData["reply"] = req.Reply
	mapData["replytime"] = utils.GetUnix()
	if err = global.DB.Model(models.McFeedback{}).Where("id", req.FeedbackId).Updates(&mapData).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return nil
}
