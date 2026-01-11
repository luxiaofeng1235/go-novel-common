package collect_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/collect_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"regexp"
	"strings"
)

func CollectThread(collectId int64) (err error) {
	//采集分页列表
	var collectListData *models.CollectListData
	var noReturn bool
	collectListData, noReturn, err = CollectListGet(collectId)
	if noReturn == false {
		//collect_rm
		CollectRm(collectId)
		//返回采集完成 'state'=>'finish'
		global.Collectlog.Infof("采集完成")
		return
	}
	if collectListData == nil {
		global.Collectlog.Errorf("获取采集列表分页失败")
		return
	}
	collectBookUrl := collectListData.Url
	//CollectListRm
	err = CollectListRm(collectId, collectBookUrl)
	if err != nil {
		return
	}
	//FieldContent
	collectInfo := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectInfo, collectId))
	collect := models.McCollect{}
	err = json.Unmarshal([]byte(collectInfo), &collect)
	if err != nil {
		global.Collectlog.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
		return
	}
	var collecBookInfoRes *models.CollecBookInfoRes
	collecBookInfoRes, err = FieldContent(&collect, collectBookUrl)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}
	if collecBookInfoRes.BookUrl != "" {
		log.Println(collecBookInfoRes.BookUrl, collecBookInfoRes.BookName)
		SaveData(&collect, collecBookInfoRes)
	}
	//log.Printf("%+v", utils.JSONString(collectInfo))
	//log.Printf("%+v", utils.JSONString(collecBookInfoRes))
	//for _, chapter := range collecBookInfoRes.Chapters {
	//	log.Println(chapter.ChapterLink, chapter.ChapterTitle)
	//}

	return
}

func CollectListGet(collectId int64) (collectListData *models.CollectListData, noReturn bool, err error) {
	noReturn = true
	list := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectList, collectId))
	collectRes := new(models.CollectRes)
	err = json.Unmarshal([]byte(list), &collectRes)
	if err != nil {
		noReturn = false
		global.Collectlog.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
		return
	}
	collectResData := collectRes.Data
	for _, val := range collectResData {
		if val.Lock == 0 {
			val.Lock = 1
			err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectList, collectId), collectRes, 0)
			if err != nil {
				global.Collectlog.Errorf("缓存采集当前分页状态失败 %v", err.Error())
				return
			}
			collectListData = &models.CollectListData{
				Url:     val.Url,
				Lock:    val.Lock,
				Count:   collectRes.Count,
				PageNum: collectRes.PageNum,
			}
			return
		}
	}

	if collectRes.PageNum < collectRes.Count {
		nextCurrent := collectRes.PageNum + 1
		collectRes, err = CollectListSet(collectId, nextCurrent)
		err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectLog, collectId), nextCurrent, 0)
		if err != nil {
			global.Collectlog.Errorf("缓存采集当前页数失败 %v", err.Error())
			return
		}
		collectListData, noReturn, err = CollectListGet(collectId)
		return
	}
	noReturn = false
	return
}

func CollectListRm(collectId int64, url string) (err error) {
	queueList := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectList, collectId))
	collectRes := new(models.CollectRes)
	err = json.Unmarshal([]byte(queueList), &collectRes)
	if err != nil {
		global.Collectlog.Errorf("解析采集信息出错 collectId=%v err=%v", collectId, err.Error())
		return
	}
	bookUrls := removeURL(collectRes.Data, url)
	collectRes.Data = bookUrls
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectList, collectId), collectRes, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集分页列表出错 %v", err.Error())
		return
	}
	return
}

func removeURL(urlList []*models.CollectDataRes, targetURL string) (newList []*models.CollectDataRes) {
	for _, item := range urlList {
		if item.Url != targetURL {
			newList = append(newList, item)
		}
	}
	return
}

func CollectRm(collectId int64) {
	var err error
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectInfo, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存采集信息失败 %v", err.Error())
		return
	}
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectSourceUrl, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存采集分页列表失败 %v", err.Error())
		return
	}
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectList, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存采集当前分页状态失败 %v", err.Error())
		return
	}
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectLog, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存采集当前页数失败 %v", err.Error())
		return
	}
	return
}

func CollectChapterThread() (err error) {
	var keys []string
	keys, _, err = redis_service.LoadRedisKeysNumber(utils.CollectChapter)
	if err != nil {
		global.Collectlog.Errorf("获取缓存章节Key出错 %v", err.Error())
		return
	}
	if len(keys) <= 0 {
		return
	}
	for _, key := range keys {
		chapterContent := redis_service.Get(key)
		if err != nil {
			global.Collectlog.Errorf("获取缓存章节列表出错 %v", err.Error())
			continue
		}
		if chapterContent == "" {
			global.Collectlog.Errorf("%v", "获取缓存章节列表出错")
			continue
		}

		var bookChapters []*models.CollectChapterInfo
		err = json.Unmarshal([]byte(chapterContent), &bookChapters)
		if err != nil {
			global.Collectlog.Errorf("解析采集信息失败  err=%v", err.Error())
			return
		}
		var chapterTextReg, charset, chapterTable string
		var collectId, bookId int64
		parts := strings.Split(key, "_")
		if len(parts) >= 4 {
			bookId = utils.FormatInt64(parts[2])
			//章节表
			chapterTable, err = book_service.GetChapterTable(bookId)
			if err != nil {
				global.Collectlog.Errorf("%v", "生成章节表失败")
				continue
			}
			log.Println(chapterTable)
			collect, _ := collect_service.GetCollectById(utils.FormatInt64(parts[3]))
			collectId = collect.Id
			chapterTextReg = collect.ChapterTextReg
			charset = collect.Charset
		}
		for _, chapter := range bookChapters {
			var chapterHtml, text string
			chapterHtml, err = utils.GetHtml(chapter.ChapterLink, charset, 0, 1)
			if err != nil {
				continue
			}
			//chapterId := book_service.GetChapterIdByName(chapterTable, chapter.ChapterTitle)
			text, err = collect_service.GetChapterContent(chapterTextReg, chapterHtml)
			if err != nil {
				continue
			}
			var book *models.McBook
			book, err = book_service.GetBookById(bookId)
			if err != nil {
				return
			}
			_, text, err = book_service.GetBookTxt(book.BookName, book.Author, "", text)
			if err != nil {
				continue
			}
			bookChapters = removeChapter(bookChapters, chapter.ChapterLink)
			err = redis_service.Set(fmt.Sprintf("%v_%v_%v", utils.CollectChapter, bookId, collectId), bookChapters, 0)
			if err != nil {
				global.Collectlog.Errorf("缓存采集分页列表出错 %v", err.Error())
				return
			}
			log.Println(bookId, chapter.ChapterLink, chapter.ChapterTitle, "写入成功")
		}
	}
	return
}

