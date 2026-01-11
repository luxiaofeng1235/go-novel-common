package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
)

func GetBookById(id int64) (book *models.McBook, err error) {
	err = global.DB.Model(models.McBook{}).Where("id", id).First(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 删除推荐数据信息
* @param package_id int 包ID
* @return res 结果集合 ,err 错误信息
 */
func DeleteBookRec(recommandType string) (res bool, err error) {
	if recommandType == "" {
		return
	}
	db := global.DB.Debug().Model(models.McBookRecommand{}).Where("recommand_type", recommandType)
	var total int64

	err = db.Count(&total).Error
	if total > 0 {
		fmt.Println("***************************删除推荐活动数据信息**************************")
		err = db.Delete(&models.McBookRecommand{}).Error
		if err != nil {
			return res, err
		}
		return res, nil
	} else {
		return res, nil
	}
}

// 设置首页推荐数据
func SetBookRecData(req *models.CreateBookRecommandReq) (res bool, err error) {
	//类型判断
	var recommandType = req.RecommandType
	if recommandType == "" {
		err = fmt.Errorf("%v", "设置推荐数据不能为空")
		return
	}
	//里面的接口数据判断
	var recommandIds = req.RecommandIds //设置的推荐ID信息
	if len(recommandIds) == 0 {
		err = fmt.Errorf("%v", "设置推荐数据不能为空")
		return
	}
	//fmt.Println(recommandType)

	_, err = DeleteBookRec(recommandType) // 每次提交删除冗余的推荐数据
	var myRecBookList []models.MySyncRecList
	insertData := make([]map[string]interface{}, 0, len(myRecBookList))
	for _, value := range recommandIds {
		//fmt.Println(value.BookId)
		insertData = append(insertData, map[string]interface{}{
			"recommand_type": recommandType,   //推荐类型
			"book_id":        value.BookId,    //小说ID
			"addtime":        utils.GetUnix(), //添加时间
		})
	}
	fmt.Println(insertData)
	err = global.DB.Model(models.McBookRecommand{}).Debug().CreateInBatches(&insertData, 100).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

// 通过类型进行检索相关类型信息
func GetBookRecByType(recommandType string) (recList []int64, err error) {
	if recommandType == "" {
		return
	}
	var list []*models.McBookRecommand
	db := global.DB.Debug().Model(models.McBookRecommand{}).Order("id asc").Debug()
	db = db.Where("recommand_type = ?", recommandType)
	err = db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	var myIds []int64 //定义对象类型进行追加
	for _, value := range list {
		myIds = append(myIds, value.BookId)
	}
	fmt.Println(myIds)
	return myIds, nil
}

// 获取推荐类型
func GetBookRecList() (bookDetailInfo map[string]interface{}, err error) {

	//重新定义接口数据信息
	bookDetailInfo = make(map[string]interface{})
	//'推荐类型  rec_serialize 推荐完结  -  rec_new 推荐新书 - hot_serialize 热门完结 - hot_rank 热门排行 - hot_new 热门新书 - hot_search 热门搜索 - classic_search 经典热搜 - classic_hight  经典高分 - classic_rq 经典人气 - classic_serialize 经典完结 - classic_new 经典新书',
	//  `book_id` int(11) DEFAULT '0' COMMENT '小说ID',
	recSerializeList, _ := GetBookRecByType("rec_serialize")
	recNewList, _ := GetBookRecByType("rec_new")
	hotSerializeList, _ := GetBookRecByType("hot_serialize")
	hotRankList, _ := GetBookRecByType("hot_rank")
	hotNewList, _ := GetBookRecByType("hot_new")
	hotSearchList, _ := GetBookRecByType("hot_search")
	classicSearchList, _ := GetBookRecByType("classic_search")
	classicHighList, _ := GetBookRecByType("classic_hight")
	classicRqList, _ := GetBookRecByType("classic_rq")
	classicSerializeList, _ := GetBookRecByType("classic_serialize")
	classicNewList, _ := GetBookRecByType("classic_new")

	bookDetailInfo["rec_serialize"] = recSerializeList //推荐完结
	bookDetailInfo["rec_new"] = recNewList             //推荐新书

	bookDetailInfo["hot_serialize"] = hotSerializeList //热门完结
	bookDetailInfo["hot_new"] = hotNewList             //热门新书
	bookDetailInfo["hot_search"] = hotSearchList       //热门搜索
	bookDetailInfo["hot_rank"] = hotRankList           //热门排行

	bookDetailInfo["classic_search"] = classicSearchList       //经典搜索
	bookDetailInfo["classic_hight"] = classicHighList          //经典高分
	bookDetailInfo["classic_rq"] = classicRqList               //经典人气
	bookDetailInfo["classic_serialize"] = classicSerializeList //经典完结
	bookDetailInfo["classic_new"] = classicNewList             //经典新书

	return bookDetailInfo, nil
}

func BookListPassSearch(req *models.BookListPassReq) (list []*models.BookListAdminRes, total int64, err error) {
	db := global.DB.Model(&models.McBook{}).Debug().Select("id,book_name,author,update_chapter_title,source_url")
	bookName := strings.TrimSpace(req.BookName) //书籍ID搜索
	db = db.Where("status = 1")                 //只查询为1的
	if bookName != "" {
		db = db.Where("book_name LIKE ?", "%"+bookName+"%")
	}
	cid := strings.TrimSpace(req.Cid)
	if cid != "" { //分类ID搜索
		db = db.Where("cid = ?", cid)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}
	return
}

func BookListSearch(req *models.BookListReq) (list []*models.McBook, total int64, err error) {
	db := global.DB.Model(&models.McBook{}).Order("id desc")

	bookId := strings.TrimSpace(req.BookId)
	if bookId != "" {
		db = db.Where("id = ?", bookId)
	}

	bookName := strings.TrimSpace(req.BookName)
	if bookName != "" {
		db = db.Where("book_name LIKE ?", "%"+bookName+"%")
	}

	author := strings.TrimSpace(req.Author)
	if author != "" {
		db = db.Where("author = ?", author)
	}

	sourceUrl := strings.TrimSpace(req.SourceUrl)
	if sourceUrl != "" {
		db = db.Where("source_url LIKE ?", "%"+sourceUrl+"%")
	}

	serialize := strings.TrimSpace(req.Serialize)
	if serialize != "" {
		db = db.Where("serialize = ?", serialize)
	}

	isRec := strings.TrimSpace(req.IsRec)
	if isRec != "" {
		db = db.Where("is_rec = ?", isRec)
	}

	isHot := strings.TrimSpace(req.IsHot)
	if isHot != "" {
		db = db.Where("is_hot = ?", isHot)
	}

	isNew := strings.TrimSpace(req.IsNew)
	if isNew != "" {
		db = db.Where("is_new = ?", isNew)
	}

	isChoice := strings.TrimSpace(req.IsChoice)
	if isChoice != "" {
		db = db.Where("is_choice = ?", isChoice)
	}

	isClassic := strings.TrimSpace(req.IsClassic)
	if isClassic != "" {
		db = db.Where("is_classic = ?", isClassic)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	isLess := strings.TrimSpace(req.IsLess)
	if isLess != "" {
		db = db.Where("is_less = ?", isLess)
	}

	cid := strings.TrimSpace(req.Cid)
	if cid != "" {
		db = db.Where("cid = ?", cid)
	}

	tid := strings.TrimSpace(req.Tid)
	if tid != "" {
		db = db.Where("tid = ?", tid)
	}

	if req.SourceDay > 0 {
		sourceDayUnix := utils.GetAgoDayUnix(req.SourceDay)
		db = db.Where("last_chapter_time >= ?", sourceDayUnix)
	}

	if req.SourceNoday > 0 {
		sourceNoDayUnix := utils.GetAgoDayUnix(req.SourceNoday)
		db = db.Where("last_chapter_time <= ?", sourceNoDayUnix)
	}

	if req.RecentlyDay > 0 {
		recentlyUnix := utils.GetAgoDayUnix(req.RecentlyDay)
		db = db.Where("update_chapter_time >= ?", recentlyUnix)
	}

	if req.EecentlyNoday > 0 {
		recentlyNoUnix := utils.GetAgoDayUnix(req.EecentlyNoday)
		db = db.Where("update_chapter_time <= ?", recentlyNoUnix)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}
	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, book := range list {
		log.Printf("pic的数据库路径 pic = %v", book.Pic)
		book.Pic = utils.GetAdminFileUrl(book.Pic)
	}
	return
}

func CreateBook(req *models.CreateBookReq) (InsertId int64, err error) {
	bookType := req.BookType
	if bookType <= 0 {
		err = fmt.Errorf("%v", "阅读类型不能为空")
		return
	}

	bookName := req.BookName
	if bookName == "" {
		err = fmt.Errorf("%v", "小说名称不能为空")
		return
	}

	pic := req.Pic
	if pic == "" {
		err = fmt.Errorf("%v", "小说封面不能为空")
		return
	}

	cid := req.Cid
	if cid <= 0 {
		err = fmt.Errorf("%v", "小说分类不能为空")
		return
	}
	isRec := req.IsRec
	isHot := req.IsHot
	isChoice := req.IsChoice
	isClassic := req.IsClassic
	isNew := req.IsNew
	isTeen := req.IsTeen
	serialize := req.Serialize
	if serialize <= 0 {
		err = fmt.Errorf("%v", "小说连载状态不能为空")
		return
	}
	author := req.Author
	if author != "" {
		err = fmt.Errorf("%v", "小说作者不能为空")
		return
	}
	tags := req.Tags
	desc := req.Desc
	if desc != "" {
		err = fmt.Errorf("%v", "小说简介不能为空")
		return
	}
	textNum := req.TextNum
	if textNum <= 0 {
		err = fmt.Errorf("%v", "小说总字数不能为空")
		return
	}
	hits := req.Hits
	hitsMonth := req.HitsMonth
	hitsWeek := req.HitsWeek
	hitsDay := req.HitsDay
	shits := req.Shits
	isPay := req.IsPay
	chapterNum := req.ChapterNum
	if chapterNum <= 0 {
		err = fmt.Errorf("%v", "章节总数不能为空")
		return
	}
	score := req.Score
	sourceUrl := req.SourceUrl
	sourceId := req.SourceId
	readCount := req.ReadCount
	searchCount := req.SearchCount
	book := models.McBook{
		BookType:    bookType,
		BookName:    bookName,
		Pic:         pic,
		IsRec:       isRec,
		IsHot:       isHot,
		IsChoice:    isChoice,
		IsClassic:   isClassic,
		IsNew:       isNew,
		IsTeen:      isTeen,
		Serialize:   serialize,
		Author:      author,
		Tags:        tags,
		Desc:        desc,
		TextNum:     textNum,
		Hits:        hits,
		HitsMonth:   hitsMonth,
		HitsWeek:    hitsWeek,
		HitsDay:     hitsDay,
		Shits:       shits,
		IsPay:       isPay,
		ChapterNum:  chapterNum,
		Score:       score,
		SourceUrl:   sourceUrl,
		SourceId:    sourceId,
		ReadCount:   readCount,
		SearchCount: searchCount,
		Addtime:     utils.GetUnix(),
	}

	if err = global.DB.Create(&book).Error; err != nil {
		return 0, err
	}

	return book.Id, nil
}

func UpdateBook(req *models.UpdateBookReq) (res bool, err error) {
	id := req.BookId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	var mapData = make(map[string]interface{})
	bookType := req.BookType
	if bookType > 0 {
		mapData["book_type"] = bookType
	}
	bookName := strings.TrimSpace(req.BookName)
	if bookName != "" {
		mapData["book_name"] = bookName
	}
	pic := strings.TrimSpace(req.Pic)
	if pic != "" {
		mapData["pic"] = pic
	}
	cid := req.Cid
	if cid > 0 {
		mapData["cid"] = cid
	}
	isRec := req.IsRec
	if isRec > 0 {
		mapData["is_rec"] = isRec
	} else {
		mapData["is_rec"] = 0 //默认推荐为0
	}
	isHot := req.IsHot
	if isHot > 0 {
		mapData["is_hot"] = isHot
	} else {
		mapData["is_hot"] = 0 //默认热门为0
	}
	isChoice := req.IsChoice
	if isChoice > 0 {
		mapData["is_choice"] = isChoice
	} else {
		mapData["is_choice"] = 0 //默认精选为0
	}
	isClassic := req.IsClassic
	if isClassic > 0 {
		mapData["is_classic"] = isClassic
	} else {
		mapData["is_classic"] = 0 //默认经典为0
	}
	isNew := req.IsNew
	if isNew > 0 {
		mapData["is_new"] = isNew
	} else {
		mapData["is_new"] = 0 //默认最新为0
	}
	isTeen := req.IsTeen
	if isTeen > 0 {
		mapData["is_teen"] = isTeen
	} else {
		mapData["is_teen"] = 0 //默认青少年为0
	}
	serialize := req.Serialize
	if serialize > 0 {
		mapData["serialize"] = serialize
	}
	author := strings.TrimSpace(req.Author)
	if author != "" {
		mapData["author"] = author
	}
	tags := strings.TrimSpace(req.Tags)
	if tags != "" {
		mapData["tags"] = tags
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		mapData["desc"] = desc
	}
	textNum := req.TextNum
	if textNum > 0 {
		mapData["text_num"] = req.TextNum
	}
	hits := req.Hits
	if hits > 0 {
		mapData["hits"] = req.Hits
	}
	hitsMonth := req.HitsMonth
	if hitsMonth > 0 {
		mapData["hits_month"] = req.HitsMonth
	}
	hitsWeek := req.HitsWeek
	if hitsWeek > 0 {
		mapData["hits_week"] = req.HitsWeek
	}
	hitsDay := req.HitsDay
	if hitsDay > 0 {
		mapData["hits_day"] = hitsDay
	}
	shits := req.Shits
	if shits > 0 {
		mapData["shits"] = req.Shits
	}
	isPay := req.IsPay
	if isPay > 0 {
		mapData["is_pay"] = req.IsPay
	}
	chapterNum := req.ChapterNum
	if chapterNum > 0 {
		mapData["chapter_num"] = req.ChapterNum
	}
	score := req.Score
	if score > 0 {
		mapData["score"] = req.Score
	}
	sourceUrl := strings.TrimSpace(req.SourceUrl)
	if sourceUrl != "" {
		mapData["source_url"] = sourceUrl
	}
	sourceId := req.SourceId
	if sourceId > 0 {
		mapData["source_id"] = sourceId
	}
	readCount := req.ReadCount
	if readCount > 0 {
		mapData["read_count"] = readCount
	}
	searchCount := req.SearchCount
	if searchCount > 0 {
		mapData["search_count"] = searchCount
	}
	mapData["uptime"] = utils.GetUnix()
	isIndex := req.IsIndex //表明是从首页进入的
	if isIndex > 0 {
		//首页进入不做任何处理
	} else {
		mapData["status"] = req.Status //修改状态
	}
	if err = global.DB.Model(models.McBook{}).Where("id", id).Debug().Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteBook(req *models.DeleteBookReq) (res bool, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说id错误")
		return
	}
	book, err := GetBookById(bookId)
	if err != nil {
		return
	}
	bookName := book.BookName
	author := book.Author
	var uploadBookChapterPath string
	uploadBookChapterPath, err = setting_service.GetValueByName(utils.UploadBookChapterPath)
	if err != nil {
		err = fmt.Errorf("获取小说章节失败 uploadBookChapterPath=%v", uploadBookChapterPath)
		return
	}
	chapterFile, err := chapter_service.GetChapterFile(bookName, author)
	err = utils.RemoveFile(chapterFile)
	if err != nil {
		return
	}
	var uploadBookTextPath string
	uploadBookTextPath, err = setting_service.GetValueByName(utils.UploadBookTextPath)
	if err != nil {
		err = fmt.Errorf("获取小说内容目录失败 uploadBookTextPath=%v", uploadBookTextPath)
		return
	}
	txtDir, err := chapter_service.GetTxtDir(bookName, author)
	if err != nil {
		return
	}
	err = utils.RemoveDir(txtDir)
	if err != nil {
		return
	}
	err = global.DB.Where("bid = ?", bookId).Delete(&models.McBookShelf{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	err = global.DB.Where("bid = ?", bookId).Delete(&models.McBookBrowse{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	err = global.DB.Where("bid = ?", bookId).Delete(&models.McBookRead{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	err = global.DB.Where("id = ?", bookId).Delete(&models.McBook{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	return true, nil
}
