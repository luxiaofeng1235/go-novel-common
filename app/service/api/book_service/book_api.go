package book_service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/api/version_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/user_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"sort"
	"strings"
)

func Info(req *models.BookInfoReq) (rbook *models.BookInfoRes, err error) {
	bookId := req.BookId
	userId := req.UserId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID为空")
		return
	}
	//判断小说是否存在
	book, err := book_service.GetBookById(bookId)
	if book.Id <= 0 {
		err = fmt.Errorf("%v", "小说不存在")
		return
	}
	if err != nil {
		return
	}
	rbook = &models.BookInfoRes{
		BookId:    bookId,
		BookName:  book.BookName,
		Author:    book.Author,
		Desc:      book.Desc,
		Serialize: book.Serialize,
		Pic:       utils.GetFileUrl(book.Pic),
		TextNum:   book.TextNum,
		Hits:      book.Hits,
		Score:     book.Score,
		IsPay:     book.IsPay,
		Addtime:   book.Addtime,
	}

	var count int64
	if userId > 0 {
		count = book_service.GetShelfCountByBookId(bookId, userId)
		if count > 0 {
			rbook.IsShelf = 1
		}
		rbook.ReadChapterId, rbook.ReadChapterName = book_service.GetReadChapterIdByBookId(bookId, userId)
	}
	rbook.CommentCount = getCommentCount(bookId)
	rbook.ReadCount = book.ReadCount

	newChapterName := chapter_service.GetBookNewChapterName(book.BookName, book.Author)
	rbook.NewChapterName = newChapterName

	if rbook.TextNum <= 0 {
		chapterFile, _ := chapter_service.GetChapterFile(book.BookName, book.Author)
		var chapters []*models.McBookChapter
		chapters, _ = chapter_service.GetChaptersByFile(chapterFile)
		rbook.TextNum = len(chapters) * 2000
	}

	var books []*models.BookInfoHighScoreRes
	commentList := getCommentByBookId(bookId, userId)
	rbook.CommentList = commentList
	books, err = GetHighScoreBook(bookId, book.Cid, book.Tid, userId, 3, "", "", "", "")
	rbook.Scores = books

	count = getBrowseCountByBookId(bookId, userId)
	if count <= 0 {
		read := models.McBookBrowse{
			Bid:     bookId,
			Uid:     userId,
			Addtime: utils.GetUnix(),
			Uptime:  utils.GetUnix(),
		}
		if err = global.DB.Create(&read).Error; err != nil {
			global.Sqllog.Errorf("记录失败，稍后再试 err=%v", err.Error())
			err = nil
			return
		}
	} else {
		err = updateBrowseByUserId(userId, bookId)
		if err != nil {
			global.Sqllog.Errorf("更新失败 err=%v", err.Error())
			return
		}
	}
	err = updateShelfByUserId(userId, bookId)
	if err != nil {
		global.Sqllog.Errorf("更新失败 err=%v", err.Error())
		return
	}
	err = book_service.UpdateHitsByBookName(book.BookName, book.Author)
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetHighScoreBook(bookId, cid, tid, userId int64, limit int, device_type string, package_name string, ip string, mark string) (books []*models.BookInfoHighScoreRes, err error) {
	key := fmt.Sprintf("%v_%v_%v", userId, cid, tid)
	pageNum := utils.HighScoreBookPage[key]
	fmt.Println("页码", pageNum)
	pageNum++
	utils.HighScoreBookPage[key] = pageNum

	fmt.Println("11111111111111111111", ip)
	db := global.DB.Model(models.McBook{}).Debug()
	bookStatus, err := GetBookCopyright(device_type, package_name, ip, mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	var condition string
	if bookStatus == 1 {
		//查询版权
		condition = "status = 1 and is_banquan = 1"
	} else {
		//所有的书籍
		condition = "status = 1"

	}
	db = db.Where(condition)

	//随机排序规则
	//if cid > 0 {
	//	db = db.Where("cid = ?", cid)
	//}
	//if tid > 0 {
	//	db = db.Where("tid = ?", tid)
	//}
	condition = condition + " and score >= 8"
	db = db.Where(condition)
	if limit <= 0 {
		limit = 3
	}
	pageSize := limit
	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&books).Error
	if len(books) <= 0 {
		utils.HighScoreBookPage[key] = 0
		return
	}
	for _, book := range books {
		book.Pic = utils.GetFileUrl(book.Pic)
	}
	return
}

