package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	Minute = "minute"
	Day    = "day"
	Month  = "month"
)

// 解析时间戳
func ParseDateTzTime(daytime string) (stamp string) {
	if daytime == "" {
		return
	}
	loc, _ := time.LoadLocation("Local")
	var to time.Time
	if strings.Contains(daytime, "+08:00") {
		to, _ = time.ParseInLocation("2006-01-02T15:04:05+08:00", daytime, loc)
	} else if strings.Contains(daytime, "Z") {
		to, _ = time.ParseInLocation("2006-01-02T15:04:05Z", daytime, loc)
	}
	stamp = to.Format("2006-01-02 15:04:05")
	return stamp
}

// 时间戳转时间
func UnixToDatetime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// 获取每周的星期一的日期
func GetThisWeekFirstDate() (andoyRes string) {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	toBeCharge := weekStartDate.Format("2006-01-02")

	timeLayout := "2006-01-02"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	thisMonday := theTime.Unix()
	andoyRes = time.Unix(thisMonday, 0).Format("200601")
	return
}

// 时间戳转日期
func UnixToDate(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02")
}

// 时间戳转日期
func UnixToDayTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("01-02 15:04")
}

// 获取当前时间  20060102
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// 获取当前时间  20060102
func GetDaytime() string {
	template := "20060102150405"
	return time.Now().Format(template)
}

func GetDate() string {
	return time.Now().Format("2006-01-02")
}

// 获取明天的日期
func GetTomorrowDate() string {
	day := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	return day
}

func GetDatetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前时间戳(毫秒)
func GetUnix() int64 {
	return time.Now().Unix()
}

// 获取当前时间戳(纳秒)
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取昨天的时间戳
func GetYesterdayUnix() int64 {
	now := time.Now()
	// 获取昨天的日期
	yesterday := now.AddDate(0, 0, -1)
	// 设置为昨天的零点
	midnight := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
	// 转换为时间戳
	return midnight.Unix()
}

// 获取今天的凌晨的时间戳
func GetTodayUnix() int64 {
	now := time.Now()
	// 获取年月日
	year, month, day := now.Date()
	// 设置为当天的零点
	midnight := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	// 转换为时间戳
	return midnight.Unix()
}

func GetWeekyUnix() (startOfWeekUnix int64) {
	// 获取当前时间
	now := time.Now()
	// 计算当前日期是周几（Sunday = 0, Monday = 1, ..., Saturday = 6）
	dayOfWeek := int(now.Weekday())
	// 计算当前时间距离本周开始的天数（负数表示过去的天数，正数表示未来的天数）
	daysSinceMonday := -dayOfWeek + 1
	// 计算本周周一的日期
	startOfWeek := now.AddDate(0, 0, daysSinceMonday)
	// 将时间设置为当天的零点
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	// 将时间戳转换为 Unix 时间戳
	startOfWeekUnix = startOfWeek.Unix()
	return
}

// 获取明天零点时间戳
func GetTomorrowUnix() time.Time {
	day := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	template := "2006-01-02"
	t, err := time.ParseInLocation(template, day, time.Local)
	if err != nil {
		return t
	}
	return t
}

func AgoTime(second int64) (agoTime int64) {
	duration := time.Duration(second) * time.Second
	agoTime = time.Now().Add(-duration).Unix()
	return
}

func GetAgoDayUnix(agoDay int) (agoDayUnix int64) {
	agoDayUnix = time.Now().AddDate(0, 0, -agoDay).Unix()
	return
}

func GetAgoDayUnixRange(daysAgo int) (int64, int64) {
	// 获取当前时间
	now := time.Now()
	// 计算几天前的日期
	targetDate := now.Add(-time.Duration(daysAgo) * 24 * time.Hour)
	// 获取几天前的 0 点时间戳
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.Local).Unix()
	// 获取几天前的 23:59:59 时间戳
	endOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, time.Local).Unix()
	return startOfDay, endOfDay
}

func DaysSince(timestamp int64) (day int) {
	if timestamp <= 0 {
		return
	}
	now := time.Now()
	tm := time.Unix(timestamp, 0)

	// 截取到天的时间，忽略时分秒
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tm = time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, tm.Location())

	duration := now.Sub(tm)
	day = int(duration.Hours() / 24)

	return
}

// 2020-05-02
func GetDayBeginAndEndTime(day string) (beginTime, endTime int64, err error) {
	start := fmt.Sprintf("%s %s", day, "00:00:00")
	end := fmt.Sprintf("%s %s", day, "23:59:59")
	beginTime = DateToUnix(start)
	endTime = DateToUnix(end)
	return
}

