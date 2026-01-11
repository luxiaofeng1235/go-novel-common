package nsq_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func getClassName(classId int64) (className string) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("class_name").Where("id = ?", classId).Scan(&className).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
