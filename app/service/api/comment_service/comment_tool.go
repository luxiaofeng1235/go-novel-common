package comment_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetReportCommentIdsByUserId(uid int64) (cids []int64) {
	err := global.DB.Model(models.McCommentReport{}).Where("uid = ? and addtime >= ?", uid, utils.GetAgoDayUnix(30)).Pluck("cid", &cids).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetFollowCountByUid(userId, byUserId int64) (count int64) {
	err := global.DB.Model(models.McUserFollow{}).Where("uid = ? and by_uid = ?", userId, byUserId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookScoreById(id int64) (score float64) {
	var err error
	err = global.DB.Model(models.McBook{}).Select("score").Where("id = ?", id).Scan(&score).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getBookById(id int64) (book *models.McBook, err error) {
	err = global.DB.Model(models.McBook{}).Where("id = ?", id).Find(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getParentLinkByCommentId(commentId int64) (parentLink string) {
	var err error
	err = global.DB.Model(models.McComment{}).Select("parent_link").Where("id", commentId).Find(&parentLink).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getTextByCommentId(commentId int64) (text string) {
	var err error
	err = global.DB.Model(models.McComment{}).Select("text").Where("id", commentId).Find(&text).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
