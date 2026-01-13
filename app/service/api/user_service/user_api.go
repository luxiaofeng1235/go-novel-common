/*
 * @Descripttion: 用户登录/注册/游客登录（脚手架最小实现）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package user_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/config"
	"go-novel/global"
	"go-novel/utils"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/goroom/rand"
)

func hashLoginPasswd(plain string) string {
	salt := strings.TrimSpace(config.GetString("auth.passwordSalt"))
	if salt == "" {
		return utils.Md5(plain)
	}
	return utils.GetMd5(plain, salt)
}

// GuestLogin 游客登录（deviceid / oaid 至少一个不为空）
func GuestLogin(c *gin.Context, req *models.GuestLoginReq) (userInfo *models.McUser, token string, expireTime int64, err error) {
	deviceid := strings.TrimSpace(req.Deviceid)
	referrer := strings.TrimSpace(req.Referrer)
	sex := req.Sex

	mark := utils.GetRequestHeaderByName(c, "Mark")
	devicePackage := utils.GetRequestHeaderByName(c, "Package")
	imei := utils.GetRequestHeaderByName(c, "Imei")
	oaid := utils.GetRequestHeaderByName(c, "Oaid")
	ip := utils.RemoteIp(c)

	if deviceid == "" && oaid == "" {
		return nil, "", 0, fmt.Errorf("%v", "设备ID或者oaid不能同时为空")
	}

	userInfo, err = GetUserByDeviceAndOaid(deviceid, oaid, imei)
	if err != nil {
		return nil, "", 0, err
	}

	if userInfo == nil || userInfo.Id <= 0 {
		parentId, parentLink, err := GetParentLinkByReffer(referrer)
		if err != nil {
			return nil, "", 0, err
		}

		username := utils.GetGuestName()
		guestPass := utils.Md5(utils.RandomString("guest", 24))

		user := models.McUser{
			ParentId:   parentId,
			ParentLink: parentLink,
			Username:   username,
			Passwd:     guestPass,
			Nickname:   rand.GetRand().ChineseName(),
			Invitation: utils.RandomString("code", 6),
			Status:     1,
			IsGuest:    1,
			Sex:        int(sex),
			Deviceid:   deviceid,
			Mark:       mark,
			Oaid:       oaid,
			Package:    devicePackage,
			Imei:       imei,
			Ip:         ip,
			Addtime:    utils.GetUnix(),
		}

		if err = global.DB.Model(models.McUser{}).Create(&user).Error; err != nil {
			return nil, "", 0, err
		}
		userInfo = &user
	} else {
		mu := map[string]interface{}{
			"mark":            mark,
			"imei":            imei,
			"package":         devicePackage,
			"ip":              ip,
			"last_login_time": utils.GetUnix(),
			"uptime":          utils.GetUnix(),
		}
		if deviceid != "" {
			mu["deviceid"] = deviceid
		}
		if oaid != "" {
			mu["oaid"] = oaid
		}
		_ = UpdateUserByUserId(userInfo.Id, mu)
	}

	token, expireTime, err = utils.GenerateToken(userInfo.Id, userInfo.Username, 1)
	if err != nil {
		return nil, "", 0, err
	}
	userInfo.Passwd = ""
	return userInfo, token, expireTime, nil
}

// Register 账号密码注册
func Register(c *gin.Context, req *models.RegisterReq) (token string, expireTime int64, err error) {
	username := strings.TrimSpace(req.Username)
	passwd := strings.TrimSpace(req.Passwd)
	nickname := strings.TrimSpace(req.Nickname)
	referrer := strings.TrimSpace(req.Referrer)
	deviceid := strings.TrimSpace(req.Deviceid)

	if username == "" {
		return "", 0, fmt.Errorf("%v", "账号不能为空")
	}
	if passwd == "" {
		return "", 0, fmt.Errorf("%v", "密码不能为空")
	}

	var count int64
	if err = global.DB.Model(models.McUser{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return "", 0, err
	}
	if count > 0 {
		return "", 0, fmt.Errorf("%v", "账号已存在")
	}

	parentId, parentLink, err := GetParentLinkByReffer(referrer)
	if err != nil {
		return "", 0, err
	}

	mark := utils.GetRequestHeaderByName(c, "Mark")
	devicePackage := utils.GetRequestHeaderByName(c, "Package")
	imei := utils.GetRequestHeaderByName(c, "Imei")
	oaid := utils.GetRequestHeaderByName(c, "Oaid")
	ip := utils.RemoteIp(c)

	if nickname == "" {
		nickname = rand.GetRand().ChineseName()
	}
	//类型注入和依赖管理
	user := models.McUser{
		ParentId:   parentId,
		ParentLink: parentLink,
		Username:   username,
		Passwd:     hashLoginPasswd(passwd),
		Nickname:   nickname,
		Invitation: utils.RandomString("code", 6),
		Status:     1,
		IsGuest:    0,
		Deviceid:   deviceid,
		Mark:       mark,
		Oaid:       oaid,
		Package:    devicePackage,
		Ip:         ip,
		Imei:       imei,
		Addtime:    utils.GetUnix(),
	}

	if err = global.DB.Model(models.McUser{}).Create(&user).Error; err != nil {
		return "", 0, err
	}

	token, expireTime, err = utils.GenerateToken(user.Id, user.Username, 1)
	if err != nil {
		return "", 0, err
	}
	return token, expireTime, nil
}

// Login 账号密码登录
func Login(c *gin.Context, req *models.LoginReq) (token string, expireTime int64, err error) {
	username := strings.TrimSpace(req.Username)
	passwd := strings.TrimSpace(req.Passwd)
	deviceid := strings.TrimSpace(req.Deviceid)

	if username == "" {
		return "", 0, fmt.Errorf("%v", "账号不能为空")
	}
	if passwd == "" {
		return "", 0, fmt.Errorf("%v", "密码不能为空")
	}

	var user models.McUser
	if err = global.DB.Model(models.McUser{}).Where("username = ?", username).First(&user).Error; err != nil {
		return "", 0, fmt.Errorf("%v", "账号不存在，请先注册")
	}
	if user.Status == 0 {
		return "", 0, fmt.Errorf("%v", "账户已被锁定~")
	}
	if user.Passwd == "" || hashLoginPasswd(passwd) != user.Passwd {
		return "", 0, fmt.Errorf("%v", "密码不正确~")
	}

	mark := utils.GetRequestHeaderByName(c, "Mark")
	devicePackage := utils.GetRequestHeaderByName(c, "Package")
	imei := utils.GetRequestHeaderByName(c, "Imei")
	oaid := utils.GetRequestHeaderByName(c, "Oaid")
	ip := utils.RemoteIp(c)
	mu := map[string]interface{}{
		"mark":            mark,
		"imei":            imei,
		"package":         devicePackage,
		"ip":              ip,
		"last_login_time": utils.GetUnix(),
		"uptime":          utils.GetUnix(),
	}
	if deviceid != "" {
		mu["deviceid"] = deviceid
	}
	if oaid != "" {
		mu["oaid"] = oaid
	}
	_ = UpdateUserByUserId(user.Id, mu)

	token, expireTime, err = utils.GenerateToken(user.Id, user.Username, 1)
	if err != nil {
		return "", 0, err
	}
	return token, expireTime, nil
}

// Logoff 账号注销：将用户状态置为 0（注销/锁定）
func Logoff(c *gin.Context, userID int64) error {
	if userID <= 0 {
		return fmt.Errorf("用户不存在")
	}

	var user models.McUser
	if err := global.DB.Model(models.McUser{}).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			global.Sqllog.Errorf("用户不能存在 ：%v", err.Error())
			return fmt.Errorf("用户不存在")
		}
		return err
	}
	if user.Status == 0 {
		return nil
	}

	update := map[string]interface{}{
		"status": 0,
		"uptime": utils.GetUnix(),
	}
	return global.DB.Model(models.McUser{}).Where("id = ?", userID).Updates(update).Error
}

// GetUserInfoByUserID 查询用户信息（返回前清空 passwd）
func GetUserInfoByUserID(userID int64) (user *models.McUser, error error) {
	if userID <= 0 {
		return nil, fmt.Errorf("用户不存在")
	}
	err := global.DB.Model(models.McUser{}).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}
	if user.Status == 0 {
		return nil, fmt.Errorf("账户已被锁定~")
	}
	user.Passwd = ""
	return user, nil
}
