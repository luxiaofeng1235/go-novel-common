package setting_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"html"
	"math"
)

func GetValueByName(name string) (value string, err error) {
	err = global.DB.Model(models.McSetting{}).Select("value").Where("name = ?", name).Scan(&value).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	value = html.UnescapeString(value)
	return
}

func GetCionRate() (rate int64, err error) {
	cionRate, err := GetValueByName("cionRate")
	if err != nil {
		global.Errlog.Errorf("获取金币对人民币汇率错误 err=%v", err.Error())
		err = fmt.Errorf("金币转换汇率配置错误 %v", err.Error())
		return
	}
	rate = utils.FormatInt64(cionRate)
	return
}

func GetMoneyByCion(gold int64) (rate int64, rmb float64, err error) {
	rate, err = GetCionRate()
	if err != nil {
		return
	}
	rmb = float64(gold) / float64(rate)
	rmb = math.Round(rmb*100) / 100
	return
}

func GetCionByMoney(rmb float64) (rate, gold int64, err error) {
	rate, err = GetCionRate()
	if err != nil {
		return
	}
	gold = int64(math.Floor(rmb * float64(rate)))
	return
}

func GetInviteGive() (inviteGiveCionInt64 int64) {
	inviteGiveCion, err := GetValueByName("inviteGiveCion")
	if err != nil {
		global.Errlog.Errorf("获取邀请赠送金币错误 err=%v", err.Error())
		return
	}
	inviteGiveCionInt64 = utils.FormatInt64(inviteGiveCion)
	return
}
