package class_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetClassTypes() (typeList []*models.McClassType, err error) {
	err = global.DB.Model(models.McClassType{}).Order("sort asc").Where("status = 1").Find(&typeList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
