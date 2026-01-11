package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/common/book_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func ChapterBuy(req *models.ChapterBuyReq) (err error) {
	bookId := req.BookId
	userId := req.UserId
	chapterId := req.ChapterId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if chapterId <= 0 {
		err = fmt.Errorf("%v", "章节ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登陆")
		return
	}
	if err != nil {
		return
	}
	//章节表
	chapterTable, err := book_service.GetChapterTable(bookId)
	if err != nil {
		err = fmt.Errorf("%v", "获取章节失败")
		return
	}

	chapter, err := GetChapterByChapterId(chapterTable, chapterId)
	if err != nil {
		return
	}
	if chapter.Id <= 0 {
		err = fmt.Errorf("%v", "该章节不存在")
		return
	}
	if chapter.Cion <= 0 {
		err = fmt.Errorf("%v", "该章节为免费章节")
		return
	}
	user, err := user_service.GetUserById(userId)
	if err != nil {
		return
	}
	if user.Id <= 0 {
		err = fmt.Errorf("%v", "账号异常")
		return
	}
	if user.Cion < chapter.Cion {
		err = fmt.Errorf("%v", "金币不足，请充值")
		return
	}
	err = appBookBuy(chapter, user, req)
	return
}

func appBookBuy(chapter *models.McBookChapter, user *models.McUser, req *models.ChapterBuyReq) (err error) {
	bookId := req.BookId
	bookName := book_service.GetBookNameById(bookId)
	userId := user.Id
	userCion := user.Cion
	chapterId := chapter.Id
	chapterName := chapter.ChapterName
	chapterCion := chapter.Cion
	ip := req.Ip
	auto := req.Auto
	if userCion < chapterCion {
		return
	}

	var count int64
	count = GetBookBuyCountByUnique(userId, bookId, chapterId)
	if count > 0 {
		err = fmt.Errorf("%v", "购买记录已存在")
		return
	}

	tx := global.DB.Begin()

	us := make(map[string]interface{})
	us["cion"] = gorm.Expr("cion - ?", chapterCion)

	err = tx.Model(models.McUser{}).Where("id = ?", user.Id).Updates(us).Error
	if err != nil {
		tx.Rollback()
		return
	}

	//写入消费记录
	buy := models.McBuy{
		Uid:     userId,
		Text:    fmt.Sprintf("购买小说《%v》章节《%v》", bookName, chapterName),
		Cion:    chapterCion,
		Bid:     bookId,
		Cid:     chapterId,
		Ip:      ip,
		Addtime: utils.GetUnix(),
	}
	if err = tx.Create(&buy).Error; err != nil {
		global.Sqllog.Errorf("写入消费记录失败 err=%v", err.Error())
		tx.Rollback()
		return
	}

	//写入购买记录
	bookBuy := models.McBookBuy{
		Uid:  userId,
		Bid:  bookId,
		Cid:  chapterId,
		Auto: auto,
	}
	if err = tx.Create(&bookBuy).Error; err != nil {
		global.Sqllog.Errorf("写入购买记录失败 err=%v", err.Error())
		tx.Rollback()
		return
	}

	tx.Commit()

	//改变所有购买模式
	err = UpdateBookBuyAuto(bookId, auto)
	if err != nil {
		return
	}

	//分成记录
	//if authorId > 0 && authorId != userId {
	//	income := models.McIncome{
	//		Uid:     userId,
	//		Text:    fmt.Sprintf("收到小说《%v》章节《%v》购买分成", bookName, chapterName),
	//		Bid:     bookId,
	//		Cion:    utils.GetCion(chapterCion, utils.Author_Fc_Book),
	//		Zcion:   chapterCion,
	//		Addtime: utils.GetUnix(),
	//	}
	//	if err = global.DB.Create(&income).Error; err != nil {
	//		global.Sqllog.Errorf("写入收入记录失败 err=%v", err.Error())
	//		return
	//	}
	//
	//	//增加收入
	//	xrmb := utils.GetCion(chapterCion/utils.Pay_Rmb_Cion, utils.Author_Fc_Book)
	//
	//	count = getUserCountById(authorId)
	//	if count > 0 && xrmb > 0 {
	//		global.DB.Model(models.McUser{}).Where("id = ?", authorId).Update("rmb", gorm.Expr("rmb + ?", xrmb))
	//	}
	//}
	return
}

// 判断小说收费，返回说明：0可以浏览，-1需要登陆，1需要购买VIP，2需要金币购买
func appBookPay(chapter *models.McBookChapter, req *models.ChapterReadReq) (costType int, err error) {
	bookId := req.BookId
	chapterId := chapter.Id
	chapterCion := chapter.Cion
	chapterVip := chapter.Vip
	userId := req.UserId
	ip := req.Ip

	if userId <= 0 {
		costType = -1
		return
	}

	user, _ := user_service.GetUserById(userId)
	userVip := user.Vip

	if chapterCion > 0 || chapterVip > 0 {
		if chapterVip > 0 && userVip == 0 {
			costType = 1
			return
		}
		var count int64
		var auto int
		if chapterCion > 0 {
			count = GetBookBuyCountByUnique(userId, bookId, chapterId)
			if count <= 0 {
				auto = GetBookBuyAutoByUnique(userId, bookId)
				if auto == 1 {
					if user.Cion < chapterCion {
						costType = 2
					}
					chapterBuy := &models.ChapterBuyReq{
						BookId:    bookId,
						ChapterId: chapterId,
						UserId:    userId,
						Auto:      auto,
						Ip:        ip,
					}
					err = appBookBuy(chapter, user, chapterBuy)
					if err != nil {
						costType = 3
					}
				} else {
					costType = 2
				}
			}
		}
	}
	return
}

func GetRank(bookId int64, sort string) (rowNo string, err error) {
	// 构建原始 SQL 查询，添加 WHERE id = 1 的条件
	//sql := fmt.Sprintf("%v", "SELECT bk.id, bk.cion, bk.rank FROM ( SELECT id, cion, (@rank := @rank + 1) AS rank FROM mc_book, (SELECT @rank := 0) r ORDER BY cion DESC ) AS bk WHERE bk.id = 1")
	var rank models.GetRankRes
	mcTab := fmt.Sprintf("%v", new(models.McBook).TableName())
	tableName := fmt.Sprintf("( SELECT id, cion, (@rank := @rank + 1) AS rank FROM %v, (SELECT @rank := 0) r ORDER BY %v DESC )", mcTab, sort)
	sql := fmt.Sprintf("SELECT bk.id, bk.cion, bk.rank FROM %v AS bk WHERE bk.id = %v", tableName, bookId)
	// 执行原始 SQL 查询
	err = global.DB.Raw(sql).Debug().Scan(&rank).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if rank.Rank > 0 {
		rowNo = fmt.Sprintf("%v", rank.Rank)
	} else {
		rowNo = "100名以外"
	}
	return
}
