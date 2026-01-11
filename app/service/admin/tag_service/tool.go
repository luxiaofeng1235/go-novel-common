package tag_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetBookCountByTagId(columnType int, tagId int64) (count int64) {
	db := global.DB.Model(models.McBook{}).Where("tid = ?", tagId)
	if columnType == 1 {
		db = db.Where("is_rec = 1")
	} else if columnType == 2 {
		db = db.Where("is_hot = 1")
	} else if columnType == 3 {
		db = db.Where("is_classic = 1")
	}
	err := db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
