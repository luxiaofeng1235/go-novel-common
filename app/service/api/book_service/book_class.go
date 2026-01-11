package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetBookClassList() (books []*models.McBookClass, err error) {
	err = global.DB.Model(models.McBookClass{}).Select("id,class_name").Order("sort desc").Limit(6).Find(&books).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetClassIdsByName(className string) (ids []int64) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Where("class_name like ?", "%"+className+"%").Pluck("id", &ids).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTagIdsByName(tagName string) (ids []int64) {
	var err error
	err = global.DB.Model(models.McTag{}).Where("tag_name like ?", "%"+tagName+"%").Pluck("id", &ids).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