// 时间转时间戳
// 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02"
	if strings.Contains(str, ".000") {
		template = "2006-01-02T15:04:05.000Z"
	} else if strings.Contains(str, "Z") {
		template = "2006-01-02T15:04:05Z"
	} else if strings.Contains(str, ":") {
		template = "2006-01-02 15:04:05"
	}
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 获取相差时间
func GetHourDiffer(startTime, endTime string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}

func GetCurrentMonth() (month int) {
	month = int(time.Now().Month())
	return
}

func GetCurrentYear() (year int) {
	year = int(time.Now().Year())
	return
}

func GetDaysInMonth(year, month int) []int {
	var days []int

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	for day := firstDay; !day.After(lastDay); day = day.AddDate(0, 0, 1) {
		days = append(days, day.Day())
	}

	return days
}

func GetMonthRange(month int) (startTimeUnix, endTimeUnix int64, err error) {
	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("Invalid month")
	}
	// 获取当前年份
	currentYear := time.Now().Year()
	// 计算月份的开始时间
	startTime := time.Date(currentYear, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// 计算月份的结束时间
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Second)
	// 转换为 Unix 时间戳
	startTimeUnix = startTime.Unix()
	endTimeUnix = endTime.Unix()
	return startTimeUnix, endTimeUnix, nil
}

func FormatBeforeUnixTime(dt int64) string {
	const (
		Minute = 60
		Hour   = 3600
		Day    = 86400
		Month  = 86400 * 30
		Year   = 86400 * 30 * 12
	)

	var res string
	now := time.Now().Unix()
	diff := now - dt
	switch {
	case diff < Minute:
		res = "刚刚"
	case diff < Hour && diff >= Minute:
		res = fmt.Sprintf("%d分钟前", diff/Minute)
	case diff < Day && diff >= Hour:
		res = fmt.Sprintf("%d小时前", diff/Hour)
	case diff < Month && diff >= Day:
		res = fmt.Sprintf("%d天前", diff/Day)
	case diff < Year && diff >= Month:
		res = fmt.Sprintf("%d个月前", diff/Month)
	case diff >= Year:
		res = fmt.Sprintf("%d年前", diff/Year)
	}
	return res
}

// 获取星期
func Getweek() string {
	var datStr string
	day := time.Now().Weekday()
	switch day {
	case 0:
		datStr = "星期天"
	case 1:
		datStr = "星期一"
	case 2:
		datStr = "星期二"
	case 3:
		datStr = "星期三"
	case 4:
		datStr = "星期四"
	case 5:
		datStr = "星期五"
	case 6:
		datStr = "星期六"
	}
	return datStr
}

func GetweekDay() (dayWeek int) {
	// 获取当前时间
	now := time.Now()
	// 计算今天是一周中的第几天
	dayWeek = int(now.Weekday())
	// 转换为从周一开始计算的格式
	if dayWeek == 0 {
		dayWeek = 7
	}
	return
}

// 判断是否可以重置在线时长和在线次数
func IsResetOnline(sType string, timestamp int64) (isBool bool) {
	t := time.Unix(timestamp, 0)
	if sType == Minute {
		diffSeconds := int(time.Since(t).Seconds())
		if diffSeconds > 300 {
			isBool = true
			return
		}
	} else if sType == Day {
		tday := t.Day()
		today := time.Now().Day()
		if tday != today {
			isBool = true
			return
		}
	} else if sType == Month {
		tmonth := t.Month().String()
		cmonth := time.Now().Month().String()
		if tmonth != cmonth {
			isBool = true
			return
		}
	}
	return isBool
}

func IsInThisWeek(timestamp int64) bool {
	now := time.Now()
	//获取本周的起始时间
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1).Truncate(24 * time.Hour)
	//获取本周的结束时间
	weekEnd := weekStart.AddDate(0, 0, 7)
	//转换时间戳为时间
	timestampTime := time.Unix(timestamp, 0)
	//判断时间戳是否在本周内
	return timestampTime.After(weekStart) && timestampTime.Before(weekEnd)
}

func GetWeekDayRange(dayWeek int) (startTime int64, endTime int64) {
	now := time.Now()
	weekday := int(now.Weekday())

	// 计算本周开始的日期
	startWeek := now.AddDate(0, 0, -weekday+1)
	startWeek = time.Date(startWeek.Year(), startWeek.Month(), startWeek.Day(), 0, 0, 0, 0, time.Local)

	// 计算目标天的日期
	targetDay := startWeek.AddDate(0, 0, dayWeek-1)
	startDay := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 0, 0, 0, 0, time.Local)
	endDay := startDay.Add(24*time.Hour - time.Second)
	startTime = startDay.Unix()
	endTime = endDay.Unix()
	return
}
