/*
 * @Descripttion: 用户信息查询（基于 token 解析出的 username）
 * @Author: red
 * @Date: 2026-01-13 09:40:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 09:40:00
 */
package user_service

import (
	"fmt"
	"strings"

	"go-novel/app/models"
	"go-novel/global"

	"gorm.io/gorm"
)

func GetUserInfoByUsername(username string) (*models.McUser, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("用户不存在")
	}

	var user models.McUser
	err := global.DB.Model(models.McUser{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	user.Passwd = ""
	return &user, nil
}

func GetUserInfoByUserID(userID int64) (*models.McUser, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	var user models.McUser
	err := global.DB.Model(models.McUser{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	user.Passwd = ""
	return &user, nil
}
