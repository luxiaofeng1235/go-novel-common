package monitor_service

import (
	"errors"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
	"time"
)

// 获取登录日志列表
func LoginLogListSearch(req *models.LoginLogListReq) ([]*models.SysLoginLog, int64, error) {
	var list []*models.SysLoginLog
	db := global.DB.Model(&models.SysLoginLog{}).Order("id desc")

	loginName := strings.TrimSpace(req.LoginName)
	if loginName != "" {
		db = db.Where("login_name like ?", "%"+loginName+"%")
	}

	ipaddr := strings.TrimSpace(req.Ipaddr)
	if ipaddr != "" {
		db = db.Where("ipaddr like ?", ipaddr)
	}

	status := req.Status
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if req.BeginTime != "" {
		db = db.Where("login_time >=?", utils.DateToUnix(req.BeginTime))
	}
	if req.EndTime != "" {
		db = db.Where("login_time <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 删除登录日志
func DeleteLoginLog(req *models.DeleteLoginLogReq) (res bool, err error) {
	ids := req.LoginLogIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.SysLoginLog{}).Error
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		err = errors.New("删除失败，参数错误")
		return
	}

	return true, nil
}

// 清空登录日志
func ClearLoginLog() (res bool, err error) {
	var loginlog models.SysLoginLog
	err = global.DB.Exec(fmt.Sprintf("truncate table %s", loginlog.TableName())).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// 创建登录日志
func LoginLog(status int, username, ip, userAgent, msg, module string) (InsertId int64, err error) {
	loginlog := new(models.SysLoginLog)
	loginlog.LoginName = username
	loginlog.Ipaddr = ip
	loginlog.LoginLocation = utils.GetCityByIp(loginlog.Ipaddr)
	loginlog.Browser, loginlog.Os = utils.GetBrowser(userAgent)
	loginlog.Status = status
	loginlog.Msg = msg
	loginlog.LoginTime = utils.GetUnix()
	loginlog.Module = module

	if err = global.DB.Create(&loginlog).Error; err != nil {
		return 0, err
	}

	return loginlog.Id, nil
}

// 判断30分钟内有输入次错误密码记录
func GetLoginError(username string) (err error) {
	login_time := time.Now().Unix() - 1800
	var count int64
	global.DB.Model(models.SysLoginLog{}).Where("login_name=?", username).Where("status=?", 0).Where("login_time>=?", login_time).Count(&count)
	//判断用户登陆账户密码错误次数是否超过3次
	if count >= 3 {
		return errors.New("密码尝试次数已超过3次，账户已被暂时锁定，请30分钟后再试！")
	}
	return nil
}
