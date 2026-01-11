package collect_service

import (
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/collect_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"regexp"
	"strings"
	"sync"
)

func StartCollect(collectId int64, isRestart bool) (err error) {
	if collectId <= 0 {
		err = fmt.Errorf("%v", "采集id不能为空")
		return
	}
	collect, err := collect_service.GetCollectById(collectId)
	if err != nil {
		return
	}
	collectId = collect.Id
	if collectId <= 0 {
		err = fmt.Errorf("采集id不存在 collectId=%v", collectId)
		return
	}
	if collect.Status != 1 {
		err = fmt.Errorf("collectId=%v 当前采集已停用", collectId)
		return
	}
	err = GetCollectPageList(collect, isRestart)
	return
}

func GetCollectPageList(collect *models.McCollect, isRestart bool) (err error) {
	collectId := collect.Id
	//获取列表失败 'state'=>'stop'
	bookUrl, progress, isEnd, err := GetCollectPageBookUrl(collect, isRestart)
	if err != nil {
		return
	}
	var collecBookInfoRes *models.CollecBookInfoRes
	collecBookInfoRes, err = FieldContent(collect, bookUrl)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}
	if collecBookInfoRes.BookUrl != "" {
		err = SaveData(collect, collecBookInfoRes)
		if err != nil {
			global.Collectlog.Errorf("%v", err.Error())
			return
		}
	}

	err = RmmoveCollectPageBooks(collectId, bookUrl)
	if err != nil {
		return
	}
	if progress > 0 {
		global.Collectlog.Infof("%v%v", "%", progress)
	}
	if isEnd {
		global.Collectlog.Infof("%v", "采集完成")
		RemoveCollect(collectId)
		return
	}
	return
}

// 获取采集网站分页列表
func GetCollectPageBookUrl(collect *models.McCollect, isRestart bool) (bookUrl string, progress float64, isEnd bool, err error) {
	collectId := collect.Id
	list := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId))
	var bookUrls []*models.CollectBookUrl
	var pageBook *models.CollectPageBook
	var pageNum, pageCount int
	if list != "" {
		err = json.Unmarshal([]byte(list), &pageBook)
		if err != nil {
			global.Collectlog.Errorf("解析采集信息失败 collectId=%v err=%v", collectId, err.Error())
			return
		}
		progress = GetCollectProgress(pageBook)
		pageNum = pageBook.PageNum
		pageCount = pageBook.PageCount
		bookUrls = pageBook.BookUrls
		if pageNum >= pageCount && len(bookUrls) <= 1 {
			isEnd = true
		}
		for _, val := range bookUrls {
			if val.Lock == 1 {
				bookUrl = val.Url
				return
			}
			if val.Lock == 0 {
				val.Lock = 1
				err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId), pageBook, 0)
				if err != nil {
					global.Collectlog.Errorf("缓存采集当前分页小说列表失败 %v", err.Error())
					return
				}
				bookUrl = val.Url
				return
			}
		}
	}
	if len(bookUrls) <= 0 {
		bookUrl, pageNum, pageCount, err = GetCollectBookUrl(collect, isRestart)
		if err != nil {
			return
		}
	}
	return
}

