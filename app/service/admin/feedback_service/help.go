package feedback_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"html"
	"strings"
)

func GetHelpById(id int64) (help *models.McFeedbackHelp, err error) {
	err = global.DB.Model(models.McFeedbackHelp{}).Where("id", id).First(&help).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func HelpListSearch(req *models.FeedbackHelpListReq) (list []*models.McFeedbackHelp, total int64, err error) {
	db := global.DB.Model(&models.McFeedbackHelp{}).Order("id desc")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title = ?", title)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", req.BeginTime)
	}

	if req.EndTime != "" {
		db = db.Where("addtime <=?", req.EndTime)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err
}

func CreateHelp(req *models.CreateHelpReq) (InsertId int64, err error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		err = fmt.Errorf("%v", "帮助标题不能为空")
		return
	}
	content := html.EscapeString(strings.TrimSpace(req.Content))
	if content == "" {
		err = fmt.Errorf("%v", "帮助内容不能为空")
		return
	}
	help := models.McFeedbackHelp{
		Title:   title,
		Content: content,
		Addtime: utils.GetUnix(),
	}

	if err = global.DB.Create(&help).Error; err != nil {
		return
	}

	return help.Id, nil
}

func UpdateHelp(req *models.UpdateHelpReq) (res bool, err error) {
	id := req.HelpId
	title := strings.TrimSpace(req.Title)
	if title == "" {
		err = fmt.Errorf("%v", "帮助标题不能为空")
		return
	}
	content := html.EscapeString(strings.TrimSpace(req.Content))
	if content == "" {
		err = fmt.Errorf("%v", "帮助内容不能为空")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["title"] = title
	mapData["content"] = content
	mapData["uptime"] = utils.GetUnix()

	if err = global.DB.Model(models.McFeedbackHelp{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DelHelp(req *models.DelHelpReq) (res bool, err error) {
	ids := req.HelpIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McFeedbackHelp{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
