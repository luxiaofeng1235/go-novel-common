package collect_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetCollectById(id int64) (collect *models.McCollect, err error) {
	err = global.DB.Model(models.McCollect{}).Where("id = ?", id).Find(&collect).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