func Chapter(bookId, sourceId int64, sortStatus string) (chapters []*models.McBookChapter, bookUrl, chapterTextReg string, err error) {
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}

	if sourceId > 0 {
		var collect *models.McCollect
		collect, err = getCollectById(sourceId)
		if err != nil {
			err = fmt.Errorf("%v", "访问书源出错啦")
			return
		}
		if collect.Id <= 0 {
			err = fmt.Errorf("%v", "该书源不存在")
			return
		}
		chapterTextReg = collect.ChapterTextReg
		bookUrl = getBookUrlByBookId(sourceId, bookId)
		chapters, err = GetSourceChapters(sourceId, bookUrl, sortStatus)
		return
	}
	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}

	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(book.BookName, book.Author)
	if err != nil {
		return
	}
	chapters, err = chapter_service.GetChaptersByFile(chapterFile)
	if err != nil {
		return
	}

	if len(chapters) <= 0 {
		return
	}

	var lastIndex = len(chapters) - 1
	chapters[0].IsFirst = 1
	chapters[lastIndex].IsLast = 1
	if sortStatus == "desc" {
		sort.Slice(chapters, func(i, j int) bool {
			return chapters[i].Sort > chapters[j].Sort
		})
		chapters[0].IsLast = 1
		chapters[lastIndex].IsFirst = 1
	}
	return
}

func ChapterRead(req *models.ChapterReadReq) (readRes *models.ReadInfoRes, err error) {
	bookId := req.BookId
	userId := req.UserId
	chapterId := req.ChapterId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if chapterId <= 0 {
		return
	}
	//判断小说是否存在
	book, err := book_service.GetBookById(bookId)
	if book.Id <= 0 {
		err = fmt.Errorf("%v", "小说不存在")
		return
	}

	var bookName = book.BookName
	var author = book.Author
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		return
	}

	var chapter *models.McBookChapter
	if chapterId > 0 {
		chapter, err = chapter_service.GetChapterByChapterId(gq, chapterId)
	} else {
		chapter, err = chapter_service.GetFirst(gq, "sort")
	}
	if err != nil {
		return
	}
	if chapter.Id <= 0 {
		err = fmt.Errorf("%v", "章节不存在")
		return
	}

	chapterSort := chapter.Sort
	readRes = &models.ReadInfoRes{
		ChapterId:   chapter.Id,
		Bid:         bookId,
		ChapterName: chapter.ChapterName,
		Vip:         chapter.Vip,
		TextNum:     chapter.TextNum,
		Addtime:     chapter.Addtime,
		AudioName:   utils.Md5(fmt.Sprintf("%v%v", "audio", strings.TrimSpace(chapter.ChapterName))), //音频字段处理客户端的请求
	}

	//是否加入书架
	var count int64
	if userId > 0 {
		count = book_service.GetShelfCountByBookId(bookId, userId)
		if count > 0 {
			readRes.IsShelf = 1
		}
	}

	var text string
	_, text, err = book_service.GetBookTxt(bookName, author, chapter.ChapterName, "")
	if err != nil {
		text = "该章缺失txt文件"
	}
	readRes.Text = text

	//上下章ID
	readRes.PrevChapterId = chapter_service.GetChapterPrev(gq, chapterSort)
	readRes.NextChapterId = chapter_service.GetChapterNext(gq, chapterSort)

	firstSort := chapter_service.GetSortFirst(bookName, author)
	lastSort := chapter_service.GetSortLast(bookName, author)
	if chapterSort == firstSort {
		readRes.IsFirst = 1
	}
	if chapterSort == lastSort {
		readRes.IsLast = 1
	}
	return
}

