package nsq_service

import (
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
)

func getSleepSecond() (sleepSecond int64) {
	var collectSleep string
	var err error
	collectSleep, err = setting_service.GetValueByName("collectSleep")
	if err != nil {
		global.Collectlog.Errorf("获取采集间隔时间失败 %v", err.Error())
		return
	}
	sleepSecond = utils.FormatInt64(collectSleep)
	return
}

func delBookSourceById(bookSourceId int64) (err error) {
	err = global.DB.Where("bid = ?", bookSourceId).Delete(&models.McBookSource{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func updateBookSource(bookSourceId int64, data map[string]interface{}) (err error) {
	err = global.DB.Model(models.McBookSource{}).Where("id = ?", bookSourceId).Updates(&data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
