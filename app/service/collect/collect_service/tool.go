package collect_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/redis_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"math"
	"regexp"
	"strings"
)

func getSleepSecond() (sleepSecond int64) {
	var collectSleep string
	var err error
	collectSleep, err = setting_service.GetValueByName("collectSleep")
	if err != nil {
		global.Collectlog.Errorf("获取采集间隔时间失败 %v", err.Error())
		return
	}
	sleepSecond = utils.FormatInt64(collectSleep)
	return
}

func removeChapter(chapterList []*models.CollectChapterInfo, chapterLink string) (newList []*models.CollectChapterInfo) {
	for _, item := range chapterList {
		if item.ChapterLink != chapterLink {
			newList = append(newList, item)
		}
	}
	return
}

func GetPageUrls(collect *models.McCollect) (pageUrls []string, err error) {
	patternListSection := collect.ListSectionReg
	if patternListSection == "" {
		err = fmt.Errorf("列表区间正则不能为空 collectId=%v", collect.Id)
		return
	}

	patternListUrl := collect.ListUrlReg
	if patternListUrl == "" {
		err = fmt.Errorf("获取小说链接地址正则为空 collectId=%v", collect.Id)
		return
	}
	listPageReg := collect.ListPageReg
	if listPageReg == "" {
		err = fmt.Errorf("获取小说分页列表正则为空 collectId=%v", collect.Id)
		return
	}
	pageUrls, err = utils.GetListSourceURL(collect.ListPageReg)
	return
}

func GetPageBookUrls(collect *models.McCollect, pageUrl string) (bookUrls []string, err error) {
	charset := collect.Charset
	urlComplete := collect.UrlComplete
	listSection := collect.ListSectionReg
	listUrl := collect.ListUrlReg
	urlReverse := collect.UrlReverse
	sleepSecond := getSleepSecond()
	listContentHtml, err := utils.GetHtml(pageUrl, charset, urlComplete, sleepSecond)

	if err != nil {
		err = fmt.Errorf("解析采集分页html出错 pageUrl=%v err=%v", pageUrl, err.Error())
		return
	}
	if listContentHtml == "" {
		//未获取到起始页面数据! 'state'=>'stop'
		err = fmt.Errorf("未获取到起始页面数据 pageUrl=%v err=%v", pageUrl, err.Error())
		return
	}

	match := regexp.MustCompile(listSection).FindStringSubmatch(listContentHtml)
	if len(match) <= 0 {
		err = fmt.Errorf("采集获取列表失败 pageUrl=%v err=%v", pageUrl, err.Error())
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
	bookUrls = bookLinks
	return
}

func RmmoveCollectPageBooks(collectId int64, bookUrl string) (err error) {
	queueList := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId))
	var pageBook *models.CollectPageBook
	var pageBookUrls []*models.CollectBookUrl
	err = json.Unmarshal([]byte(queueList), &pageBook)
	if err != nil {
		global.Collectlog.Errorf("解析采集小说列表出错 collectId=%v err=%v", collectId, err.Error())
		return
	}
	pageBookUrls = pageBook.BookUrls
	bookUrls := removeBookUrl(pageBookUrls, bookUrl)
	pageBook.BookUrls = bookUrls
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId), pageBook, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集小说列表出错 %v", err.Error())
		return
	}
	return
}

func GetCollectProgress(pageBook *models.CollectPageBook) (progress float64) {
	pageCount := pageBook.PageCount
	pageNum := pageBook.PageNum
	bookCount := pageBook.BookCount
	bookNum := len(pageBook.BookUrls)

	bookReadNum := bookCount - bookNum
	totalPages := pageCount * bookCount
	readPages := (pageNum-1)*bookCount + bookReadNum

	progress = float64(readPages) / float64(totalPages) * 100
	progress = math.Round(progress*100) / 100
	return
}

func removeBookUrl(urls []*models.CollectBookUrl, targetURL string) (newList []*models.CollectBookUrl) {
	for _, info := range urls {
		if info.Url != targetURL {
			newList = append(newList, info)
		}
	}
	return
}

func removeChapterUrl(urls []*models.CollectChapterInfo, targetURL string) (newList []*models.CollectChapterInfo) {
	for _, info := range urls {
		if info.ChapterLink != targetURL {
			newList = append(newList, info)
		}
	}
	return
}

func resetPage(collect *models.McCollect) (err error) {
	var collectPageUrl = new(models.CollectPageUrlRes)
	collectId := collect.Id
	var pageUrls []string
	pageUrls, err = GetPageUrls(collect)
	if err != nil {
		return
	}
	collectPageUrl.PageNum = 1
	collectPageUrl.PageCount = len(pageUrls)
	collectPageUrl.PageUrls = pageUrls
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageUrl, collectId), collectPageUrl, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集当前页数失败 %v", err.Error())
		return
	}
	return
}

func RemoveCollect(collectId int64) {
	var err error
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectPageUrl, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存采集分页列表失败 %v", err.Error())
		return
	}
	err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId))
	if err != nil {
		global.Collectlog.Errorf("删除缓存小说列表失败 %v", err.Error())
		return
	}
	return
}

func GetSourceIdsByBookId(bookId int64) (sourceIds []int64) {
	var err error
	err = global.DB.Model(models.McBookSource{}).Where("bid = ?", bookId).Pluck("sid", &sourceIds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetSourceCountById(bookId, sourceId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McBookSource{}).Where("bid = ? and sid = ?", bookId, sourceId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func BookUrlUnLock(collectId int64) (err error) {
	list := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId))
	var pageBook *models.CollectPageBook
	if list != "" {
		err = json.Unmarshal([]byte(list), &pageBook)
		if err != nil {
			err = fmt.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
			return
		}
		bookUrls := pageBook.BookUrls
		for _, val := range bookUrls {
			if val.Lock == 1 {
				val.Lock = 0
			}
		}
		err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId), pageBook, 0)
		if err != nil {
			err = fmt.Errorf("缓存采集当前分页小说列表失败 %v", err.Error())
			return
		}
	}
	return
}
