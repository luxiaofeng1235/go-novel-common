package feedback_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetFeedbackBookById(id int64) (feedback *models.McBookFeedback, err error) {
	err = global.DB.Model(models.McBookFeedback{}).Where("id", id).First(&feedback).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func FeedBackBookListSearch(req *models.FeedbackBookListSearchReq) (list []*models.McBookFeedback, total int64, err error) {
	db := global.DB.Model(&models.McBookFeedback{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	bookName := strings.TrimSpace(req.BookName)
	if bookName != "" {
		db = db.Where("book_name = ?", bookName)
	}

	author := strings.TrimSpace(req.Author)
	if author != "" {
		db = db.Where("author = ?", author)
	}

	chapterName := strings.TrimSpace(req.ChapterName)
	if chapterName != "" {
		db = db.Where("chapter_name = ?", chapterName)
	}

	feedbackType := strings.TrimSpace(req.FeedbackType)
	if feedbackType != "" {
		db = db.Where("feedback_type = ?", feedbackType)
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

func UpdateFeedBackBook(req *models.UpdateFeedbackBookReq) (res bool, err error) {
	id := req.FeedbackBookId
	bookName := strings.TrimSpace(req.BookName)
	author := strings.TrimSpace(req.Author)
	chapterName := strings.TrimSpace(req.ChapterName)
	feedbackType := req.FeedbackType
	status := req.Status
	mapData := make(map[string]interface{})
	if bookName != "" {
		mapData["book_name"] = bookName
	}
	if author != "" {
		mapData["author"] = author
	}
	if chapterName != "" {
		mapData["chapter_name"] = chapterName
	}
	mapData["feedback_type"] = feedbackType
	mapData["status"] = status
	if req.Status == 2 {
		mapData["handtime"] = utils.GetUnix()
	}
	if err = global.DB.Model(models.McBookFeedback{}).Debug().Where("id", id).Updates(mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DelFeedbackBook(req *models.DelFeedbackBookReq) (res bool, err error) {
	ids := req.FeedbackBookIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McBookFeedback{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