// 获取当前分页链接
func GetCollectPageUrl(collect *models.McCollect, isRestart bool) (pageUrl string, pageNum, pageCount int, err error) {
	collectId := collect.Id
	pageVal := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectPageUrl, collectId))
	var collectPageUrl = new(models.CollectPageUrlRes)
	if pageVal != "" {
		err = json.Unmarshal([]byte(pageVal), &collectPageUrl)
		if err != nil {
			err = fmt.Errorf("解析采集分页链接失败 collectId=%v err=%v", collectId, err.Error())
			return
		}
	}

	if collectPageUrl.PageNum <= 0 {
		collectPageUrl.PageNum = 1
	}

	pageUrls := collectPageUrl.PageUrls
	pageNum = collectPageUrl.PageNum
	pageCount = collectPageUrl.PageCount
	index := collectPageUrl.PageNum - 1
	if len(pageUrls) <= 0 {
		err = resetPage(collect)
		if err != nil {
			return
		}
		pageVal = redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectPageUrl, collectId))
		err = json.Unmarshal([]byte(pageVal), &collectPageUrl)
		if err != nil {
			err = fmt.Errorf("解析采集分页链接失败 collectId=%v err=%v", collectId, err.Error())
			return
		}
		pageUrls = collectPageUrl.PageUrls
		pageCount = collectPageUrl.PageCount
		index = collectPageUrl.PageNum - 1
	}
	if len(pageUrls) <= index {
		err = fmt.Errorf("已经采集完毕 pageUrls=%v index=%v", pageUrls, index)
		err = resetPage(collect)
		if err != nil {
			return
		}
		return
	}
	pageCount = collectPageUrl.PageCount
	pageUrl = collectPageUrl.PageUrls[index]
	if collectPageUrl.PageNum <= collectPageUrl.PageCount {
		nextPage := collectPageUrl.PageNum + 1
		if isRestart == false {
			if nextPage > collectPageUrl.PageCount {
				err = fmt.Errorf("%v", "已采集完毕")
				return
			}
		}
		collectPageUrl.PageNum = nextPage
	}

	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageUrl, collectId), collectPageUrl, 0)
	if err != nil {
		err = fmt.Errorf("缓存采集当前页数失败 %v", err.Error())
		return
	}

	if collectPageUrl.PageNum > collectPageUrl.PageCount {
		err = resetPage(collect)
		if err != nil {
			return
		}
	}
	return
}

// 获取当前分页所有小说链接地址
func GetCollectBookUrl(collect *models.McCollect, isRestart bool) (bookUrl string, pageNum, pageCount int, err error) {
	collectId := collect.Id
	charset := collect.Charset
	urlComplete := collect.UrlComplete
	listSection := collect.ListSectionReg
	listUrl := collect.ListUrlReg
	urlReverse := collect.UrlReverse
	var pageUrl string
	pageUrl, pageNum, pageCount, err = GetCollectPageUrl(collect, isRestart)
	if err != nil {
		return
	}
	sleepSecond := getSleepSecond()
	var listContentHtml string
	listContentHtml, err = utils.GetHtml(pageUrl, charset, urlComplete, sleepSecond)
	if err != nil {
		err = fmt.Errorf("解析采集分页html出错 pageUrl=%v err=%v", pageUrl, err.Error())
		return
	}
	if listContentHtml == "" {
		//未获取到起始页面数据! 'state'=>'stop'
		err = fmt.Errorf("未获取到起始页面数据! pageUrl=%v", pageUrl)
		return
	}
	match := regexp.MustCompile(listSection).FindStringSubmatch(listContentHtml)
	if len(match) <= 0 {
		err = fmt.Errorf("%v", "获取采集列表规则失败")
		return
	}

	linkPattern := listUrl
	linksRe := regexp.MustCompile(linkPattern)
	links := linksRe.FindAllStringSubmatch(match[0], -1)
	if len(links) <= 0 {
		err = fmt.Errorf("采集列表失败 %v", err.Error())
		return
	}
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

	var pageBook = new(models.CollectPageBook)

	var CollectBookUrls []*models.CollectBookUrl
	for _, link := range bookLinks {
		collectData := &models.CollectBookUrl{
			Url:  link,
			Lock: 0,
		}
		CollectBookUrls = append(CollectBookUrls, collectData)
	}
	pageBook.PageNum = pageNum
	pageBook.PageCount = pageCount
	pageBook.BookCount = len(CollectBookUrls)
	pageBook.BookUrls = CollectBookUrls

	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectPageBookUrl, collectId), pageBook, 0)
	if err != nil {
		global.Collectlog.Errorf("缓存采集当前分页小说列表失败 %v", err.Error())
		return
	}
	bookUrl = CollectBookUrls[0].Url
	return
}