func GetRankList(req *models.ApiRankListReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})

	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1").Debug()
	}
	num := req.Size
	fmt.Println("分页的页码", num)
	if num <= 0 {
		log.Println("当前分页页码为空，会默认配置一个num =20")
		num = 20 //如果没有获取到分页就20条
	}

	bookType := req.BookType
	if bookType > 0 {
		if bookType == 1 || bookType == 2 {
			db = db.Where("book_type = ?", bookType)
		} else {
			db = db.Where("is_classic = 1")
		}
	}

	columnType := req.ColumnType
	if columnType == 1 {
		db = db.Where("is_rec = 1") //推荐书籍
	} else if columnType == 2 {
		db = db.Where("is_hot = 1") //热门书籍
	} else if columnType == 3 {
		db = db.Where("is_classic = 1") //经典书籍
	}

	sort := req.Sort

	//如果是推荐排序就按照推荐的排序进行查询
	if columnType == 3 {
		if sort == utils.Search {
			sortData, _ := GetApiBookRecByType("classic_search") //经典热搜
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.Score {
			sortData, _ := GetApiBookRecByType("classic_hight") //经典高分
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.Hits {
			sortData, _ := GetApiBookRecByType("classic_rq") //经典人气
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.Serialize {
			sortData, _ := GetApiBookRecByType("classic_serialize") //经典完结
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.New {
			sortData, _ := GetApiBookRecByType("classic_new") //经典新书
			if sortData != "" {
				db = db.Order(sortData)
			}
		}
	} else {
		if sort == utils.Rec {
			db = db.Order("is_rec desc,uptime desc") //如果推荐排序一致，按照更新时间排序
		} else if sort == utils.New {
			//热门新书
			sortData, _ := GetApiBookRecByType("hot_new") //热门新书
			//根据热门排行进行推荐
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.Search {
			//热门搜索
			sortData, _ := GetApiBookRecByType("hot_search") //热门搜索
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else if sort == utils.Serialize {
			//热门完结
			sortData, _ := GetApiBookRecByType("hot_serialize") //热门完结
			if sortData != "" {
				db = db.Order(sortData)
			}
		} else {
			//默认热门排行
			sortData, _ := GetApiBookRecByType("hot_rank")
			//默认热门排行进行推荐排序
			if sortData != "" {
				db = db.Order(sortData)
			}
		}
	}

	if sort == utils.Serialize {
		db = db.Where("serialize = 2")
	} else if sort == utils.New {
		db = db.Where("is_new = 1")
	} else if sort == utils.Rec {
		db = db.Where("is_rec = 1")
	} else if sort == utils.Hot {
		db = db.Where("is_hot = 1")
	} else if sort == utils.Hits {
		db = db.Order("hits desc")
	} else if sort == utils.Search {
		db = db.Order("search_count desc")
	} else if sort == utils.Score {
		db = db.Order("score desc")
	}

	//是否为最热
	isHot := req.IsHot
	if isHot == 1 {
		db = db.Where("is_hot = 1")
	}
	//知否为最新
	isNew := req.IsNew
	if isNew == 1 {
		db = db.Where("is_new = 1")
	}

	tid := req.Tid
	if tid > 0 {
		db = db.Where("tid = ?", tid)
	}

	err = db.Limit(num).Find(&list).Error

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

// 为你推荐处理
func GetSectionForYouRec(req *models.SectionForYouRecReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})
	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1").Debug()
	}
	num := req.Size
	if num <= 0 {
		num = 8
	}

	bookType := user_service.GetBookTypeByUserId(req.UserId)
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}
	var total int64
	db = db.Where("is_rec = 1")
	db.Count(&total) //统计总数
	global.Requestlog.Infof("为你推荐的总书籍总数 total = %v", total)
	//获取随机的数量
	offsetNum := utils.GetBookRandPosition(total)
	err = db.Offset(offsetNum).Limit(num).Find(&list).Error

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

// 通过类型进行检索相关类型信息
func GetApiBookRecByType(recommandType string) (sortString string, err error) {
	if recommandType == "" {
		return
	}
	var list []*models.McBookRecommand
	db := global.DB.Model(models.McBookRecommand{}).Debug()
	db = db.Where("recommand_type = ?", recommandType)
	err = db.Find(&list).Error
	if err != nil {
		return "", err
	}
	var myIds []int64 //定义对象类型进行追加
	for _, value := range list {
		//如果没有值为空或者什么就直接退出本次循环
		if value.BookId == 0 {
			continue
		}
		myIds = append(myIds, value.BookId)
	}
	var sortDataString string
	if len(myIds) > 0 {
		sortString = utils.JoinInt64ToString(myIds)
		sortDataString = fmt.Sprintf("CASE WHEN id in(%s) THEN 1 ELSE 999999 END", sortString)
	}
	//特殊为空的条件
	if sortDataString == "" {
		return "", nil
	}
	log.Printf("获取类型recommand_type = 【%v】 对应的推荐的列表数据为: %v\n", recommandType, sortString)
	log.Printf("获取对应的排序语句order by = %s", sortDataString)
	return sortDataString, nil
}

// 获取随机排名的N本书-定时计算
func GetHighRankRand(limit int, rec_type string) (ids []int64, err error) {
	var list []*models.McBookSearchRank
	db := global.DB.Model(models.McBookSearchRank{}).Debug()
	db = db.Where("rec_type = ?", rec_type)
	var total int64
	db.Count(&total) //统计总数
	var newtotal int64
	if total >= 100 {
		newtotal = total - 10 //默认-10保证都有数据去随机
	} else {
		newtotal = total
	}
	global.Requestlog.Infof("高分统计的总次数 total = %v, newTotal = %v", total, newtotal)
	//获取随机的总次数
	offsetNum := utils.GetBookRandPosition(newtotal)
	//如果随机的时候还存在有不够8本的情况，就需要重新随机一次
	//if offsetNum < limit {
	//	offsetNum = utils.GetBookRandPosition(total) //重新请求一次这次需要看下几率问题
	//}
	err = db.Offset(offsetNum).Limit(limit).Find(&list).Error
	if err != nil {
		return
	}
	var myIds []int64 //定义对象类型进行追加
	for _, value := range list {
		//如果没有值为空或者什么就直接退出本次循环
		if value.Bid == 0 {
			continue
		}
		myIds = append(myIds, value.Bid)
	}
	log.Printf("获取高分书籍的TOP8的排行数据为: %v\n", myIds)
	return myIds, nil
}

func GetSectionHighScore(req *models.SectionHighScoreReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})

	fmt.Println("11111111111111111111", req.Ip)
	//获取正版授权的图书
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1  and is_rec= 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1 and is_rec= 1").Debug()
	}
	num := req.Size
	if num <= 0 {
		num = 6 //随机六条
	}
	bookType := user_service.GetBookTypeByUserId(req.UserId)
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}
	//获取高分统计的随机八条
	hightBooksId, err := GetHighRankRand(num, "high")
	if err != nil {
		global.Sqllog.Errorf("GetHighRankRand err:%v", err)
	}
	if len(hightBooksId) > 0 {
		db = db.Where("id in (?)", hightBooksId)
	} else {
		//如果没有查到高分的关联统计数据，就默认用>=8进行查看统计
		db = db.Where("score >=8")
	}
	//判断当前有缓存就从缓存中进行读取
	if len(hightBooksId) > 0 {
		//有高峰佳作就不需要limit，限制查询的语句就可以了
		err = db.Find(&list).Error
	} else {
		log.Println("获取的高峰佳作数据为空，会优先走本身的book表的随机排序")
		//如果没有匹配到相关的关联数据，就用这个里面再随机
		var total int64
		db.Count(&total) //统计总数
		offsetNum := utils.GetBookRandPosition(total)
		//默认查八条
		err = db.Offset(offsetNum).Limit(num).Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

func GetSectionEnd(req *models.SectionEndReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})

	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1").Debug()
	}
	num := req.Size
	if num <= 0 {
		num = 8
	}

	bookType := user_service.GetBookTypeByUserId(req.UserId)
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}

	db = db.Where("is_hot = 1 and serialize = 2")
	var total int64
	db.Count(&total)
	global.Requestlog.Infof("完结的精品好书总数 total = %v", total)
	offsetNum := utils.GetBookRandPosition(total)
	err = db.Offset(offsetNum).Limit(num).Find(&list).Error
	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

