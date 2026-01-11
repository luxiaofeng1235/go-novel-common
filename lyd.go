package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/ants/v2"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/collect/collect_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)
	addr, passwd, defaultdb := db.GetRedis()
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	db.InitNsqProducer()
	db.InitNsqConsumer()
	var siteUrl string = "https://www.27k.net/sort/"
	for {
		LydCollect(siteUrl)
	}
}

func LydCollect(siteUrl string) {
	categorys := LydGetCacheCategory(siteUrl)
	if len(categorys) <= 0 {
		return
	}
	var category = new(models.LydCategory)
	for _, val := range categorys {
		if val.Use == 0 {
			category = val
			break
		}
	}
	LydGetBooksByCategory(category)
	return
}

func LydGetBooksByCategory(category *models.LydCategory) {
	if category.CategoryHref == "" {
		err := redis_service.Del(utils.LydCategory)
		if err != nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
	}
	pagebook, pageBookKey, err := LydGetCacheCategoryBookList(category.CategoryKey, category.CategoryHref)
	if err != nil {
		global.Lydlog.Errorf("%v", err.Error())
		return
	}
	if pagebook == nil {
		return
	}
	books := pagebook.Books
	var bookInfo *models.LydBookDesc
	var isSkip bool
	bookInfo, isSkip, err = LydGetBookInfo(books)
	if isSkip {
		err = LydRemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.BookLink, true)
		if err != nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
	}
	if bookInfo == nil {
		global.Lydlog.Errorln("bookInfo为nil")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			err = LydRemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.BookLink, true)
			if err != nil {
				global.Lydlog.Errorf("%v", err.Error())
				return
			}
		}
	}
	chapters, isSkip, err := LydGetFilterChapter(pageBookKey, pagebook, bookInfo)
	if err != nil {
		global.Lydlog.Errorf("LydGetFilterChapter err=%v", err.Error())
		return
	}
	if isSkip {
		global.Lydlog.Errorf("isSkip 跳过 %v %v", bookInfo.BookName, bookInfo.BookLink)
		return
	}
	err = LydSaveChapterJson(pageBookKey, pagebook, bookInfo, chapters)
	if err != nil {
		return
	}
}

func LydRemoveBook(pagebook *models.LydPageBooks, pageBookKey, bookName, author, bookLink string, isSkip bool) (err error) {
	books := pagebook.Books
	for _, val := range books {
		if val.BookLink == bookLink {
			val.Use = 1
			if isSkip {
				val.Use = 2
			}
		}
	}
	err = redis_service.Set(pageBookKey, pagebook, 0)
	if err != nil {
		global.Lydlog.Errorf("LydRemoveBook bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}
	return
}

func LydGetBookInfo(books []*models.LydPageBook) (bookInfo *models.LydBookDesc, isSkip bool, err error) {
	if len(books) <= 0 {
		return
	}
	var bookLink, bookName, author string
	for _, val := range books {
		if val.Use == 0 {
			bookLink = val.BookLink
			bookName = val.BookName
			author = val.Author
			break
		}
	}
	bookInfo, isSkip, err = LydGetBookDesc(bookLink)
	if bookInfo == nil {
		bookInfo = new(models.LydBookDesc)
		bookInfo.BookName = bookName
		bookInfo.Author = author
		bookInfo.BookLink = bookLink
		bookInfo.SourceUrl = ""
		isSkip = true
	}
	return
}

func LydGetFilterChapter(pageBookKey string, pagebook *models.LydPageBooks, bookInfo *models.LydBookDesc) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if pagebook == nil || bookInfo == nil {
		global.Lydlog.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	sourceUrl := bookInfo.SourceUrl
	chapters, isSkip, err = LydGetChapters(sourceUrl)
	if err != nil {
		global.Lydlog.Errorf("%v", err.Error())
		if bookInfo == nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
		err = LydRemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.BookLink, true)
		if err != nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
		return
	}
	if isSkip {
		log.Println(bookInfo.BookName, bookInfo.Author, bookInfo.BookLink)
		err = LydRemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.BookLink, true)
		if err != nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
		return
	}
	return
}

