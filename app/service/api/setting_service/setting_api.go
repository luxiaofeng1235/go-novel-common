package setting_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"html"
)

func GetValue(req *models.GetValueReq) (value string, err error) {
	name := req.Name
	if name == "" {
		err = fmt.Errorf("%v", "参数错误")
		return
	}
	value, err = setting_service.GetValueByName(name)
	return
}

func GetValueByNameInfo(name string) (value string, err error) {
	err = global.DB.Model(models.McSetting{}).Select("value").Where("name = ?", name).Scan(&value).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	value = html.UnescapeString(value)
	return
}

func GetRanksName() (ranks []*models.GetRanksRes, err error) {
	var value string
	err = global.DB.Model(models.McSetting{}).Select("value").Where("name = ?", "rank").Scan(&value).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	err = json.Unmarshal([]byte(value), &ranks)
	if err != nil {
		return
	}
	return
}
