package book_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/common_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
	"log"
	"strings"
)

func GetBookById(id int64) (book *models.McBook, err error) {
	err = global.DB.Model(models.McBook{}).Where("id = ?", id).Find(&book).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetSimpleBooksByTag(tagId int64, isNew, limit int) (books []*models.SimpleBookRes) {
	var err error
	db := global.DB.Model(models.McBook{}).Debug()

	//处理排序推荐问题
	db = db.Order("id desc")
	db = db.Where("status = 1")
	if tagId > 0 {
		db = db.Where("tid = ?", tagId)
	}
	if isNew > 0 {
		db = db.Where("is_new = 1")
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Find(&books).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if len(books) > 0 {
		for _, book := range books {
			book.Pic = utils.GetFileUrl(book.Pic)
		}
	}
	return
}

func GetBooksByTag(tid int64, pageNum, pageSize int, ip, device_type, package_name, mark string) (list []*models.McBook, total int64, err error) {
	db := global.DB.Model(models.McBook{})

	//获取当前是否为版权书
	bookStatus, _ := GetBookCopyright(device_type, package_name, ip, mark)
	if bookStatus == 1 { //判断版权
		db = db.Where("status = 1 and is_banquan=1")
	} else {
		db = db.Where("status = 1 ")
	}
	db = db.Where("is_new = 1").Debug()
	//优先展示推荐最新的靠前
	sortData, _ := GetApiBookRecByType("rec_new")
	if sortData != "" { //查推荐的新书
		db = db.Order(sortData)
	} else {
		db = db.Order("id desc")
	}

	//根据分类标签进行查询
	if tid > 0 {
		db = db.Where("tid = ?", tid)

	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if len(list) > 0 {
		for _, book := range list {
			book.Pic = utils.GetFileUrl(book.Pic)
		}
	}
	return
}

// 获取迅搜的全文检索的关键字数据信息
func GetSearchBookIds(title string, page int64, limit int64) (myString []int64, err error) {
	if title == "" {
		return
	}
	//搜索接口中的新书
	jsonData := common_service.XunSearchByBookName("searchList", title)
	var resp *common_service.XunResponse
	err = json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var slice []int64
	if resp.Data != nil {
		for _, val := range resp.Data {
			slice = append(slice, val.BookID) //追加关联数据信息
		}
	}
	log.Printf("获取全文检索后的书籍ID = %v\n", slice)
	return slice, nil
}

// 根据分类ID统计列表数据
func BookCateListRes(req *models.ApiCateBookReq) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{}).Debug()
	bookId := req.Bid
	if bookId == 0 {
		return
	}
	//获取小说的详情信息
	bookInfo, err := book_service.GetBookById(bookId)
	if bookInfo.Id <= 0 {
		err = fmt.Errorf("%v", "小说不存在")
		return
	}
	cid := bookInfo.Cid //获取分类ID
	global.Requestlog.Infof("获取book_id = %v ,对应的 cid = %v", bookId, cid)
	//判断cid不能为空
	if cid == 0 {
		cid = 1 //默认取1如果没有获取到的话
	}
	limit := 3  //线束数量
	status := 1 //正常显示状态
	//获取版全面的配置判断
	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	//处理版全面的显示，否则显示正常的分类下的数据
	if bookStatus == 1 {
		db = db.Where("status = ? and cid = ? and is_banquan= 1", status, cid)
	} else {
		db = db.Where("status = ? and cid = ?", status, cid)
	}

	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	global.Requestlog.Infof("当前cid =%v 对应的总数 total = %v", cid, total)
	//获取随机的数量
	offsetNum := utils.GetBookRandPosition(total)
	err = db.Offset(offsetNum).Limit(limit).Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}

	if len(list) <= 0 {
		return
	}
	//自动获取对应的图片信息
	for _, val := range list {
		val.Pic = utils.GetFileUrl(val.Pic)
	}
	return
}

