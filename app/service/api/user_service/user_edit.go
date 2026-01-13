/*
 * @Descripttion: 用户信息编辑（脚手架最小实现）
 * @Author: red
 * @Date: 2026-01-13 10:50:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-13 10:50:00
 */
package user_service

import (
	"errors"
	"fmt"
	"strings"

	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EditUser 支持按 type 编辑用户字段：
// - tel/email/nickname/sex/pic/pass/invite/book_type
// 说明：当前脚手架不接短信/邮箱验证码，仅校验 code 非空（若后续接入可在此处做验证码校验）。
func EditUser(c *gin.Context, req *models.EditUserReq) (*models.McUser, error) {
	if req == nil {
		return nil, fmt.Errorf("参数错误")
	}
	if req.UserId <= 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	editType := strings.TrimSpace(req.Type)
	if editType == "" {
		return nil, fmt.Errorf("type不能为空")
	}

	var user models.McUser
	if err := global.DB.Model(models.McUser{}).Where("id = ?", req.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	if user.Status == 0 {
		return nil, fmt.Errorf("账户已被锁定~")
	}

	update := map[string]interface{}{
		"uptime": utils.GetUnix(),
	}

	switch editType {
	case "tel":
		tel := strings.TrimSpace(req.Tel)
		if tel == "" {
			return nil, fmt.Errorf("tel不能为空")
		}
		if strings.TrimSpace(req.Code) == "" {
			return nil, fmt.Errorf("code不能为空")
		}
		var cnt int64
		if err := global.DB.Model(models.McUser{}).Where("tel = ? AND id <> ?", tel, req.UserId).Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, fmt.Errorf("手机号已被使用")
		}
		update["tel"] = tel
	case "email":
		email := strings.TrimSpace(req.Email)
		if email == "" {
			return nil, fmt.Errorf("email不能为空")
		}
		if strings.TrimSpace(req.Code) == "" {
			return nil, fmt.Errorf("code不能为空")
		}
		var cnt int64
		if err := global.DB.Model(models.McUser{}).Where("email = ? AND id <> ?", email, req.UserId).Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, fmt.Errorf("邮箱已被使用")
		}
		update["email"] = email
	case "nickname":
		nickname := strings.TrimSpace(req.Nickname)
		if nickname == "" {
			return nil, fmt.Errorf("nickname不能为空")
		}
		update["nickname"] = nickname
	case "sex":
		if req.Sex != 1 && req.Sex != 2 {
			return nil, fmt.Errorf("sex不合法")
		}
		update["sex"] = req.Sex
	case "pic":
		pic := strings.TrimSpace(req.Pic)
		if pic == "" {
			return nil, fmt.Errorf("pic不能为空")
		}
		update["pic"] = pic
	case "book_type":
		update["book_type"] = req.BookType
	case "pass":
		oldPass := strings.TrimSpace(req.OldPasswd)
		newPass := strings.TrimSpace(req.Passwd)
		if oldPass == "" {
			return nil, fmt.Errorf("old_passwd不能为空")
		}
		if newPass == "" {
			return nil, fmt.Errorf("passwd不能为空")
		}
		if user.Passwd == "" || hashLoginPasswd(oldPass) != user.Passwd {
			return nil, fmt.Errorf("旧密码不正确~")
		}
		update["passwd"] = hashLoginPasswd(newPass)
	case "invite":
		// 约定：invite 使用 code 作为邀请码
		referrer := strings.TrimSpace(req.Code)
		if referrer == "" {
			return nil, fmt.Errorf("code不能为空")
		}
		if user.ParentId > 0 {
			return nil, fmt.Errorf("已绑定邀请码")
		}
		parentId, parentLink, err := GetParentLinkByReffer(referrer)
		if err != nil {
			return nil, err
		}
		update["parent_id"] = parentId
		update["parent_link"] = parentLink
		update["referrer"] = referrer
	default:
		return nil, fmt.Errorf("type不支持")
	}

	if err := global.DB.Model(models.McUser{}).Where("id = ?", req.UserId).Updates(update).Error; err != nil {
		return nil, err
	}

	// 返回更新后的用户信息（清空密码）
	return GetUserInfoByUserID(req.UserId)
}
