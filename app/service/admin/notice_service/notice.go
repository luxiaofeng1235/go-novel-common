package notice_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"html"
	"strings"
)

func GetNoticeById(id int64) (notice *models.McNotice, err error) {
	err = global.DB.Model(models.McNotice{}).Where("id", id).First(&notice).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func NoticeListSearch(req *models.NoticeListReq) (list []*models.McNotice, total int64, err error) {
	db := global.DB.Model(&models.McNotice{}).Order("id desc")

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

func CreateNotice(req *models.CreateNoticeReq) (InsertId int64, err error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		err = fmt.Errorf("%v", "公告标题不能为空")
		return
	}
	content := html.EscapeString(strings.TrimSpace(req.Content))
	if content == "" {
		err = fmt.Errorf("%v", "公告内容不能为空")
		return
	}
	link := strings.TrimSpace(req.Link)
	status := req.Status

	notice := models.McNotice{
		Title:   title,
		Content: content,
		Link:    link,
		Status:  status,
		Addtime: utils.GetUnix(),
	}

	if err = global.DB.Create(&notice).Error; err != nil {
		return
	}

	return notice.Id, nil
}

func UpdateNotice(req *models.UpdateNoticeReq) (res bool, err error) {
	id := req.NoticeId
	title := strings.TrimSpace(req.Title)
	if title == "" {
		err = fmt.Errorf("%v", "公告标题不能为空")
		return
	}
	content := html.EscapeString(strings.TrimSpace(req.Content))
	if content == "" {
		err = fmt.Errorf("%v", "公告内容不能为空")
		return
	}
	link := strings.TrimSpace(req.Link)
	status := req.Status
	var mapData = make(map[string]interface{})
	if title != "" {
		mapData["title"] = title
	}
	if content != "" {
		mapData["content"] = content
	}
	mapData["status"] = status
	mapData["link"] = link
	mapData["uptime"] = utils.GetUnix()

	if err = global.DB.Model(models.McNotice{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DelNotice(req *models.DelNoticeReq) (res bool, err error) {
	ids := req.NoticeIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McNotice{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