// 书籍搜索
func BookListSearch(req *models.ApiBookListReq) (list []*models.McBook, total int64, pageNum, pageSize int, err error) {
	db := global.DB.Model(&models.McBook{}).Debug()

	bookStatus, err := GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	//判断是否为版权的书
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan=1")
	} else {
		db = db.Where("status = 1")
	}

	cid := req.Cid
	if cid > 0 {
		db = db.Where("cid = ?", cid)
	}

	serialize := req.Serialize
	if serialize > 0 {
		db = db.Where("serialize = ?", serialize)
	}

	sort := req.Sort
	if sort == utils.Hot {
		db = db.Where("is_hot = 1")
	} else if sort == utils.Hits {
		db = db.Order("hits desc")
	} else if sort == utils.Score {
		db = db.Order("score desc")
	} else if sort == utils.New {
		db = db.Where("is_new = 1") //查询新书
		db = db.Order("addtime desc")
	} else if sort == utils.Classic {
		db = db.Where("is_classic = 1")
	} else if serialize == 2 { //推荐完结的排序
		//推荐完结
		sortStr, _ := GetApiBookRecByType("rec_serialize")
		if sortStr != "" {
			db = db.Order(sortStr)
		}
	}

	isPay := req.IsPay
	if isPay > 0 {
		db = db.Where("is_pay = ?", isPay)
	}
	textNumType := req.TextNumType
	if textNumType == 3 {
		db = db.Where("text_num > ?", 2000000)
	} else if textNumType == 2 {
		db = db.Where("text_num > ? and text_num < ?", 500000, 1000000)
	} else if textNumType == 1 {
		db = db.Where("text_num < ?", 500000)
	}

	bookName := strings.TrimSpace(req.BookName)
	userId := req.UserId
	if bookName != "" {
		go SearchLog(bookName, userId)
		//搜索全文索引内容信息
		var bookIds []int64
		//根据全文检索来进行搜索
		bookIds, err = GetSearchBookIds(bookName, int64(req.Page), int64(req.Size))
		if bookIds != nil {
			if len(bookIds) > 1 {
				db = db.Where("id in(?)", bookIds)
			} else { //只有一个的时候直接用=
				db = db.Where("id = ?", bookIds[0])
			}
			/////排序流程处理
			//转换成字符串进行排序
			bookStr := utils.JoinInt64ToString(bookIds)
			orderStr := fmt.Sprintf("field(id,%v)", bookStr)
			db = db.Order(orderStr)
		} else {
			db = db.Where("id in (-1)")
		}
		//先使用回表索引进行查询优化，后期用搜索引擎如：gofound来进行替换分词进行搜索
		//采用Mysql5.7以后支持全文检索来进行搜索，替代原有的like查询效率提升3倍
		//searchBookString := fmt.Sprintf("MATCH (book_name) AGAINST ('\"%s\"' IN BOOLEAN MODE)", bookName)
		////db = db.Where("book_name LIKE ?", "%"+bookName+"%")
		//db = db.Where(searchBookString)
	} else {
		db = db.Order("id desc")
	}

	tid := req.Tid
	if tid > 0 {
		db = db.Where("tid = ?", tid)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error

	if err != nil {
		return
	}
	pageNum = req.Page
	pageSize = req.Size

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	//搜索结果查询
	if pageNum > 0 && pageSize > 0 {
		err = db.Debug().Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Debug().Find(&list).Error
	}

	if len(list) > 0 {
		for _, val := range list {
			if val.ReadCount <= 0 {
				val.ReadCount = GetReadCountByBookId(val.Id)
			}
			if val.SearchCount <= 0 {
				val.SearchCount = getSearchCountByBookName(val.BookName)
			}
			log.Printf("数据库查出来的pic = %v \n", val.Pic)
			val.Pic = utils.GetFileUrl(val.Pic)
		}
	}
	return
}

func SearchLog(searchName string, userId int64) {
	var err error
	var count int64
	//更新搜索mc_book表里的search_count字段信息
	_ = book_service.UpdateSearchNumByName(searchName)
	db := global.DB.Model(models.McBookSearch{}).Where("search_name = ? and uid = ?", searchName, userId)
	db.Count(&count)
	if count > 0 {
		da := make(map[string]interface{})
		da["num"] = gorm.Expr("num + ?", 1)
		da["uptime"] = utils.GetUnix()
		err = db.Updates(da).Error
		if err != nil {
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	} else {
		search := models.McBookSearch{
			Uid:        userId,
			SearchName: searchName,
			Num:        1,
			Addtime:    utils.GetUnix(),
			Uptime:     utils.GetUnix(),
		}
		if err = global.DB.Create(&search).Error; err != nil {
			global.Sqllog.Errorf("写入搜索记录失败 err=%v", err.Error())
			return
		}
	}
	return
}

func GetEndCountByHot(device_type, package_name, ip, mark string) (count int64) {
	var err error

	fmt.Println("11111111111111111111", ip)
	db := global.DB.Model(models.McBook{})

	bookStatus, err := GetBookCopyright(device_type, package_name, ip, mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status =1 and is_banquan = 1 and is_hot = 1 and serialize = 2").Debug()
	} else {
		db = db.Where("status =1 and is_hot = 1 and serialize = 2").Debug()
	}
	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetRankCountByHot(device_type, package_name, ip, mark string) (count int64) {
	var err error
	fmt.Println("11111111111111111111", ip)
	db := global.DB.Model(models.McBook{})

	bookStatus, err := GetBookCopyright(device_type, package_name, ip, mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status =1 and is_banquan = 1 and is_hot = 1 and is_rec = 1").Debug()
	} else {
		db = db.Where("status =1 and is_hot = 1 and is_rec = 1").Debug()
	}

	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetNewCountByHot(device_type, package_name, ip, mark string) (count int64) {
	var err error
	fmt.Println("11111111111111111111", ip)
	db := global.DB.Model(models.McBook{})

	bookStatus, err := GetBookCopyright(device_type, package_name, ip, mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status =1 and is_banquan = 1 and is_hot = 1 and is_new = 1").Debug()
	} else {
		db = db.Where("status =1 and is_hot = 1 and is_new = 1").Debug()
	}
	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