/*
* @note 判断是否在数组中的元素信息
* @param items object 切片对象
* @param item integer 元素值
* @return bool
 */
func IsContainCity(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 测试排序
func GetRandTest() (aa string) {
	db := global.DB.Table("mc_book as a").Debug()
	var list []*models.McBook
	withjoin := "join (SELECT round((SELECT max(id) - min(id) FROM mc_book) * rand() + (SELECT min(id) FROM mc_book ) ) AS rid FROM mc_book  WHERE is_rec = 1 LIMIT 8) b ON a.id = b.rid"
	db = db.Joins(withjoin)
	err := db.Select("id").Find(&list).Limit(5)
	if err != nil {

	}
	return "1122"
}

/*
* @note 获取IP的归属地
* @param ip string ip地址
* @return bool
 */
func GetIpString(ip string) (ipstring string) {
	if ip == "" {
		return
	}
	ipstring = utils.GetGeoCiyByIp(ip)
	//通过geno来进行解析
	//判断是否为ipv6
	if strings.Contains(ip, string(":")) != false {
		ipstring = utils.GetGeoCiyByIp(ip) //解析IPv6的地址
	} else {
		ipstring = utils.GetIpDbNameByIp(ip) //用商业版的IP
	}
	log.Printf("ip = %v 解析出来的城市city_name=%v\n", ip, ipstring)

	return
}

// 获取当前的城市列表信息
func GetBookCityList() (cityList []string, err error) {
	var list []*models.McCity
	db := global.DB.Model(&models.McCity{}).Debug()
	//查询当前的城市列表信息
	err = db.Find(&list).Error
	if err != nil {
		return
	}
	if len(list) <= 0 {
		return
	}
	//获取当前的配置的城市列表信息
	var cityStringObject []string
	for _, val := range list {
		cityStringObject = append(cityStringObject, val.CityName)
	}
	log.Printf("cityList的版权城市信息, %v\n", cityStringObject)
	return cityStringObject, nil
}

/*
* @note 获取书的版权信息-通用函数
* @param device_type string  设备类型
* @param package_name string 包名
* @param ip string 客户端IP
* @param mark 渠道号
* @return status  , err
 */
func GetBookCopyright(device_type, package_name, ip, mark string) (status int, err error) {
	//如果端号和包名为空，说明是老的APP直接开启审核模式
	//status :1 审核模式开启模式 0：未审核模式
	log.Printf("接受请求的具体参数信息**************************** device_type=%v package_name=%v mark=%v***************\n", device_type, package_name, mark)
	device_type = strings.TrimSpace(device_type)
	package_name = strings.TrimSpace(package_name)

	mark = strings.TrimSpace(mark)
	//判断端号、包名、渠道号任何一个为空就返回版权页

	log.Printf("*****请求的客户端IP：%v", ip)
	//通过设备号+包名+渠道来关联
	versionInfo, err := version_service.GetVersionByQdh(device_type, package_name, mark)
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return 0, nil
	}
	log.Printf("获取到的渠道包信息:%+v\n", versionInfo)
	//获取版本获取为空，返回0
	if versionInfo == nil {
		return 0, nil
	}
	if versionInfo.Id != 0 {
		copyright := versionInfo.CopyrightStatus
		if copyright == 1 { //审核流程

			//只有开启版权判断才进行判断
			log.Printf("后台配置的审核状态开启，进入审核流程，接下来判断地区是否满足")
			if device_type == "" || package_name == "" || mark == "" {
				log.Printf("当前没有包名或端号，默认流程为开启审核")
				return 1, nil //如果没有传默认为免审核模式
			}
			//只有在这里面显示正版
			//cityList := []string{"郑州市", "北京市", "上海市", "广州市", "深圳市", "东莞市"}
			cityList, _ := GetBookCityList() //获取数据库配置的城市列表
			cityName := GetIpString(ip)      //根据IP地址反向解析所在城市
			//cityName := utils.GetCityNameByIp(ip) //根据配置获取IP请求信息
			if cityName == "" {
				log.Printf("当前IP=【%v】***************** 未解析定位到城市,可浏览盗版图书", ip)
				return 0, nil
			}
			log.Printf("根据IP =【%v】 解析到的city_name = 【%v】 查询的城市列表：【%v】\n", ip, cityName, cityList)
			index := IsContainCity(cityList, cityName) //判断返回的类型是否包含在类目中
			if index {
				status = 1
				log.Printf("当前解析的城市city_name=【%v】 满足免审核要求，可浏览正版书\n", cityName)
			} else {
				status = 0
				log.Printf("当前解析城市city_name=【%v】 不满足审核要求，可浏览盗版书\n", cityName)
			}
		} else {
			//版权标识配置未开启查询
			log.Printf("当前配置的审核状态未开启，可浏览盗版书！")
			status = 0 //非审核
		}
	} else {
		log.Printf("获取当前记录为空，默认浏览盗版书籍")
		status = 0
	}
	return status, nil
}

