package read_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/task_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetReadRes(req *models.BookReadListReq) (readRes *models.BookReadRes, err error) {
	pageNum := req.Page
	pageSize := req.Size
	req.Page = 0
	req.Size = 0
	var today []*models.BookReadListRes
	var yesterday []*models.BookReadListRes
	var agoday []*models.BookReadListRes
	req.Day = utils.Today
	today, _, err = GetReadList(req)
	if err != nil {
		return
	}
	req.Day = utils.Yesterday
	yesterday, _, err = GetReadList(req)
	if err != nil {
		return
	}

	var total int64
	req.Day = utils.Agoday
	req.Page = pageNum
	req.Size = pageSize
	agoday, total, err = GetReadList(req)
	if err != nil {
		return
	}
	readRes = &models.BookReadRes{
		Today:     today,
		Yesterday: yesterday,
		Agoday:    agoday,
		Total:     total,
	}
	return
}

func ReadAdd(req *models.ReadAddReq) (err error) {
	bookId := req.BookId
	chapterId := req.ChapterId
	chapterName := req.ChapterName
	textNum := req.TextNum
	second := req.Second
	userId := req.UserId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	//if chapterId <= 0 {
	//	err = fmt.Errorf("%v", "章节ID不能为空")
	//	return
	//}
	if textNum <= 0 {
		err = fmt.Errorf("%v", "阅读字数不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	var count int64
	count = GetReadCountByBookId(bookId, userId)
	if count <= 0 {
		read := models.McBookRead{
			Uid:         userId,
			Bid:         bookId,
			Cid:         chapterId,
			ChapterName: chapterName,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
			Uptime:      utils.GetUnix(),
		}
		if err = global.DB.Create(&read).Error; err != nil {
			err = fmt.Errorf("记录失败，稍后再试 err=%v", err.Error())
			return
		}
		mread := make(map[string]interface{})
		mread["read_count"] = gorm.Expr("read_count + ?", 1)
		err = global.DB.Model(models.McBook{}).Where("id = ?", bookId).Updates(mread).Error
		if err != nil {
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	} else {
		err = UpdateReadCidByUserId(userId, bookId, chapterId, chapterName, textNum)
		if err != nil {
			return
		}
	}

	if second > 0 {
		today := utils.GetDate()
		count = GetTimeCountByBookId(bookId, userId, today)
		if count <= 0 {
			time := models.McBookTime{
				Uid:     userId,
				Bid:     bookId,
				Second:  0,
				Day:     today,
				Addtime: utils.GetUnix(),
				Uptime:  utils.GetUnix(),
			}
			if err = global.DB.Create(&time).Error; err != nil {
				err = fmt.Errorf("记录失败，请稍后再试 err=%v", err.Error())
				return
			}
		}
		err = UpdateReadTimeByUserId(today, userId, bookId, second)
		if err != nil {
			return
		}

		//判断是否完成阅读任务
		todayUnix := utils.GetTodayUnix()
		seconds := GetTodaySecondsByUserId(userId, todayUnix)

		var tasks []*models.McTask
		tasks, err = GetReadTaskList()
		if err != nil {
			return
		}
		for _, task := range tasks {
			if seconds > task.ReadMinute*60 {
				_ = task_service.CompleteTask(task, userId)
			}
		}
	}
	return
}

func ReadInfo(req *models.ReadInfoReq) (read *models.McBookRead, err error) {
	bookId := req.BookId
	userId := req.UserId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID为空")
		return
	}
	//判断小说阅读记录是否存在
	read, err = GetReadById(bookId, userId)
	if read.Id <= 0 {
		err = fmt.Errorf("%v", "小说不存在")
		return
	}
	if err != nil {
		return
	}
	return
}

func ReadDel(req *models.ReadDelReq) (err error) {
	bookIds := req.BookIds
	userId := req.UserId
	if len(bookIds) <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	err = DeleteReadByBookIds(bookIds, userId)
	return
}

func GetBrowseRes(req *models.BrowseListReq) (readRes *models.BookReadRes, err error) {
	pageNum := req.Page
	pageSize := req.Size
	req.Page = 0
	req.Size = 0
	var today []*models.BookReadListRes
	var yesterday []*models.BookReadListRes
	var agoday []*models.BookReadListRes
	req.Day = utils.Today
	today, _, err = GetBrowseList(req)
	if err != nil {
		return
	}
	req.Day = utils.Yesterday
	yesterday, _, err = GetBrowseList(req)
	if err != nil {
		return
	}

	var total int64
	req.Day = utils.Agoday
	req.Page = pageNum
	req.Size = pageSize
	agoday, total, err = GetBrowseList(req)
	if err != nil {
		return
	}
	readRes = &models.BookReadRes{
		Today:     today,
		Yesterday: yesterday,
		Agoday:    agoday,
		Total:     total,
	}
	return
}

func BrowseDel(req *models.BrowseDelReq) (err error) {
	bookIds := req.BookIds
	userId := req.UserId
	if len(bookIds) <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	err = DeleteBrowseByBookIds(bookIds, userId)
	return
}