func LydSaveChapterJson(pageBookKey string, pagebook *models.LydPageBooks, bookInfo *models.LydBookDesc, chapters []*models.CollectChapterInfo) (err error) {
	if pagebook == nil || bookInfo == nil {
		global.Lydlog.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	bookName := bookInfo.BookName
	author := bookInfo.Author
	//sourceUrl := bookInfo.SourceUrl
	bookLink := bookInfo.BookLink
	LydManyThreadChapter(bookName, author, chapters)
	var chapterAll []*models.McBookChapter
	var id int64
	var sort int
	var totalTextNum int
	for _, val := range chapters {
		id++
		sort++
		chapterName := val.ChapterTitle
		chapterLink := val.ChapterLink
		var textNum, isLess int
		textNum, err = LydCollectChapterText(bookName, author, chapterName, chapterLink)
		if err != nil {
			global.Lydlog.Errorf("获取章节内容失败 bookName=%v author=%v chapterName=%v chapterLink=%v err=%v", bookName, author, chapterName, chapterLink, err.Error())
			continue
		}
		if textNum <= 1000 {
			isLess = 1
		}
		if textNum <= 100 {
			continue
		}
		totalTextNum += textNum
		updatedChapter := &models.McBookChapter{
			Id:          id,
			ChapterLink: chapterLink,
			ChapterName: chapterName,
			IsLess:      isLess,
			Sort:        sort,
			Vip:         0,
			Cion:        0,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
		}
		chapterAll = append(chapterAll, updatedChapter)
	}
	if len(chapterAll) <= 0 {
		return
	}
	bookInfo.TextNum = totalTextNum
	bookInfo.ChapterNum = len(chapterAll)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	newJsonData, err := json.MarshalIndent(chapterAll, "", "  ")
	if err != nil {
		err = fmt.Errorf("美化章节格式错误 %v", "")
		return
	}
	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	err = utils.WriteFile(chapterFile, string(newJsonData))

	err = LydRemoveBook(pagebook, pageBookKey, bookName, author, bookLink, false)
	if err != nil {
		return
	}

	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Lydlog.Errorf("获取JSONQ对象失败 %v", err.Error())
		return
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
		Pic:                bookInfo.Pic,
		ClassId:            bookInfo.ClassId,
		CategoryName:       bookInfo.CategoryName,
		Tags:               bookInfo.CategoryName,
		Desc:               bookInfo.Desc,
		ChapterNum:         bookInfo.ChapterNum,
		SourceId:           0,
		SourceUrl:          bookInfo.SourceUrl,
		LastChapterTitle:   bookInfo.LastChapterTitle,
		LastChapterTime:    bookInfo.LastChapterTime,
		UpdateChapterId:    updateChapterId,
		UpdateChapterTitle: updateChapterTitle,
		UpdateChapterTime:  utils.GetUnix(),
		TextNum:            bookInfo.TextNum,
		Serialize:          bookInfo.Serialize,
		IsClassic:          bookInfo.IsClassic,
	}

	err = collect_service.NsqCollectBookPush(msg)
	if err != nil {
		global.Lydlog.Errorf("%v", err.Error())
		return
	}
	return
}

func LydManyThreadChapter(bookName, author string, chapters []*models.CollectChapterInfo) {
	var err error
	pool, err := ants.NewPool(30)
	if err != nil {
		return
	}
	defer pool.Release()
	var wg sync.WaitGroup
	for _, val := range chapters {
		wg.Add(1)
		err = pool.Submit(func() {
			defer wg.Done()
			chapterName := val.ChapterTitle
			chapterLink := val.ChapterLink
			_, err = LydCollectChapterTextThread(bookName, author, chapterName, chapterLink)
			if err != nil {
				return
			}
		})
		if err != nil {
			return
		}
	}
	wg.Wait()
	return
}