func GetSectionNew(req *models.SectionNewReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})

	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1").Debug()
	}
	num := req.Size
	if num <= 0 {
		num = 8
	}

	bookType := user_service.GetBookTypeByUserId(req.UserId)
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}

	db = db.Where("is_hot = 1 and is_new = 1")

	//修改随机排序替换mysql中的rand排序
	var total int64
	db.Count(&total)
	global.Requestlog.Infof("热门新书列表的总数 total = %v", total)
	offsetNum := utils.GetBookRandPosition(total)
	err = db.Offset(offsetNum).Limit(num).Find(&list).Error
	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

func GetTags(req *models.GetTagsReq) (tags []*models.TagsRank, err error) {
	// 自定义 SQL 查询
	columnType := req.ColumnType
	limit := req.Size
	//bookTable := new(models.McBook).TableName()
	//sql := fmt.Sprintf(`
	//	SELECT tag, COUNT(*) as tag_count
	//	FROM (
	//		SELECT TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', n.digit + 1), ',', -1)) as tag
	//		FROM %v
	//		JOIN (
	//			SELECT 0 as digit UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
	//		) as n
	//		ON CHAR_LENGTH(tags) - CHAR_LENGTH(REPLACE(tags, ',', '')) >= n.digit
	//		WHERE book_type = %v
	//	) as tag_split
	//	GROUP BY tag
	//	ORDER BY tag_count DESC
	//	LIMIT %v;
	//`, bookTable, bookType, limit)
	//err = global.DB.Raw(sql).Scan(&tags).Error
	bookType := user_service.GetBookTypeByUserId(req.UserId)

	var list []*models.McTag
	db := global.DB.Model(&models.McTag{}).Order("sort desc")
	db = db.Where("status = 1 and is_new = 0")
	if bookType > 0 {
		db = db.Where("book_type = ?", bookType)
	}
	if columnType > 0 {
		db = db.Where("column_type = ?", columnType)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Find(&list).Error
	if len(list) <= 0 {
		return
	}
	for _, val := range list {
		tag := &models.TagsRank{
			TagId:    val.Id,
			TagName:  val.TagName,
			TagCount: getBookCountByTagId(val.Id),
		}
		tags = append(tags, tag)
	}
	return
}

func GetTeenZoneList(req *models.TeenZoneListReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})
	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	//获取是否有版权
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan =1").Debug()
	} else {
		db = db.Where("status = 1 and is_teen = 1").Debug()
	}
	num := req.Size
	if num <= 0 {
		num = 8
	}

	teenType := req.TeenType
	var teenName string
	if teenType == 1 {
		teenName = "名著"
	} else if teenType == 2 {
		teenName = "传记"
	} else if teenType == 3 {
		teenName = "文学"
	}

	cids := GetClassIdsByName(teenName)
	if len(cids) > 0 {
		db = db.Where("cid in ?", cids)
	}

	tids := GetTagIdsByName(teenName)
	if len(tids) > 0 {
		db = db.Where("tid in ?", tids)
	}
	var (
		total    int64
		newSize  int64
		newTotal int64
	)
	//强制转换
	newSize = int64(num) //强制转换
	db.Count(&total)
	if total > newSize {
		newTotal = total - newSize
	}
	global.Requestlog.Infof("青少年的书籍总数 total = %v newTotal = %v", total, newTotal)
	offsetNum := utils.GetBookRandPosition(newTotal)
	//按照对应的进行设置配置
	err = db.Offset(offsetNum).Limit(num).Find(&list).Error

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