// 重新采集 current=0 继续采集current=page
func Collect(collectId int64, pageNum int) (collectRes *models.CollectRes, err error) {
	var collect *models.McCollect
	collect, err = collect_service.GetCollectById(collectId)
	if err != nil {
		return
	}
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectInfo, collectId), collect, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集信息失败 %v", err.Error())
		return
	}
	listPageUrls, err := utils.GetListSourceURL(collect.ListPageReg)
	if err != nil {
		return
	}
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectSourceUrl, collectId), listPageUrls, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集分页列表失败 %v", err.Error())
		return
	}
	collectRes, err = CollectListSet(collectId, pageNum)
	return
}

// 单本采集
func Alone(collectId int64, collectBookUrl string) (err error) {
	var collect *models.McCollect
	collect, err = collect_service.GetCollectById(collectId)
	if err != nil {
		return
	}
	var collecBookInfoRes *models.CollecBookInfoRes
	collecBookInfoRes, err = FieldContent(collect, collectBookUrl)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}
	if collecBookInfoRes.BookUrl != "" {
		SaveData(collect, collecBookInfoRes)
	}
	//['title'=>$return['title'],'url'=>$return['reurl'],'status'=>$return['status']]
	return
}

// 询问重新采集还是继续采集
func CollectContinuation(collectId int64) (err error) {
	nextCurrent := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectLog, collectId))
	var currentPage int
	if nextCurrent != "" {
		err = json.Unmarshal([]byte(nextCurrent), &currentPage)
		if err != nil {
			global.Collectlog.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
			return
		}
	}
	if currentPage > 0 {
		global.Collectlog.Infof("继续采集 page=%v", currentPage)
	} else {
		global.Collectlog.Infof("重新采集 page=%v", collectId)
	}
	return
}

func CollectListSet(collectId int64, pageNum int) (collectRes *models.CollectRes, err error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	collectInfo := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectInfo, collectId))
	collect := models.McCollect{}
	err = json.Unmarshal([]byte(collectInfo), &collect)
	if err != nil {
		global.Collectlog.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
		return
	}
	charset := collect.Charset
	urlComplete := collect.UrlComplete
	listSection := collect.ListSectionReg
	listUrl := collect.ListUrlReg
	urlReverse := collect.UrlReverse

	listPageVal := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectSourceUrl, collectId))
	var listPageUrls []string
	err = json.Unmarshal([]byte(listPageVal), &listPageUrls)
	if err != nil {
		global.Collectlog.Errorf("解析采集分页链接失败 collectId=%v err=%v", collectId, err.Error())
		return
	}
	index := pageNum - 1
	//获取列表失败 'state'=>'stop'
	if len(listPageUrls) <= index {
		global.Collectlog.Errorf("解析采集分页链接失败 listPageUrls=%v sliceIndex=%v", listPageUrls, index)
		return
	}

	collectRes = new(models.CollectRes)

	collectRes.PageNum = pageNum
	collectRes.Count = len(listPageUrls)
	sourceUrl := listPageUrls[index]
	sleepSecond := getSleepSecond()
	listContentHtml, err := utils.GetHtml(sourceUrl, charset, urlComplete, sleepSecond)

	if err != nil {
		global.Collectlog.Errorf("解析采集分页html出错 sourceUrl=%v err=%v", sourceUrl, err.Error())
		return
	}
	if listContentHtml == "" {
		//未获取到起始页面数据! 'state'=>'stop'
		global.Collectlog.Errorf("未获取到起始页面数据! sourceUrl=%v err=%v", sourceUrl, err.Error())
		return
	}
	match := regexp.MustCompile(listSection).FindStringSubmatch(listContentHtml)
	if len(match) <= 0 {
		global.Collectlog.Errorf("采集获取列表失败 err=%v", err.Error())
		return
	}
	linkPattern := listUrl
	linksRe := regexp.MustCompile(linkPattern)
	links := linksRe.FindAllStringSubmatch(match[0], -1)
	var bookLinks []string
	for _, link := range links {
		if len(link) > 1 {
			href := strings.TrimSpace(link[1])
			bookLinks = append(bookLinks, href)
		}
	}

	if urlReverse > 0 {
		bookLinks = utils.ArrayReverse(bookLinks)
	}
	if len(bookLinks) <= 0 {
		return
	}

	var collectDataRes []*models.CollectDataRes
	for _, link := range bookLinks {
		collectData := &models.CollectDataRes{
			Url:  link,
			Lock: 0,
		}
		collectDataRes = append(collectDataRes, collectData)
	}
	collectRes.Data = collectDataRes
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectList, collectId), collectRes, 0)
	return
}
