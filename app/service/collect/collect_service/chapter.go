package collect_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func getBookSources() (list []*models.McBookSource, err error) {
	err = global.DB.Model(models.McBookSource{}).Order("uptime asc").Where("is_update = 1").Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getBookSourceById(sourceId int64) (bookSource *models.McBookSource, err error) {
	err = global.DB.Model(models.McBookSource{}).Order("uptime desc").Where("id = ?", sourceId).First(&bookSource).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
