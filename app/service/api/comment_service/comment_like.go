package comment_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetIsPraiseByUserId(commentId, userId int64) (count int64) {
	err := global.DB.Model(models.McCommentPraise{}).Where("cid = ? and uid = ?", commentId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
