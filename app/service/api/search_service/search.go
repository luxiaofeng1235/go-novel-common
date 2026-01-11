package search_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetBookSearchList() (searchs []*models.McBookSearch, err error) {
	err = global.DB.Model(models.McBookSearch{}).Order("num desc").Find(&searchs).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetSearchCountByBookName(bookName string) (count int64) {
	var err error
	err = global.DB.Model(models.McBookSearch{}).Distinct("uid").Where("search_name = ?", bookName).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
