/*
 * @Descripttion: 用户验证码校验工具
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package user_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

// 完善用户信息
func GetSaveCountByUserId(uid, tid int64) (count int64) {
	err := global.DB.Model(models.McTaskList{}).Where("uid = ? and tid = ?", uid, tid).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func checkEmailCode(email, code string) (err error) {
	//判断是否开启超级验证码
	if utils.SuperCodeOpenStatus {
		SuperCode := utils.SuperCode
		if code != SuperCode {
			//效验邮箱验证码
			err = utils.IsEmailCode(email, code)
			if err != nil {
				err = fmt.Errorf("%v", "验证码不正确或已过期")
				return err
			}
		}
	} else {
		//效验邮箱验证码
		err = utils.IsEmailCode(email, code)
		if err != nil {
			err = fmt.Errorf("%v", "验证码不正确或已过期")
			return err
		}
	}
	return
}

func checkTelCode(tel, code string) (err error) {
	//判断是否开启超级验证码
	if utils.SuperCodeOpenStatus {
		SuperCode := utils.SuperCode
		if code != SuperCode {
			//效验短信验证码
			err = utils.IsYzm(tel, code)
			if err != nil {
				err = fmt.Errorf("%v", "验证码不正确或已过期")
				return err
			}
		}
	} else {
		//效验短信验证码
		err = utils.IsYzm(tel, code)
		if err != nil {
			err = fmt.Errorf("%v", "验证码不正确或已过期")
			return err
		}
	}
	return
}
