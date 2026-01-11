package withdraw_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func getUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