func LydGetCacheCategoryBookList(categoryKey, categoryHref string) (pagebook *models.LydPageBooks, pageBookKey string, err error) {
	pageBookKey = fmt.Sprintf("%v_%v", utils.LydBooks, categoryKey)
	pagebookVal := redis_service.Get(pageBookKey)
	if pagebookVal == "" || pagebookVal == "null" {
		pagebook, err = LydGetPageBookList(categoryHref)
		if err != nil {
			global.Lydlog.Errorf("%v", err.Error())
			return
		}
		err = redis_service.Set(pageBookKey, pagebook, 0)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
	} else {
		err = json.Unmarshal([]byte(pagebookVal), &pagebook)
		if err != nil {
			global.Lydlog.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
		book := pagebook.Books
		totalPage := pagebook.TotalPage
		nextPageNum := pagebook.Page + 1
		for _, val := range book {
			if val.Use == 0 {
				return
			}
		}
		domain := utils.GetUrlDomain(categoryHref)
		var nextPageLink string
		if nextPageNum > 0 {
			cateNum := LydGetCateNum(categoryHref)
			nextPageLink = fmt.Sprintf("%v/sort/%v/%v.html", domain, cateNum, nextPageNum)
			if nextPageNum <= totalPage {
				pagebook, err = LydGetPageBookList(nextPageLink)
				if err != nil {
					return
				}
				err = redis_service.Set(pageBookKey, pagebook, 0)
				if err != nil {
					err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
					return
				}
			} else {
				err = LydRemoveCategory(pageBookKey, categoryHref)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func LydRemoveCategory(pageBookKey, categoryHref string) (err error) {
	err = redis_service.Del(pageBookKey)
	if err != nil {
		global.Lydlog.Errorf("删除pageBookKey=%v 出错 %v", pageBookKey, err.Error())
		return
	}
	categoryVal := redis_service.Get(utils.LydCategory)
	var categorys []*models.LydCategory
	err = json.Unmarshal([]byte(categoryVal), &categorys)
	if err != nil {
		global.Lydlog.Errorf("获取分类缓存失败 err=%v", err.Error())
		return
	}
	for _, category := range categorys {
		if category.CategoryHref == categoryHref {
			category.Use = 1
			break
		}
	}
	err = redis_service.Set(utils.LydCategory, categorys, 0)
	if err != nil {
		err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
		return
	}
	return
}

func LydGetCacheCategory(siteUrl string) (categorys []*models.LydCategory) {
	categoryVal := redis_service.Get(utils.LydCategory)
	var err error
	if categoryVal == "" || categoryVal == "null" {
		categorys = LydGetCategory(siteUrl)
		err = redis_service.Set(utils.LydCategory, categorys, time.Hour*24*30)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
	} else {
		err = json.Unmarshal([]byte(categoryVal), &categorys)
		if err != nil {
			global.Lydlog.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func LydGetCategory(siteUrl string) (categorys []*models.LydCategory) {
	var html string
	var err error
	html, err = utils.GetHtmlcolly(siteUrl)
	if err != nil {
		global.Lydlog.Errorf("LydGetCategory err:%v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		global.Lydlog.Errorf("err:%v", err.Error())
		return
	}
	domain := utils.GetUrlDomain(siteUrl)
	document.Find(".weekl_yrank ul li").Each(func(i int, s *goquery.Selection) {
		categoryName := s.Find("a").Text()
		categoryHref, _ := s.Find("a").Attr("href")
		categoryKey, _ := s.Find("a").Attr("id")
		if categoryName == "" || categoryKey == "" {
			return
		}
		if categoryHref != "" {
			categoryHref = fmt.Sprintf("%v%v", domain, categoryHref)
		}
		category := &models.LydCategory{
			CategoryName: categoryName,
			CategoryKey:  categoryKey,
			CategoryHref: categoryHref,
			Use:          0,
		}
		categorys = append(categorys, category)
		return
	})
	return
}

func LydGetTotalPage(categoryHref string) (totalPage int) {
	var err error
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(categoryHref)
	if err != nil {
		return
	}
	if tempHtml == "" {
		global.Lydlog.Errorln("GetTotalPage html为空")
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		global.Lydlog.Errorf("GetTotalPage err:%v", err.Error())
		return
	}

	pageString := document.Find(".articlepage em").Text()
	pageInfo := strings.Split(pageString, "/")
	if len(pageInfo) > 1 {
		totalPageStr := utils.GetUrlBookNum(pageInfo[1])
		// 转换为整数
		totalPage, err = strconv.Atoi(totalPageStr)
		if err != nil {
			global.Lydlog.Errorf("无法解析页数:%v", err.Error())
			return
		}
	}
	return
}

func LydGetPageBookList(pageLink string) (pagebook *models.LydPageBooks, err error) {
	if pageLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(pageLink)
	if html == "" {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	domain := utils.GetUrlDomain(pageLink)
	var books []*models.LydPageBook
	pageNum := LydGetPageNum(pageLink)
	totalPage := document.Find("#pagelink").Find("a").Last().Text()
	totalPageInt, _ := strconv.Atoi(totalPage)
	document.Find("#article_list_content li").Each(func(i int, s *goquery.Selection) {
		pic, _ := s.Find("img").Attr("src")
		bookLink, _ := s.Find("a").Attr("href")
		bookName := s.Find(".newnav h3 a").Text()
		author := s.Find(".labelbox label").First().Text()
		if bookLink != "" {
			bookLink = fmt.Sprintf("%v%v", domain, bookLink)
		}
		book := &models.LydPageBook{
			BookName: bookName,
			Author:   author,
			Pic:      pic,
			BookLink: bookLink,
		}
		books = append(books, book)
		return
	})
	pagebook = &models.LydPageBooks{
		PageLink:  pageLink,
		Page:      pageNum,
		TotalPage: totalPageInt,
		BookNum:   len(books),
		Books:     books,
	}
	return
}

func LydGetChapters(bookLink string) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Lydlog.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if html == "" {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v", bookLink)
		return
	}
	domain := utils.GetUrlDomain(bookLink)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var chaptersAll []*models.CollectChapterInfo
	document.Find("#chapterList li").Each(func(i int, s *goquery.Selection) {
		chapterHref, _ := s.Find("a").Attr("href")
		if chapterHref != "" {
			chapterHref = fmt.Sprintf("%v%v", domain, chapterHref)
		}
		chapterName := s.Find("a").Text()
		chapter := &models.CollectChapterInfo{
			ChapterLink:  chapterHref,
			ChapterTitle: chapterName,
		}
		chaptersAll = append(chaptersAll, chapter)
		return
	})
	if len(chaptersAll) <= 100 {
		isSkip = true
		return
	}
	for _, val := range chaptersAll {
		isChapterName := utils.IsChapterName(val.ChapterTitle)
		if !isChapterName {
			continue
		}
		chapters = append(chapters, val)
	}
	lenAll := len(chaptersAll)
	lenc := len(chapters)

	per := int(float64(lenc) / float64(lenAll) * 100.0)
	if per < 85 {
		isSkip = true
	}
	return
}

func LydCollectChapterText(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Lydlog.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := LydGetChapterText(chapterLink)
	textNum = len([]rune(text))
	if textNum <= 10 {
		err = fmt.Errorf("获取内容失败 %v", text)
		return
	}
	//log.Println("text", bookName, author, chapterName, chapterLink, text, textNum)
	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func LydCollectChapterTextThread(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Lydlog.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := LydGetChapterText(chapterLink)
	textNum = len([]rune(text))
	if textNum <= 10 {
		err = fmt.Errorf("获取内容失败 %v", text)
		return
	}
	//log.Println("text", bookName, author, chapterName, chapterLink, text, textNum)
	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func LydGetChapterText(chapterLink string) (text string) {
	var tempHtml string
	var err error
	tempHtml, err = utils.GetHtmlcolly(chapterLink)
	if err != nil {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v %v", chapterLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v", chapterLink)
		return
	}

	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	chapterName := document.Find(".txtnav h1").Text()
	chapterTime := document.Find(".txtinfo span").First().Text()
	author := document.Find(".txtinfo span").Last().Text()
	tempText := document.Find(".txtnav").Text()
	tempText = strings.ReplaceAll(tempText, chapterName, "")
	tempText = strings.ReplaceAll(tempText, fmt.Sprintf("%v %v", chapterTime, author), "")
	tempText = strings.ReplaceAll(tempText, " 笔趣阁顶点.，最快万相之王！", "")
	tempText = strings.ReplaceAll(tempText, "很多书已经很难再找到，且看且珍惜吧", "")
	tempText = strings.ReplaceAll(tempText, "<script>loadAdv(2,0);</script>", "")
	tempText = strings.ReplaceAll(tempText, "loadAdv(2,0);", "")
	tempText = strings.ReplaceAll(tempText, "loadAdv(3,0);", "")
	text = strings.TrimSpace(tempText)
	return
}

func LydGetBookDesc(bookLink string) (info *models.LydBookDesc, isSkip bool, err error) {
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Lydlog.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v %v", bookLink, err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	domain := utils.GetUrlDomain(bookLink)
	sourceUrl, _ := document.Find("#a_addbookcase").Attr("href")
	if sourceUrl != "" {
		sourceUrl = fmt.Sprintf("%v%v", domain, sourceUrl)
	}
	tempHtml, err = utils.GetHtmlcolly(sourceUrl)
	if err != nil {
		global.Lydlog.Errorf("%v %v", sourceUrl, err.Error())
		return
	}
	if tempHtml == "" {
		global.Lydlog.Errorf("GetHtmlcolly为空 %v %v", sourceUrl, err.Error())
		return
	}
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	bookName := document.Find("#article_list_content .newnav h3").Text()
	author, _ := document.Find("#article_list_content label a").Attr("title")
	if author != "" {
		author = strings.ReplaceAll(author, "作者：", "")
	}
	desc := document.Find(".ellipsis_2").Text()

	pic, _ := document.Find("#article_list_content .imgbox img").Attr("src")
	filePath, err := utils.DownImg(bookName, author, pic, "/data/pic/")
	pic = strings.TrimLeft(filePath, ".")
	categoryName := document.Find("#article_list_content .newnav .labelbox label").Eq(1).Find("a").Text()

	serializeName := document.Find("#article_list_content .newnav .labelbox label").Last().Text()
	if serializeName != "" {
		serializeName = strings.ReplaceAll(serializeName, "小说状态：", "")
	}
	lastChapterTitle := document.Find("#article_list_content .zxzj h4").First().Find("a").Text()
	lastChapterTimeStr := document.Find("#article_list_content .zxzj h4").Last().Find("p").Text()
	if lastChapterTimeStr != "" {
		lastChapterTimeStr = strings.ReplaceAll(lastChapterTimeStr, "更新时间：", "")
	}
	lastChapterTime := utils.DateToUnix(lastChapterTimeStr)

	categorys := []*models.CategoryReg{
		{"玄幻小说", 12},
		{"都市小说", 2},
		{"历史小说", 23},
		{"科幻小说", 63},
		{"悬疑小说", 24},
		{"网游小说", 8},
		{"穿越小说", 21},
		{"现代言情小说", 45},
		{"古代言情小说", 45},
		{"豪门总裁小说", 64},
		{"青春校园小说", 41},
		{"其他类别小说", 2},
	}
	classId := utils.CategoryEquiv(categorys, categoryName)
	//log.Println("sourceUrl", sourceUrl)
	//log.Println("bookName", bookName)
	//log.Println("author", author)
	//log.Println("pic", pic)
	//log.Println("desc", desc)
	//log.Println("categoryName", categoryName)
	//log.Println("serializeName", serializeName)
	//log.Println("lastChapterTime", lastChapterTimeStr, lastChapterTime)
	//log.Println("lastChapterTitle", lastChapterTitle)
	var serialize = 1
	if strings.Contains(serializeName, "完本") {
		serialize = 2
	}
	//log.Println("serialize", serialize)
	if bookName == "" || author == "" {
		global.Lydlog.Errorf("小说名称或作者为空 %v %v %v", bookName, author, bookLink)
		isSkip = true
		return
	}
	info = &models.LydBookDesc{
		BookName:         bookName,
		Author:           author,
		Pic:              pic,
		Desc:             desc,
		BookLink:         bookLink,
		SourceUrl:        sourceUrl,
		Serialize:        serialize,
		CategoryName:     categoryName,
		ClassId:          classId,
		LastChapterTime:  utils.UnixToDatetime(lastChapterTime),
		LastChapterTitle: lastChapterTitle,
		IsClassic:        0,
		Use:              0,
	}
	return
}

func LydGetPageNum(pageUrl string) (num int) {
	// 创建正则表达式模式，匹配数字部分
	re := regexp.MustCompile(`(\d+)\.html`)
	matches := re.FindStringSubmatch(pageUrl)
	if len(matches) > 1 {
		num, _ = strconv.Atoi(matches[1])
	}
	return
}

func LydGetCateNum(pageUrl string) (num int) {
	urlObj, _ := url.Parse(pageUrl)
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(urlObj.Path, -1)
	if len(matches) > 0 {
		num, _ = strconv.Atoi(matches[0])
	}
	return
}
