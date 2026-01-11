package user_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetUserIdByUsername(username string) (userId int64) {
	var err error
	err = global.DB.Model(models.McUser{}).Select("id").Where("username", username).First(&userId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookTypeByUserId(userId int64) (bookType int) {
	var err error
	err = global.DB.Model(models.McUser{}).Select("sex").Where("id = ?", userId).First(&bookType).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserByTel(tel string) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("tel = ?", tel).Last(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserByEmail(email string) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("email = ?", email).Last(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 根据用户id获取用户信息
func GetUserById(id int64) (*models.McUser, error) {
	var sysuser *models.McUser
	err := global.DB.Model(models.McUser{}).Where("id", id).First(&sysuser).Error
	return sysuser, err
}

func CheckTelUnique(tel string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McUser{}).Where("tel = ?", tel)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func CheckEmailUnique(email string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McUser{}).Where("email = ?", email)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func UserLogoff(userId int64) (err error) {
	data := make(map[string]interface{})
	data["tel"] = ""
	data["email"] = ""
	data["is_guest"] = ""
	data["regist_id"] = ""
	data["is_checkin_remind"] = ""
	data["vip"] = 0
	data["rmb"] = 0
	data["cion"] = 0
	err = global.DB.Model(models.McUser{}).Where("id", userId).Updates(&data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
