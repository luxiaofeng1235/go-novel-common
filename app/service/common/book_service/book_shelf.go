package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetShelfCountByBookId(bookId, userId int64) (count int64) {
	err := global.DB.Model(models.McBookShelf{}).Where("bid = ? and uid = ?", bookId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
