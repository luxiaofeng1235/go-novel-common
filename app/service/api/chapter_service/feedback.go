package chapter_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetFeedbackCountByUid(bookName, author, chapterName string, feedbackType int, userId int64) (count int64) {
	err := global.DB.Model(models.McBookFeedback{}).Where("book_name = ? and author >= ? and chapter_name = ? and feedback_type = ? and uid = ?", bookName, author, chapterName, feedbackType, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func FeedBackAdd(req *models.ChapterFeedBackAddReq) (err error) {
	bookName := strings.TrimSpace(req.BookName)
	author := strings.TrimSpace(req.Author)
	chapterName := strings.TrimSpace(req.ChapterName)
	feedbackType := req.FeedbackType
	bid := req.Bid
	cid := req.Cid
	ip := req.Ip
	userId := req.UserId

	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	if bookName == "" || author == "" {
		err = fmt.Errorf("%v", "反馈小说名称或作者不能为空")
		return
	}
	if chapterName == "" {
		err = fmt.Errorf("%v", "章节名称不能为空")
		return
	}
	count := GetFeedbackCountByUid(bookName, author, chapterName, feedbackType, userId)
	if count > 0 {
		err = fmt.Errorf("%v", "您已经反馈过啦")
		return
	}
	userinfo, _ := user_service.GetUserById(userId)
	feedback := models.McBookFeedback{
		BookName:     bookName,
		Author:       author,
		ChapterName:  chapterName,
		FeedbackType: feedbackType,
		Bid:          bid,
		Cid:          cid,
		Username:     userinfo.Username,
		Ip:           ip,
		Uid:          userId,
		Status:       0,
		Addtime:      utils.GetUnix(),
	}
	if err = global.DB.Create(&feedback).Error; err != nil {
		return
	}
	return
}
