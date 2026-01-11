package read_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetReadById(bookId, userId int64) (read *models.McBookRead, err error) {
	err = global.DB.Model(models.McBookRead{}).Where("bid = ? and uid = ?", bookId, userId).Find(&read).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetReadCountByBookId(bookId, userId int64) (count int64) {
	err := global.DB.Model(models.McBookRead{}).Where("bid = ? and uid = ?", bookId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTimeCountByBookId(bookId, userId int64, today string) (count int64) {
	err := global.DB.Model(models.McBookTime{}).Where("bid = ? and uid = ? and day = ?", bookId, userId, today).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTodaySecondsByUserId(userId, todayUnix int64) (seconds int64) {
	err := global.DB.Model(models.McBookTime{}).Select("coalesce(sum(second), 0)").Where("uid = ? and addtime >= ?", userId, todayUnix).Scan(&seconds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 获取阅读任务
func GetReadTaskList() (tasks []*models.McTask, err error) {
	err = global.DB.Model(models.McTask{}).Order("sort desc").Where("status = 1 and welfare_type = 3").Find(&tasks).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateReadCidByUserId(userId, bookId, chapterId int64, chapterName string, textNum int64) (err error) {
	data := make(map[string]interface{})
	data["cid"] = chapterId
	data["chapter_name"] = chapterName
	data["text_num"] = textNum
	data["uptime"] = utils.GetUnix()
	err = global.DB.Model(models.McBookRead{}).Where("uid = ? and bid = ?", userId, bookId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateReadTimeByUserId(today string, userId, bookId, second int64) (err error) {
	data := make(map[string]interface{})
	data["second"] = gorm.Expr("second + ?", second)
	data["uptime"] = utils.GetUnix()
	err = global.DB.Model(models.McBookTime{}).Where("day = ? and uid = ? and bid = ?", today, userId, bookId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteReadByBookIds(bookIds []int64, userId int64) (err error) {
	if len(bookIds) <= 0 || userId <= 0 {
		return
	}
	err = global.DB.Where("bid in ? and uid = ?", bookIds, userId).Delete(&models.McBookRead{}).Error
	return
}

func GetReadList(req *models.BookReadListReq) (reads []*models.BookReadListRes, total int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	var list []*models.McBookRead
	db := global.DB.Model(&models.McBookRead{}).Order("id DESC")

	db = db.Where("uid = ?", userId)

	day := req.Day
	todayUnix := utils.GetTodayUnix()
	yesterdayUnix := utils.GetYesterdayUnix()
	if day == utils.Today {
		db = db.Where("uptime >= ?", todayUnix)
	} else if day == utils.Yesterday {
		db = db.Where("uptime < ? and uptime >= ?", todayUnix, yesterdayUnix)
	} else if day == utils.Agoday {
		db = db.Where("uptime < ?", yesterdayUnix)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}
	if pageSize == 0 {
		pageSize = 1
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		var read *models.BookReadListRes
		read, err = getReadBook(val.Bid, userId)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			continue
		}
		read.Id = val.Id
		read.TextNum = val.TextNum
		read.Addtime = val.Addtime
		reads = append(reads, read)
	}
	return
}

func GetBrowseList(req *models.BrowseListReq) (reads []*models.BookReadListRes, total int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	var list []*models.McBookBrowse
	db := global.DB.Model(&models.McBookBrowse{}).Order("id DESC")

	db = db.Where("uid = ?", userId)

	day := req.Day
	todayUnix := utils.GetTodayUnix()
	yesterdayUnix := utils.GetYesterdayUnix()
	if day == utils.Today {
		db = db.Where("uptime >= ?", todayUnix)
	} else if day == utils.Yesterday {
		db = db.Where("uptime < ? and uptime >= ?", todayUnix, yesterdayUnix)
	} else if day == utils.Agoday {
		db = db.Where("uptime < ?", yesterdayUnix)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}
	if pageSize == 0 {
		pageSize = 1
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		var read *models.BookReadListRes
		read, err = getReadBook(val.Bid, userId)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			continue
		}
		read.Id = val.Id
		read.Addtime = val.Addtime
		reads = append(reads, read)
	}
	return
}
