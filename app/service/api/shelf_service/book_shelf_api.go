package shelf_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/task_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/global"
	"go-novel/utils"
)

func GetShelfList(req *models.BookShelfListReq) (shelfs []*models.BookShelfListRes, total, seconds int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	var list []*models.McBookShelf
	db := global.DB.Model(&models.McBookShelf{}).Order("top desc,uptime desc")

	db = db.Where("uid = ?", userId)

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
		var shelf *models.BookShelfListRes
		shelf, err = getShelfBook(val.Bid, userId)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			continue
		}
		shelf.Top = val.Top
		if shelf.ChapterNum <= 0 {
			chapterFile, _ := chapter_service.GetChapterFile(shelf.BookName, shelf.Author)
			var chapters []*models.McBookChapter
			chapters, _ = chapter_service.GetChaptersByFile(chapterFile)
			shelf.ChapterNum = len(chapters)
		}
		shelfs = append(shelfs, shelf)
	}
	startWeekUnix := utils.GetWeekyUnix()
	seconds = getShelfSecondByUserId(userId, startWeekUnix)
	return
}

func getShelfBook(bookId, userId int64) (shelf *models.BookShelfListRes, err error) {
	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}
	lastChapterId, lastChapterName := chapter_service.GetBookNewChapterId(book.BookName, book.Author)
	shelf = &models.BookShelfListRes{
		Uid:             userId,
		Author:          book.Author,
		Bid:             book.Id,
		BookName:        book.BookName,
		Pic:             utils.GetFileUrl(book.Pic),
		ReadChapterId:   1,
		NewsChapterId:   lastChapterId,
		NewsChapterName: lastChapterName,
		ChapterNum:      book.ChapterNum,
		Serialize:       book.Serialize,
	}
	shelf.ReadChapterId, shelf.ReadChapterName = book_service.GetReadChapterIdByBookId(bookId, userId)
	shelf.ReadChapterTextNum = book_service.GetReadTextNumByBookId(bookId, userId)
	return
}

func BookShelfAdd(req *models.BookShelfAddReq) (err error) {
	bookId := req.BookId
	userId := req.UserId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	var count int64
	count = GetShelfCountByBookId(bookId, userId)
	if count > 0 {
		err = fmt.Errorf("%v", "已在书架")
		return
	}

	shelf := models.McBookShelf{
		Bid:     bookId,
		Uid:     userId,
		Addtime: utils.GetUnix(),
		Uptime:  utils.GetUnix(),
	}
	if err = global.DB.Create(&shelf).Error; err != nil {
		err = fmt.Errorf("加入失败，稍后再试 err=%v", err.Error())
		return
	}
	err = UpdateShitsByBookId(userId)
	if err != nil {
		return
	}
	var taskId int64 = 3
	task, err := task_service.GetTaskById(taskId)
	if err != nil {
		return
	}
	err = task_service.CompleteTask(task, userId)
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		return
	}
	return
}

func BookFavDel(req *models.BookShelfDelReq) (err error) {
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
	var res bool
	res, err = DeleteBookShelf(bookIds, userId)
	if err != nil {
		return
	}
	if !res {
		err = fmt.Errorf("%v", "删除失败")
		return
	}
	return
}

func BookShelfTop(req *models.BookShelfTopReq) (err error) {
	userId := req.UserId
	bookId := req.BookId
	isTop := req.Top
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不正确")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["top"] = isTop
	mapData["uptime"] = utils.GetUnix()

	if err = global.DB.Model(models.McBookShelf{}).Where("uid = ? and bid = ?", userId, bookId).Updates(&mapData).Error; err != nil {
		return
	}
	return
}

func IsBookShelf(req *models.IsBookShelfReq) (isExist int, err error) {
	userId := req.UserId
	bookId := req.BookId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	count := GetShelfCountByBookId(bookId, userId)
	if count > 0 {
		isExist = 1
	}
	return
}
