package collect_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetCollectList() (collects []*models.McCollect, err error) {
	err = global.DB.Model(models.McCollect{}).Find(&collects).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getClassType(classId int64) (classType int) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("book_type").Where("id = ?", classId).Scan(&classType).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getClassName(classId int64) (className string) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("class_name").Where("id = ?", classId).Scan(&className).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
