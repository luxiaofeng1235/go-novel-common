package read_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func DeleteBrowseByBookIds(bookIds []int64, userId int64) (err error) {
	if len(bookIds) <= 0 || userId <= 0 {
		return
	}
	err = global.DB.Where("bid in ? and uid = ?", bookIds, userId).Delete(&models.McBookBrowse{}).Error
	return
}