func SaveData(collect *models.McCollect, info *models.CollecBookInfoRes) (err error) {
	collectId := collect.Id
	sourceUrl := info.BookUrl
	bookName := info.BookName
	author := info.Author

	if len(info.Chapters) <= 0 {
		global.Collectlog.Errorf("采集章节为空 sourceUrl=%v bookName=%v", sourceUrl, bookName)
		return
	}

	var serialize int
	if strings.Contains(info.Serialize, "连载") {
		serialize = 1
	} else if strings.Contains(info.Serialize, "完本") {
		serialize = 2
	}
	chapterNum := len(info.Chapters)
	textNum := chapterNum * 2000

	info.Desc = utils.ReplaceText(info.Desc)

	//章节表
	var gq *gojsonq.JSONQ
	var chapterFile string
	gq, chapterFile, err = chapter_service.GetJsonqByBookName(info.BookName, info.Author)
	if err != nil {
		global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
		return
	}

	_, dbChapters, err := chapter_service.GetChapterNamesByFile(chapterFile)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}

	// 存储需要更新的章节列表
	var chapters []*models.CollectChapterInfo
	// 检查采集的章节是否需要更新
	for _, val := range info.Chapters {
		chapterLink := val.ChapterLink
		chapterName := val.ChapterTitle
		//lastSort += 1
		// 查找章节是否存在于数据库中
		found := false
		for _, dbChapter := range dbChapters {
			if chapterName == dbChapter {
				found = true
				break
			}
		}
		// 如果章节不存在于数据库中，则需要更新
		if !found {
			chapter := &models.CollectChapterInfo{
				ChapterLink:  chapterLink,
				ChapterTitle: chapterName,
			}
			chapters = append(chapters, chapter)
		}
	}
	//if len(updatedChapters) > 0 {
	//	batchSize := 1000 // 每批数据的大小
	//	if err = global.DB.Table(chapterTable).CreateInBatches(updatedChapters, batchSize).Error; err != nil {
	//		global.Collectlog.Errorf("sql 书籍章节添加失败，稍后再试 err=%v ", err.Error())
	//		return
	//	}
	//}

	collectChapter := models.CollectPageBookChapter{
		Collect:  collect,
		BookName: bookName,
		Author:   author,
		Chapters: chapters,
	}
	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectChapterUrl, collectId), collectChapter, 0)
	if err != nil {
		err = fmt.Errorf("缓存章节失败 err=%v", err.Error())
		return
	}

	threadMaxCount := utils.CollectThreadCount
	chapterCount := len(chapters)
	if chapterCount > 0 {
		var wg sync.WaitGroup
		// 计算线程数量
		if chapterCount > threadMaxCount {
			threadMaxCount = threadMaxCount
		}

		// 计算每个线程处理的任务数量
		tasksPerThread := chapterCount / threadMaxCount
		remainingTasks := chapterCount % threadMaxCount

		// 设置等待的任务数量
		wg.Add(chapterCount)

		// 分配任务给线程
		for i := 0; i < threadMaxCount; i++ {
			threadTasks := tasksPerThread
			if i < remainingTasks {
				threadTasks++
			}
			for j := 0; j < threadTasks; j++ {
				startIndex := j * threadTasks
				endIndex := (j + 1) * threadTasks
				if endIndex > chapterCount {
					startIndex = chapterCount
					endIndex = chapterCount
					wg.Done()
					continue
				}
				global.Collectlog.Errorf("正在采集 线程=%v 目标数量=%v 当前线程数量=%v startIndex=%v endIndex=%v", j, chapterCount, len(chapters[startIndex:endIndex]), startIndex, endIndex)
				go processChapter(collect, bookName, author, chapters[startIndex:endIndex], &wg)
			}
		}
		wg.Wait()
		global.Collectlog.Infof("小说%v章节采集任务完成 ", bookName)
	}

	//lastChapter := info.Chapters[len(info.Chapters)-1]

	for {
		isChapterEnd := getNextChapterLink(collectId)
		if isChapterEnd {
			break
		}
	}

	var updateChapterId int64
	var updateChapterTitle string
	var lastSortChapter *models.McBookChapter
	lastSortChapter, _ = chapter_service.GetLast(gq, "sort")
	if lastSortChapter != nil {
		updateChapterId = lastSortChapter.Id
		updateChapterTitle = lastSortChapter.ChapterName
	}

	msg := &models.NsqCollectBookPush{
		BookName:           bookName,
		Author:             author,
		Pic:                info.Pic,
		ClassId:            info.ClassId,
		CategoryName:       info.CategoryName,
		Desc:               info.Desc,
		Tags:               info.TagName,
		LastChapterTime:    info.UpdateTime,
		ChapterNum:         chapterNum,
		SourceId:           collectId,
		SourceUrl:          sourceUrl,
		UpdateChapterId:    updateChapterId,
		UpdateChapterTitle: updateChapterTitle,
		TextNum:            textNum,
		Serialize:          serialize,
	}

	err = NsqCollectBookPush(msg)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}
	return
}

