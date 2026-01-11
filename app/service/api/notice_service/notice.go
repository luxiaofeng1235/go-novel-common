package notice_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetNotifyCommentIdsByUserId(uid int64, notifyType string) (cids []int64) {
	err := global.DB.Model(models.McNotify{}).Where("is_read = 0 and receive_uid = ? and notify_type = ?", uid, notifyType).Pluck("target_id", &cids).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetNoReadNotifyCount(userId int64, notifyType string) (count int64) {
	var err error
	db := global.DB.Model(models.McNotify{}).Where("is_read = 0 and receive_uid = ? and notify_type != ?", userId, "follow")
	if notifyType != "" {
		db = db.Where("notify_type = ?", notifyType)
	}
	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateIsReadNotify(userId, notifyId int64, notifyType string) (err error) {
	db := global.DB.Model(models.McNotify{})
	if notifyId > 0 {
		db = db.Where("id = ?", notifyId)
	}
	if notifyType != "" {
		db = db.Where("notify_type = ?", notifyType)
	}
	err = db.Where("receive_uid = ?", userId).Update("is_read", 1).Error
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

func getCommentTextById(id int64) (text string) {
	var err error
	err = global.DB.Model(models.McComment{}).Select("text").Where("id = ?", id).Scan(&text).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getCommentsByReplyUid(replyUid int64) (user []*models.McComment, err error) {
	err = global.DB.Model(models.McComment{}).Where("reply_uid", replyUid).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
