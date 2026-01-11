package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"strings"
)

func getCommentCount(bookId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McComment{}).Where("bid = ?", bookId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 获取书评的基本信息-单条
func GetBookCommentInfo(bookId int64) (commentInfo *models.McBookCommnetInfo) {
	var err error
	err = global.DB.Model(&models.McBookCommnetInfo{}).Where("book_id", bookId).Debug().First(&commentInfo).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	return
}

// 获取书评的多条数据-列表
func GetBookCommnetLists(req *models.BookCommentListReq) (commentBookList []*models.McBookCommnetInfo, err error) {
	db := global.DB.Model(&models.McBookCommnetInfo{}).Order("id desc")
	//主键ID
	Id := req.Id
	if Id != 0 {
		db = db.Where("book_id = ?", Id)
	}
	//书籍ID
	BookId := req.BookId
	if BookId != 0 {
		db = db.Where("book_id = ?", BookId)
	}
	//小说标题
	Title := strings.TrimSpace(req.Title)
	if Title != "" {
		likeName := "%" + Title + "%"
		db = db.Where("title LIKE ?", likeName)
	}
	//根据用户ID查询
	UserId := req.UserId
	if UserId != 0 {
		db = db.Where("user_id = ?", UserId)
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Debug().Find(&commentBookList).Error
	} else {
		err = db.Debug().Find(&commentBookList).Error
	}
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return nil, err
	}
	if len(commentBookList) < 0 {
		return
	}
	return
}

// 获取小说基本信息
// commentBookInfo *models.McBookCommnetInfo
func GetBookDetailInfo(req *models.BookCommentInfoReq) (bookDetailInfo map[string]interface{}, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		return
	}
	sid := int64(bookId)
	commentBookInfo := GetBookCommentInfo(sid)
	if commentBookInfo == nil {
		err = fmt.Errorf("%v", "id不正确")
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//重新定义一个map对象
	bookDetailInfo = make(map[string]interface{})
	bookDetailInfo["id"] = commentBookInfo.Id
	bookDetailInfo["book_id"] = bookId
	bookDetailInfo["title"] = commentBookInfo.Title
	bookDetailInfo["author"] = commentBookInfo.Author
	bookDetailInfo["book_url"] = commentBookInfo.BookUrl
	bookDetailInfo["comment_count"] = commentBookInfo.CommentCount
	bookDetailInfo["cover_logo"] = commentBookInfo.CoverLogo
	bookDetailInfo["category"] = commentBookInfo.Category
	bookDetailInfo["score"] = commentBookInfo.Score
	bookDetailInfo["neary_time"] = commentBookInfo.NearyTime
	bookDetailInfo["addtime"] = commentBookInfo.Addtime
	bookDetailInfo["utime"] = commentBookInfo.Uptime
	//查询书评的关联列表
	commentList := GetALlCoommentBookList(sid)
	bookDetailInfo["comment_list"] = commentList
	return

	//return fieldParams, nil
}

// 根据ID 获取获取所有的书评信息 ==和上面的GetBookCommnetLists类似,
func GetALlCoommentBookList(bookId int64) (commentBookList []*models.McBookCommentDetail) {
	var err error
	err = global.DB.Model(&models.McBookCommentDetail{}).Where("book_id", bookId).Order("id asc").Debug().Find(&commentBookList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	return
}

// 获取精彩书评列表
func GetWonderfulCommentList(req *models.BookCommnetUserReq) (commentBookList []*models.BookCommentListRes, err error) {
	pageSize := req.PageSize

	var typeList []*models.BookCommentListRes
	//先找book_id排序，然后再根据ID正序排列保证输出顺序
	db := global.DB.Model(&models.McBookCommentDetail{}).Order("score desc")
	if pageSize == 0 {
		err = db.Debug().Find(&typeList).Error
	} else {
		err = db.Debug().Limit(pageSize).Find(&typeList).Error
	}
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if len(typeList) == 0 {
		return
	}
	//遍历数据信息
	for _, val := range typeList {
		bookId := int64(val.BookId)
		// 如果没有bookId就不用去处理了
		if bookId == 0 {
			continue
		}
		//获取书的详情信息
		bookInfo := GetBookCommentInfo(bookId)
		checkdays := &models.BookCommentDetailRes{
			Title:     bookInfo.Title,     //小说标题
			Author:    bookInfo.Author,    //小说作者
			CoverLogo: bookInfo.CoverLogo, //小说封面
		}
		val.BookList = append(val.BookList, checkdays)
	}
	return typeList, nil
}

// 根据用户获取书评信息
func GetBookCommentByUid(user_id int64) (commentList []*models.BookCommentListRes) {
	var err error
	db := global.DB.Model(&models.McBookCommentDetail{}).Where("user_id = ?", user_id)
	//查询关联的用户信息
	err = db.Debug().Find(&commentList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if len(commentList) == 0 {
		return
	}
	for _, val := range commentList {
		bookId := int64(val.BookId)
		if bookId == 0 {
			continue
		}
		bookInfo := GetBookCommentInfo(bookId)
		checkdays := &models.BookCommentDetailRes{
			Title:     bookInfo.Title,     //小说标题
			Author:    bookInfo.Author,    //小说作者
			CoverLogo: bookInfo.CoverLogo, //小说封面
		}
		val.BookList = append(val.BookList, checkdays)
	}
	return
}

// 关联查询对应的分页页码信息
func GetCommentUserList(req *models.BookCommnetUserReq) (userList []*models.BookCommentUserListRes, err error) {
	pageSize := req.PageSize
	db := global.DB.Model(&models.McBookCommentDetail{}).Select("count(user_id) as num,user_id,username,avtar_url").Group("user_id")
	db = db.Having("num>1").Order("num desc") //使用字句查询，利用having来进行拼接
	if pageSize == 0 {
		//查询所有的
		err = db.Debug().Find(&userList).Error
	} else {
		//根据limit限制查询
		err = db.Debug().Limit(pageSize).Find(&userList).Error
	}
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if len(userList) <= 0 {
		return
	}
	return
}
