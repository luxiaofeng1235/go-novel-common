package user_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/user_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetUserById(id int64) (*models.McUser, error) {
	var sysuser *models.McUser
	err := global.DB.Model(models.McUser{}).Where("id", id).First(&sysuser).Error
	return sysuser, err
}

// 获取用户列表
func UserListSearch(req *models.UserListReq) (list []*models.McUser, total int64, err error) {
	db := global.DB.Model(&models.McUser{}).Order("id desc")

	id := req.Id
	if id > 0 {
		db = db.Where("id =  ?", id)
	}

	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname = ?", nickname)
	}

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username = ?", username)
	}

	referrer := strings.TrimSpace(req.Referrer)
	if referrer != "" {
		db = db.Where("referrer = ?", referrer)
	}

	tel := strings.TrimSpace(req.Tel)
	if tel != "" {
		db = db.Where("tel = ?", tel)
	}

	email := strings.TrimSpace(req.Email)
	if email != "" {
		db = db.Where("email = ?", email)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}

	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	pageNum := req.PageNum
	pageSize := req.PageSize

	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}
	for _, val := range list {
		val.Pic = utils.GetAdminFileUrl(val.Pic)
	}
	return list, total, err
}

func UpdateUser(req *models.UpdateUserReq) (res bool, err error) {
	nickname := req.Nickname
	tel := req.Tel
	email := req.Email
	rmb := req.Rmb
	cion := req.Cion
	status := req.Status
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "用户ID不能为空")
		return
	}
	if !user_service.CheckTelUnique(tel, userId) {
		err = fmt.Errorf("%v", "手机号已经存在")
		return
	}
	if !user_service.CheckEmailUnique(email, userId) {
		err = fmt.Errorf("%v", "邮箱已经存在")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["nickname"] = nickname
	mapData["tel"] = tel
	mapData["email"] = email
	mapData["rmb"] = rmb
	mapData["cion"] = cion
	mapData["status"] = status
	mapData["uptime"] = utils.GetUnix()

	if err = global.DB.Model(models.McUser{}).Where("id", userId).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DelUser(req *models.DelUserReq) (res bool, err error) {
	ids := req.UserIds
	if len(ids) > 0 {
		var mapData = make(map[string]interface{})
		mapData["status"] = 0
		err = global.DB.Model(models.McUser{}).Where("id in(?)", ids).Updates(&mapData).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "修改失败，参数错误")
		return
	}
	return true, nil
}
