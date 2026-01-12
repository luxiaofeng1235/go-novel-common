/*
 * @Descripttion: 用户数据访问（脚手架最小实现）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package user_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
)

// GetUserByDeviceAndOaid deviceid/oaid 不会同时为空；优先按 oaid 查找。
func GetUserByDeviceAndOaid(deviceid, oaid, imei string) (user *models.McUser, err error) {
	if deviceid == "" && oaid == "" {
		return &models.McUser{}, nil
	}

	db := global.DB.Model(models.McUser{})
	if oaid != "" {
		db = db.Where("oaid = ?", oaid)
	}
	if deviceid != "" {
		db = db.Where("deviceid = ?", deviceid)
	}
	if imei != "" {
		db = db.Where("imei = ?", imei)
	}

	err = db.First(&user).Error
	if err != nil {
		// 没查到时直接返回空对象，交由上层决定是否创建游客
		return &models.McUser{}, nil
	}
	return user, nil
}

func GetParentLinkByReffer(referrer string) (parentId int64, parentLink string, err error) {
	if referrer == "" {
		return 0, ",", nil
	}
	var parentUser models.McUser
	err = global.DB.Model(models.McUser{}).Select("id,parent_link").Where("invitation = ?", referrer).First(&parentUser).Error
	if err != nil || parentUser.Id <= 0 {
		return 0, "", fmt.Errorf("%v", "邀请码不存在")
	}
	parentId = parentUser.Id
	parentLink = fmt.Sprintf("%v%v,", parentUser.ParentLink, parentId)
	return parentId, parentLink, nil
}

func UpdateUserByUserId(userId int64, md map[string]interface{}) (err error) {
	if userId <= 0 {
		return nil
	}
	return global.DB.Model(models.McUser{}).Where("id", userId).Updates(md).Error
}