func processChapter(collect *models.McCollect, bookName, author string, chapters []*models.CollectChapterInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(chapters) <= 0 {
		return
	}
	var err error
	for _, chapter := range chapters {
		_, err = CollectChapterText(collect, bookName, author, chapter.ChapterTitle, chapter.ChapterLink)
		if err != nil {
			global.Collectlog.Errorf("解析collect小说替换出错 collectId=%v err=%v", collect.Id, err.Error())
		}
	}
	return
}

func getNextChapterLink(collectId int64) (isEnd bool) {
	var err error
	var collectPageChapter = new(models.CollectPageBookChapter)
	chapterVal := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectChapterUrl, collectId))
	if chapterVal == "" {
		return
	}
	err = json.Unmarshal([]byte(chapterVal), &collectPageChapter)
	if err != nil {
		global.Collectlog.Errorf("解析采集当前章节链接失败 collectId=%v err=%v", collectId, err.Error())
		return
	}
	collect := collectPageChapter.Collect
	chapters := collectPageChapter.Chapters
	bookName := collectPageChapter.BookName
	author := collectPageChapter.Author
	if len(chapters) <= 0 {
		err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectChapterUrl, collectId))
		if err != nil {
			global.Collectlog.Errorf("删除当前章节链接失败 %v", err.Error())
			return
		}
		isEnd = true
		return
	}

	var chapterTitle, chapterLink string
	for _, val := range chapters {
		chapterTitle = val.ChapterTitle
		chapterLink = val.ChapterLink
		var textNum int
		textNum, err = CollectChapterText(collect, bookName, author, chapterTitle, chapterLink)
		if err != nil && textNum > 0 {
			global.Collectlog.Errorf("解析collect小说替换出错 collectId=%v err=%v", collect.Id, err.Error())
			utils.GetS5()
			return
		}
		updatedChapter := &models.McBookChapter{
			ChapterLink: chapterLink,
			ChapterName: chapterTitle,
			Vip:         0,
			Cion:        0,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
		}
		err = chapter_service.CreateChapter(bookName, author, updatedChapter)
		if err != nil {
			global.Collectlog.Errorf("%v", err.Error())
			return
		}
		chapterLinks := removeChapterUrl(chapters, chapterLink)
		collectPageChapter.Chapters = chapterLinks
		err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectChapterUrl, collectId), collectPageChapter, 0)
		if err != nil {
			global.Collectlog.Errorf("缓存采集当前章节链接失败 %v", err.Error())
			return
		}
		return
	}
	return
}

func CollectChapterText(collect *models.McCollect, bookName, author, chapterTitle, chapterLink string) (textNum int, err error) {
	charset := collect.Charset
	textReplace := collect.TextReplaceReg
	ChapterTextReg := collect.ChapterTextReg

	var replaces []*models.TextReplace
	if textReplace != "" {
		err = json.Unmarshal([]byte(textReplace), &replaces)
		if err != nil {
			err = fmt.Errorf("解析collect小说内容替换出错 collectId=%v err=%v", collect.Id, err.Error())
			return
		}
	}

	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterTitle)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		global.Collectlog.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterTitle, txtFile, textNum)
		return
	}

	var chapterHtml, text string
	chapterHtml, err = utils.GetHtml(chapterLink, charset, 0, 0)
	if err != nil {
		return
	}
	text, err = collect_service.GetChapterContent(ChapterTextReg, chapterHtml)
	if err != nil {
		global.Collectlog.Errorf("bookName=%v chapterTitle=%v chapterLink = %v err=%v", bookName, chapterTitle, chapterLink, err.Error())
		return
	}
	if len(replaces) > 0 {
		text = utils.ReplaceWords(text, replaces)
	}
	textNum = len([]rune(text))
	var chapterNameMd5 string
	chapterNameMd5, text, err = book_service.GetBookTxt(bookName, author, chapterTitle, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterTitle, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func NsqCollectBookPush(msg *models.NsqCollectBookPush) (err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(msg)
	if err != nil {
		err = fmt.Errorf("转换json数据失败: %v", err.Error())
		return
	}
	err = global.NsqPro.Publish(utils.UpdateBook, jsonData)
	if err != nil {
		err = fmt.Errorf("队列发送数据失败 %v", err.Error())
		return
	}
	return
}