func GetNewBookRec(req *models.GetTagsReq) (tagBooks []*models.GetNewBookRecRes, err error) {
	var tags []*models.McTag
	bookType := user_service.GetBookTypeByUserId(req.UserId)
	if bookType <= 0 {
		bookType = 1
	}
	tags, err = getNewBookTag(bookType, req.ColumnType)
	if err != nil {
		return
	}
	global.Requestlog.Infof("共有分类的总数:%v", len(tags))
	if len(tags) <= 0 {
		return
	}
	for _, tag := range tags {
		var books []*models.SimpleBookRes
		//获取书籍的ID信息
		books = GetSimpleBooksByTag(tag.Id, 1, 8)
		if len(books) > 0 {
			global.Requestlog.Infof("当前分类标签ID= %v Tag_name= %v 共查到的书籍总数：%v", tag.Id, tag.TagName, len(books))
			book := &models.GetNewBookRecRes{
				TagId:    tag.Id,
				TagName:  tag.TagName,
				BookList: books,
			}
			tagBooks = append(tagBooks, book)
		} else {
			//如果没有书籍就提示下
			global.Requestlog.Infof("当前分类标签ID= %v Tag_name= %v,并无分类的书籍关联信息", tag.Id, tag.TagName)
		}
	}
	fmt.Println("共有分类的数据", len(tagBooks))
	return
}

