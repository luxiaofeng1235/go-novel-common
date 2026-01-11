package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func getClassTypeByName(className string) (classId int64, classType int) {
	var err error
	var info models.McBookClass
	err = global.DB.Model(models.McBookClass{}).Select("id,book_type").Where("class_name = ?", className).Last(&info).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	classId = info.Id
	classType = info.BookType
	return
}
