package book_service

import (
	"github.com/tidwall/gjson"
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"log"
	"strconv"
)

func GetReadTextNumByBookId(bookId, userId int64) (textNum int64) {
	var err error
	err = global.DB.Model(models.McBookRead{}).Select("text_num").Where("bid = ? and uid = ?", bookId, userId).Scan(&textNum).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetReadChapterIdByBookId(bookId, userId int64) (cid int64, chapterName string) {
	var err error
	var read *models.McBookRead
	err = global.DB.Model(models.McBookRead{}).Where("bid = ? and uid = ?", bookId, userId).First(&read).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	cid = read.Cid
	chapterName = read.ChapterName
	return
}

// 获取当前用户的所有的阅读时间统计
func GetReadChapterAllTime(userId int64) (seconds int64) {
	var err error
	err = global.DB.Model(models.McBookTime{}).Debug().Select("coalesce(sum(second), 0)").Where("uid = ? ", userId).Scan(&seconds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 获取当前用户的所有的阅读是否为新手
func GetReadToUserBeyond(userId int64) (residueTime int64) {
	itemData, err := setting_service.GetValueByName("advertNewbieSet")
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//解析数组对象
	duration_val := gjson.Get(itemData, "setting_time").Value()
	value := ""
	if duration_val != nil {
		value = duration_val.(string)
	}
	//再换成int对象来处理
	systemValue, _ := strconv.Atoi(value) //获取配置的并转换为整型方便进行计算
	residueTime = 0                       //默认新手保护期未过期为0
	if systemValue > 0 {
		var unix_time int64
		switch systemValue {
		case 0:
			unix_time = 0 //0表示无新手保护期
		case 1:
			unix_time = 5 //5分钟
		case 2:
			unix_time = 10 //10分钟
		case 3:
			unix_time = 30 //30分钟
		case 4:
			unix_time = 45 //45分钟
		case 5:
			unix_time = 60 //60分钟--一小时
		}
		//获取当前的已经阅读的总数信息
		totalNum := GetReadChapterAllTime(userId)
		log.Println("用户阅读的总时长", totalNum)
		second := totalNum / 60 //换算成分钟数，看用户的在线分钟时间
		res := unix_time - second
		//利用系统的时间减掉 用户已阅读的分钟做计算。如果剩余时间小于0 说明已经不是新手保护期了，否则还在保护器内
		if res > 0 {
			//如果大于0 说明还是新手保护期，把剩余时间返回
			residueTime = res
		}
	}
	return
}