// 获取新书的推荐
func GetNewBookList(c *gin.Context, req *models.GetNewBookListReq) (list []*models.McBook, total int64, err error) {
	tid := req.Tid
	if tid <= 0 {
		err = fmt.Errorf("%v", "新书标签不能为空")
		return
	}
	page := req.Page
	size := req.Size
	ip := req.Ip                    //ip地址
	device_type := req.DeviceType   //客户端类型
	package_name := req.PackageName //包名
	mark := req.Mark                //渠道号

	//端号如果获取不到从header=os取
	if device_type == "" {
		device_type = utils.GetRequestHeaderByName(c, "Os")
		log.Printf("端号为空，接下来会从header【Os =%v】 重新再获取一次", device_type)
	}
	//包名如果获取不到从header=Package取
	if package_name == "" {
		package_name = utils.GetRequestHeaderByName(c, "Package")
		log.Printf("包名为空，接下来会从header 【Package =%v】  重新再获取一次", package_name)
	}
	//根据标签获取新书列表
	list, total, err = GetBooksByTag(tid, page, size, ip, device_type, package_name, mark)
	return
}

/*
* @note 通过条件查询今日的新书
* @param condition string 搜索条件
8 @param pageNum int 页码
* @param pageSize  int最大显示数量
* @return object ,total , err
*/
func GetBookListByCondition(condition string, pageNum int, pageSize int) (list []*models.McBook, total int64, err error) {
	if condition == "" {
		return
	}
	db := global.DB.Model(&models.McBook{}).Order("id desc").Debug()
	db = db.Where(condition)
	db.Count(&total)
	//页码特殊处理,判断默认值
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 8
	}
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Limit(pageSize).Find(&list).Error
	}
	return

}

