package search_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/book_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"reflect"
)

func SearchHistory(req *models.SearchHistoryReq) (searchs []*models.McBookSearch, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	searchs, err = GetBookSearchList()
	if err != nil {
		return
	}
	return
}

// 获取排行的榜单数据信息的TOP10榜单数据-定时计算
func GetRankSearchTop10(limit int, rec_type string) (ids []int64, err error) {
	var list []*models.McBookSearchRank
	db := global.DB.Model(models.McBookSearchRank{}).Select("bid").Order("id asc").Debug()
	db = db.Where("rec_type = ?", rec_type)
	db = db.Limit(limit) //默认取TOP10
	err = db.Find(&list).Error
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
	log.Printf("获取热搜榜单TOP10的排行数据为: %v\n", myIds)
	return myIds, nil
}

func GetHotSearchRank(req *models.SearchHotReq) (list []*models.McBook, err error) {
	num := 10
	db := global.DB.Model(&models.McBook{}).Debug()
	//获取排行的热热门榜单数据
	searchBookIds, err := GetRankSearchTop10(num, "hot_search")
	if err != nil {
		global.Sqllog.Errorf("GetHotSearchRank err:%v", err)
	}
	//如果存在TOP10榜单数据的话
	if len(searchBookIds) > 0 {
		db = db.Where("id in (?)", searchBookIds)
		bookStr := utils.JoinInt64ToString(searchBookIds)
		orderStr := fmt.Sprintf("field(id,%v)", bookStr)
		db = db.Order(orderStr) //按照此进行排序设置
	} else {
		db = db.Order("search_count desc")
	}
	log.Printf("rank list:%v type:%v\n", searchBookIds, reflect.TypeOf(searchBookIds))
	fmt.Println("11111111111111111111", req.Ip)
	bookStatus, err := book_service.GetBookCopyright(req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1")
	} else {
		db = db.Where("status = 1") //查所有
	}

	err = db.Limit(num).Find(&list).Error
	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		if val.Pic != "" {
			val.SearchNum = val.SearchCount //兼容search_num字段
			val.Pic = utils.GetFileUrl(val.Pic)
		}
	}
	return
}
