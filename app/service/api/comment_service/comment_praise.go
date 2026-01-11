package comment_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"gorm.io/gorm"
)

func GetPraiseCountByUid(userId, commentId int64) (count int64) {
	err := global.DB.Model(models.McCommentPraise{}).Where("uid = ? and cid = ?", userId, commentId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdatePraiseByCommentId(commentId int64, isPraise int) (err error) {
	ma := make(map[string]interface{})
	if isPraise > 0 {
		ma["praise_count"] = gorm.Expr("praise_count + 1")
	} else {
		ma["praise_count"] = gorm.Expr("praise_count - 1")
	}
	err = global.DB.Model(models.McComment{}).Where("id = ?", commentId).Updates(&ma).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeletePraiseUid(userId, commentId int64) (err error) {
	err = global.DB.Where("uid = ? and cid = ?", userId, commentId).Delete(&models.McCommentPraise{}).Error
	if err != nil {
		return
	}
	return
}