func TodayUpdateBooks(req *models.TodayUpdateBooksReq) (list []*models.McBook, total int64, err error) {
	//db := global.DB.Model(&models.McBook{}).Order("id desc")
	//.Order("uptime desc")
	//今日推荐排序：优先显示今日入库的，如果没有按照推荐进行排序 ,
	//搜索今日推荐的如果没有，再按照今日更新继续搜索一遍
	//如果今日更新都没有，就查全部的推荐的书
	fmt.Println("333324234", req.Ip, req.DeviceType, req.PackageName)
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	var bookCondition string
	if bookStatus == 1 {
		bookCondition = "status = 1 and is_banquan =1"
	} else {
		bookCondition = "status = 1"
	}
	//查今天的
	todayCondition := fmt.Sprintf("(%v and FROM_UNIXTIME(addtime) BETWEEN '%v' and '%v') or  (%v and is_rec=1) ", bookCondition, utils.GetDate(), utils.GetTomorrowDate(), bookCondition)
	//获取今日的的入库书籍，
	list, total, err = GetBookListByCondition(todayCondition, req.Page, req.Size)
	if total == 0 {
		////如果为空，查所有的status=1 and is_rec=1的书，保证列表数据不会为空
		//allBookCondition := fmt.Sprintf("%v and is_rec=1 ", bookCondition)
		//log.Printf("今日数据为空，查询线上所有状态开启的推荐书籍 condition 【 %v】\n", allBookCondition)
		//list, total, err = GetBookListByCondition(allBookCondition, req.Page, req.Size)
	}
	//特殊判断
	if err != nil {
		return nil, 0, nil
	}
	//db = db.Where(todayCondition).Debug()
	//db.Count(&total)
	////如果今日的书籍为空，会继续查一次的所有推荐的书籍
	//if total == 0 {
	//	fmt.Println("今日更新书籍为空，重新请求推荐的数据")
	//	allRecCondition := fmt.Sprintf("%v and is_rec=1", bookCondition)
	//	db = db.Where(allRecCondition)
	//	db.Count(&totalAll)
	//	total = totalAll
	//}
	//
	//pageNum := req.Page
	//pageSize := req.Size
	//
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	//if pageSize == 0 {
	//	pageSize = 8
	//}
	//
	////if pageNum > 0 && pageSize > 0 {
	////	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	////} else {
	////	err = db.Limit(pageSize).Find(&list).Error
	////}
	//
	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

func GetHotCount(req *models.HotBookCountReq) (endCount, rankCount, newCount int64) {
	ip := req.Ip
	deviceType := req.DeviceType
	packageName := req.PackageName
	mark := req.Mark
	fmt.Println(ip, deviceType, packageName)
	endCount = GetEndCountByHot(deviceType, packageName, ip, mark)
	rankCount = GetRankCountByHot(deviceType, packageName, ip, mark)
	newCount = GetNewCountByHot(deviceType, packageName, ip, mark)
	return
}

func getBookCountByTagId(tagId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McBook{}).Where("tid = ?", tagId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
