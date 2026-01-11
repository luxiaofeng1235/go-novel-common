package class_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

// bookType 1=男生 2=女生
func GetClassByClassType(bookType int, typeId int64) (classList []*models.McBookClass, err error) {
	db := global.DB.Model(models.McBookClass{}).Order("sort asc")
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}
	if typeId > 0 {
		db = db.Where("type_id = ?", typeId)
	}
	db = db.Where("status = 1")
	err = db.Find(&classList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
