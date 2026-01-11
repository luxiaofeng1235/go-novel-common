package user_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetUserCountById(id int64) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("is_guest != 1 and id = ?", id).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChildCountById(id int64) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("parent_id = ?", id).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getInviteCion(userId int64) (todayCion float64) {
	err := global.DB.Model(models.McCionChange{}).Select("coalesce(sum(cion), 0)").Where("change_type = 1 and change_type = 4 and uid = ?", userId).Scan(&todayCion).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getTodayCion(userId int64) (todayCion float64) {
	err := global.DB.Model(models.McCionChange{}).Select("coalesce(sum(cion), 0)").Where("change_type = 1 and uid = ? and addtime >= ?", userId, utils.GetTodayUnix()).Scan(&todayCion).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getFollowCount(userId int64) (followCount int64) {
	err := global.DB.Model(models.McUserFollow{}).Where("uid = ?", userId).Count(&followCount).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getFansCount(userId int64) (followCount int64) {
	err := global.DB.Model(models.McUserFollow{}).Where("by_uid = ?", userId).Count(&followCount).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserCountByReffer(referrer string) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("invitation = ?", referrer).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetGuestUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("is_guest = 1 and id = ?", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserCountByTel(tel string) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("tel = ?", tel).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserCountByEmail(email string) (count int64) {
	err := global.DB.Model(models.McUser{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 根据手机号获取相关的信息
func GetUserInfoByMailOrTel(email, tel string) (user *models.McUser, err error) {

	db := global.DB.Model(models.McUser{})
	if email != "" {
		db = db.Where("email = ?", email)
	}
	if tel != "" {
		db = db.Where("tel = ?", tel)
	}
	err = db.Debug().Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserByDeviceid(deviceid string) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("deviceid = ?", deviceid).Debug().Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 通过deviceid和oaid查询
func GetUserByDeviceAndOaid(deviceid, oaid, imei string) (user *models.McUser, err error) {
	//deviceid和oaid不会同时为空
	if deviceid == "" && oaid == "" {
		return &models.McUser{}, err
	}
	db := global.DB.Model(models.McUser{}).Debug()
	//Oaid & imei  同属字段 都为空使用deviceid
	if oaid != "" {
		fmt.Println("oaid不为空。。。。。。。。。。。")
		db = db.Where("oaid = ? ", oaid)
	}
	//如果devideId不为空,查对应的信息（有可能oaid也为空）
	if deviceid != "" {
		db = db.Where("deviceid = ? ", deviceid)
	}
	err = db.Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 根据oaid查询用户
func GetUserByOaid(oaid string) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("oaid = ?", oaid).Debug().Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetParentLinkByReffer(referrer string) (parentId int64, parentLink string, err error) {
	if referrer == "" {
		parentLink = ","
		return
	}
	var parentUser *models.McUser
	err = global.DB.Model(models.McUser{}).Select("id,parent_link").Where("invitation = ?", referrer).Find(&parentUser).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	parentId = parentUser.Id
	parentLink = fmt.Sprintf("%v%v,", parentUser.ParentLink, parentId)
	if parentId <= 0 {
		err = fmt.Errorf("%v", "邀请码不存在")
		return
	}
	return
}

func GetUserByReferrer(referrer string) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("invitation = ?", referrer).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserListByIds(ids []int64) (user []*models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id in ?", ids).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserPasswdIdById(id int64) (pass string) {
	var err error
	err = global.DB.Model(models.McUser{}).Select("pass").Where("id", id).First(&pass).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateUserByUserId(userId int64, md map[string]interface{}) (err error) {
	err = global.DB.Model(models.McUser{}).Debug().Where("id", userId).Updates(md).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdatePassByTel(tel, passwd string) (err error) {
	err = global.DB.Model(models.McUser{}).Where("tel", tel).Update("passwd", passwd).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateReferrerPassById(userId, parentId int64, referrer, parentLink string) (err error) {
	data := make(map[string]interface{})
	data["parent_id"] = parentId
	data["referrer"] = referrer
	data["parent_link"] = parentLink
	err = global.DB.Model(models.McUser{}).Where("id", userId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
