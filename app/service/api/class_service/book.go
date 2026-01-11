package class_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetBookCountByClassId(classId int64) (count int64) {
	err := global.DB.Model(models.McBook{}).Where("status = 1 and cid = ?", classId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookPicByClassId(classId int64) (pic string) {
	err := global.DB.Model(models.McBook{}).Order("id asc").Select("pic").Where("cid = ?", classId).First(&pic).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
