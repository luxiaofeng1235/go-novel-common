package checkin_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
	"time"
)

func GetCheckinList(req *models.CheckinListReq) (checkinListRes *models.CheckinListRes, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	checkinListRes = new(models.CheckinListRes)
	today := utils.GetweekDay()
	checkinListRes.Today = today

	checkinList, err := GetCheckinCionList()
	if err != nil {
		return
	}
	if len(checkinList) <= 0 {
		return
	}
	var checkdays []*models.CheckinDay
	for _, val := range checkinList {
		var isCheck int
		day := val.Day
		if today >= day {
			startTime, endTime := utils.GetWeekDayRange(day)
			count := GetCountByDayTime(userId, startTime, endTime)
			if count > 0 {
				isCheck = 1
			}
		}

		checkday := &models.CheckinDay{
			Day:     val.Day,
			Cion:    val.Cion,
			Vip:     val.Vip,
			IsCheck: isCheck,
		}
		checkdays = append(checkdays, checkday)
	}
	checkinListRes.List = checkdays
	return
}

func GetCheckin(req *models.CheckinReq) (cion int64, vip int64, err error) {
	userId := req.UserId
	isReissue := req.IsReissue
	day := req.Day
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	today := utils.GetweekDay()
	var info *models.McCheckinReward

	if isReissue > 0 {
		if day > today {
			err = fmt.Errorf("%v", "签到还未开始 无法补签")
			return
		}

		startTime, endTime := utils.GetWeekDayRange(day)
		count := GetCountByDayTime(userId, startTime, endTime)
		if count > 0 {
			err = fmt.Errorf("%v", "补签失败 已经签到过啦~")
			return
		}

		info, err = GetCheckinCionByDay(day)
		if err != nil {
			return
		}
		if info.Id <= 0 {
			return
		}
		cion = info.Cion
		vip = info.Vip
		checkin := models.McCheckin{
			Uid:       userId,
			Day:       day,
			Cion:      cion,
			Vip:       vip,
			IsReissue: 1,
			Retime:    utils.GetUnix(),
			Addtime:   utils.GetAgoDayUnix(today - day),
		}
		if err = global.DB.Create(&checkin).Error; err != nil {
			global.Sqllog.Errorf("补签失败 err=%v", err.Error())
			return
		}
		err = checkinReward(cion, vip, userId)
		return
	}

	startTime, endTime := utils.GetWeekDayRange(today)
	count := GetCountByDayTime(userId, startTime, endTime)
	if count > 0 {
		err = fmt.Errorf("%v", "今日您已经签到过啦")
		return
	}

	info, err = GetCheckinCionByDay(today)
	if err != nil {
		return
	}
	if info.Id <= 0 {
		return
	}
	cion = info.Cion
	vip = info.Vip
	checkin := models.McCheckin{
		Uid:       userId,
		Day:       today,
		Cion:      cion,
		Vip:       vip,
		IsReissue: 0,
		Retime:    0,
		Addtime:   utils.GetUnix(),
	}
	if err = global.DB.Create(&checkin).Error; err != nil {
		global.Sqllog.Errorf("签到失败 err=%v", err.Error())
		return
	}

	err = checkinReward(cion, vip, userId)
	return
}

func checkinReward(cion, vip, userId int64) (err error) {
	tx := global.DB.Begin()

	if cion > 0 {
		change := models.McCionChange{
			Tid:        0,
			Uid:        userId,
			Cion:       cion,
			ChangeType: 1,
			OperatType: 1,
			Addtime:    utils.GetUnix(),
		}
		err = tx.Model(models.McCionChange{}).Create(&change).Error
		if err != nil {
			tx.Rollback()
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	}

	err = tx.Model(models.McUser{}).Where("id = ?", userId).Update("cion", gorm.Expr("cion +?", cion)).Error

	us := make(map[string]interface{})
	if cion > 0 {
		us["cion"] = gorm.Expr("cion + ?", cion)
	}
	if vip > 0 {
		us["vip"] = 1
		viptime := getUserVipTimeById(userId)
		timeStamp := utils.GetUnix()
		if viptime > timeStamp {
			us["viptime"] = gorm.Expr("viptime + ?", 86400*vip)
		} else {
			us["viptime"] = utils.GetUnix() + 86400*vip
		}
	}
	err = tx.Model(models.McUser{}).Where("id = ?", userId).Updates(us).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func CheckinHistory(req *models.CheckinHistoryReq) (historys []*models.CheckinHistoryRes, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	year := req.Year
	month := req.Month
	if year <= 0 {
		year = utils.GetCurrentYear()
	}
	if month <= 0 || month > 12 {
		month = utils.GetCurrentMonth()
	}
	days := utils.GetDaysInMonth(year, month)
	var checkins []*models.CheckinHistoryRes
	checkins, err = getChekinHistory(userId, year, month)

	for _, day := range days {
		var cion, vip, isCheckin int64
		for _, checkin := range checkins {
			if day == checkin.Day {
				cion = checkin.Cion
				vip = checkin.Vip
				isCheckin = 1
			}
		}
		history := &models.CheckinHistoryRes{
			Day:       day,
			Cion:      cion,
			Vip:       vip,
			IsCheckin: isCheckin,
		}
		historys = append(historys, history)
	}
	return
}

func CheckRemind() {
	ticker := time.NewTicker(time.Minute * 10) // 每 10 检查一次
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := CheckRemindPush()
			if err != nil {
				global.Errlog.Errorf("%v", err.Error())
			}
		}
	}
}
func CheckRemindPush() (err error) {
	now := time.Now()
	year, month, day := now.Date()
	targetTime := time.Date(year, month, day, 21, 0, 0, 0, now.Location())

	if !now.After(targetTime) {
		return
	}

	todayUnix := utils.GetTodayUnix()
	var registIds []string
	db := global.DB.Model(models.McUser{}).Where("is_checkin_remind = 1 and regist_id != '' and last_remind_time <= ?", todayUnix).Pluck("regist_id", &registIds)
	err = db.Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if len(registIds) <= 0 {
		return
	}
	_, err = utils.JpushMsg("您今天还未进行签到", registIds)
	if err != nil {
		return
	}
	unix := utils.GetUnix()
	err = db.Update("last_remind_time", unix).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	time.Sleep(time.Minute * 10)
	return
}

func OpenCheckinRemind(req *models.OpenCheckinRemindReq) (err error) {
	isCheckinRemind := req.IsCheckinRemind
	registId := req.RegistId
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	data := make(map[string]interface{})
	data["is_checkin_remind"] = isCheckinRemind
	if registId != "" {
		data["regist_id"] = registId
	}
	err = global.DB.Model(models.McUser{}).Where("id = ?", userId).Updates(data).Error
	if err != nil {
		return
	}
	return
}
