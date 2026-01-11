package comment_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetCommentPidByCid(commentId int64) (parentId int64) {
	var err error
	err = global.DB.Model(&models.McComment{}).Select("pid").Where("id = ?", commentId).Last(&parentId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
