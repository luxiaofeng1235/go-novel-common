package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/mozillazg/go-pinyin"
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
	var siteUrl string = "https://www.bqg24.net"
	for {
		Bqg24Collect(siteUrl)
	}
}

func Bqg24Collect(siteUrl string) {
	categorys := Bqg24GetCacheCategory(siteUrl)
	if len(categorys) <= 0 {
		return
	}
	var category = new(models.Bqg24Category)
	for _, val := range categorys {
		if val.Use == 0 {
			category = val
			break
		}
	}
	Bqg24GetBooksByCategory(category)
	return
}

func Bqg24GetCacheCategory(siteUrl string) (categorys []*models.Bqg24Category) {
	categoryVal := redis_service.Get(utils.Bqg24Category)
	var err error
	if categoryVal == "" || categoryVal == "null" {
		categorys = Bqg24GetCategory(siteUrl)
		err = redis_service.Set(utils.Bqg24Category, categorys, 0)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
		return
	} else {
		err = json.Unmarshal([]byte(categoryVal), &categorys)
		if err != nil {
			global.Bqg24log.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func Bqg24GetCategory(siteUrl string) (categorys []*models.Bqg24Category) {
	var html string
	var err error
	html, err = utils.GetHtmlcolly(siteUrl)
	if err != nil {
		global.Bqg24log.Errorf("Bqg24GetCategory err:%v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		global.Bqg24log.Errorf("err:%v", err.Error())
		return
	}
	domain := utils.GetUrlDomain(siteUrl)
	document.Find(".nav_cont ul li").Each(func(i int, s *goquery.Selection) {
		categoryName := s.Find("a").Text()
		categoryHref, _ := s.Find("a").Attr("href")
		if categoryName == "" {
			return
		}
		if categoryName == "首页" {
			return
		}
		categoryKey := strings.Join(pinyin.LazyPinyin(categoryName, pinyin.NewArgs()), "")
		if categoryHref != "" {
			categoryHref = fmt.Sprintf("%v%v", domain, categoryHref)
		}
		category := &models.Bqg24Category{
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

func Bqg24GetBooksByCategory(category *models.Bqg24Category) {
	if category.CategoryHref == "" {
		err := redis_service.Del(utils.Bqg24Category)
		if err != nil {
			global.Bqg24log.Errorf("%v", err.Error())
			return
		}
	}
	pagebook, pageBookKey, err := Bqg24GetCacheCategoryBookList(category.CategoryKey, category.CategoryHref)
	if err != nil {
		global.Bqg24log.Errorf("%v", err.Error())
		return
	}
	if pagebook == nil {
		return
	}
	books := pagebook.Books
	var bookInfo *models.Bqg24BookDesc
	var isSkip bool
	bookInfo, isSkip, err = Bqg24GetBookInfo(books)
	if isSkip {
		err = Bqg24RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Bqg24log.Errorf("%v", err.Error())
			return
		}
	}
	if bookInfo == nil {
		global.Bqg24log.Errorln("bookInfo为nil")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			err = Bqg24RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
			if err != nil {
				global.Bqg24log.Errorf("%v", err.Error())
				return
			}
		}
	}
	chapters, isSkip, err := Bqg24GetFilterChapter(pageBookKey, pagebook, bookInfo)
	if err != nil {
		global.Bqg24log.Errorf("Bqg24GetFilterChapter err=%v", err.Error())
		return
	}
	if isSkip {
		return
	}
	err = Bqg24SaveChapterJson(pageBookKey, pagebook, bookInfo, chapters)
	if err != nil {
		return
	}
}

func Bqg24GetCacheCategoryBookList(categoryKey, categoryHref string) (pagebook *models.Bqg24PageBooks, pageBookKey string, err error) {
	pageBookKey = fmt.Sprintf("%v_%v", utils.Bqg24Books, categoryKey)
	pagebookVal := redis_service.Get(pageBookKey)
	if pagebookVal == "" || pagebookVal == "null" {
		pagebook, err = Bqg24GetPageBookList(categoryHref)
		if err != nil {
			global.Bqg24log.Errorf("%v", err.Error())
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
			global.Bqg24log.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
		book := pagebook.Books
		nextPageLink := pagebook.NextLink
		for _, val := range book {
			if val.Use == 0 {
				return
			}
		}

		if nextPageLink != "" {
			pagebook, err = Bqg24GetPageBookList(nextPageLink)
			if err != nil {
				return
			}
			err = redis_service.Set(pageBookKey, pagebook, 0)
			if err != nil {
				err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
				return
			}
			if len(pagebook.Books) <= 0 {
				err = Bqg24RemoveCategory(pageBookKey, categoryHref)
				if err != nil {
					return
				}
			}
		} else {
			err = Bqg24RemoveCategory(pageBookKey, categoryHref)
			if err != nil {
				return
			}
		}
	}
	return
}

func Bqg24RemoveCategory(pageBookKey, categoryHref string) (err error) {
	err = redis_service.Del(pageBookKey)
	if err != nil {
		global.Bqg24log.Errorf("删除pageBookKey=%v 出错 %v", pageBookKey, err.Error())
		return
	}
	categoryVal := redis_service.Get(utils.Bqg24Category)
	var categorys []*models.Bqg24Category
	err = json.Unmarshal([]byte(categoryVal), &categorys)
	if err != nil {
		global.Bqg24log.Errorf("获取分类缓存失败 err=%v", err.Error())
		return
	}
	for _, category := range categorys {
		if category.CategoryHref == categoryHref {
			category.Use = 1
			break
		}
	}
	err = redis_service.Set(utils.Bqg24Category, categorys, 0)
	if err != nil {
		err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
		return
	}
	return
}
func Bqg24RemoveBook(pagebook *models.Bqg24PageBooks, pageBookKey, bookName, author, sourceUrl string, isSkip bool) (err error) {
	books := pagebook.Books
	for _, val := range books {
		if val.BookLink == sourceUrl {
			val.Use = 1
			if isSkip {
				val.Use = 2
			}
			break
		}
	}
	pagebook.Books = books
	err = redis_service.Set(pageBookKey, pagebook, 0)
	if err != nil {
		global.Bqg24log.Errorf("Bqg24RemoveBook bookName=%v author=%v sourceUrl=%v err=%v", bookName, author, sourceUrl, err.Error())
		return
	}
	return
}

func Bqg24GetBookInfo(books []*models.Bqg24PageBook) (bookInfo *models.Bqg24BookDesc, isSkip bool, err error) {
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
	bookInfo, isSkip, err = Bqg24GetBookDesc(bookLink)
	if bookInfo == nil {
		bookInfo = new(models.Bqg24BookDesc)
		bookInfo.BookName = bookName
		bookInfo.Author = author
		bookInfo.SourceUrl = bookLink
	}
	return
}

func Bqg24GetFilterChapter(pageBookKey string, pagebook *models.Bqg24PageBooks, bookInfo *models.Bqg24BookDesc) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if pagebook == nil || bookInfo == nil {
		global.Bqg24log.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	sourceUrl := bookInfo.SourceUrl
	chapters, isSkip, err = Bqg24GetChapters(sourceUrl)
	if err != nil {
		global.Bqg24log.Errorf("%v", err.Error())
		if bookInfo == nil {
			global.Bqg24log.Errorf("%v", err.Error())
			return
		}
		err = Bqg24RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Bqg24log.Errorf("%v", err.Error())
			return
		}
		return
	}
	if isSkip {
		err = Bqg24RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Bqg24log.Errorf("%v", err.Error())
			return
		}
		return
	}
	return
}

func Bqg24SaveChapterJson(pageBookKey string, pagebook *models.Bqg24PageBooks, bookInfo *models.Bqg24BookDesc, chapters []*models.CollectChapterInfo) (err error) {
	if pagebook == nil || bookInfo == nil {
		global.Bqg24log.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	bookName := bookInfo.BookName
	author := bookInfo.Author
	sourceUrl := bookInfo.SourceUrl
	Bqg24ManyThreadChapter(bookName, author, chapters)
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
		textNum, err = Bqg24CollectChapterText(bookName, author, chapterName, chapterLink)
		if err != nil {
			global.Bqg24log.Errorf("获取章节内容失败 bookName=%v author=%v chapterName=%v chapterLink=%v err=%v", bookName, author, chapterName, chapterLink, err.Error())
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

	err = Bqg24RemoveBook(pagebook, pageBookKey, bookName, author, sourceUrl, false)
	if err != nil {
		return
	}

	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Bqg24log.Errorf("获取JSONQ对象失败 %v", err.Error())
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
		global.Bqg24log.Errorf("%v", err.Error())
		return
	}
	return
}

func Bqg24ManyThreadChapter(bookName, author string, chapters []*models.CollectChapterInfo) {
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
			_, err = Bqg24CollectChapterTextThread(bookName, author, chapterName, chapterLink)
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

func Bqg24GetTotalPage(categoryHref string) (totalPage int) {
	var err error
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(categoryHref)
	if err != nil {
		return
	}
	if tempHtml == "" {
		global.Bqg24log.Errorln("GetTotalPage html为空")
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		global.Bqg24log.Errorf("GetTotalPage err:%v", err.Error())
		return
	}

	pageString := document.Find(".articlepage em").Text()
	pageInfo := strings.Split(pageString, "/")
	if len(pageInfo) > 1 {
		totalPageStr := utils.GetUrlBookNum(pageInfo[1])
		// 转换为整数
		totalPage, err = strconv.Atoi(totalPageStr)
		if err != nil {
			global.Bqg24log.Errorf("无法解析页数:%v", err.Error())
			return
		}
	}
	return
}

func Bqg24GetPageBookList(pageLink string) (pagebook *models.Bqg24PageBooks, err error) {
	if pageLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(pageLink)
	if html == "" {
		global.Bqg24log.Errorf("GetHtmlcolly为空 %v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var books []*models.Bqg24PageBook
	cateNum, pageNum := Bqg24GetPageNum(pageLink)
	nextLink, _ := document.Find("#pagelink .next").Attr("href")
	if nextLink != "" {
		if !strings.Contains(nextLink, "http") {
			domain := utils.GetUrlDomain(pageLink)
			nextLink = fmt.Sprintf("%v%v", domain, nextLink)
		}
	}
	totalPage := document.Find("#pagelink .last").Text()
	totalPageInt, _ := strconv.Atoi(totalPage)
	document.Find("#mm_14 ul li").Each(func(i int, s *goquery.Selection) {
		bookLink, _ := s.Find(".sp_2 a").Attr("href")
		bookName := s.Find(".sp_2 a").Text()
		author := s.Find(".sp_4").Text()
		if bookName == "" || author == "" {
			return
		}
		book := &models.Bqg24PageBook{
			BookName: bookName,
			Author:   author,
			BookLink: bookLink,
		}
		books = append(books, book)
		return
	})
	pagebook = &models.Bqg24PageBooks{
		PageLink:  pageLink,
		NextLink:  nextLink,
		CateNum:   cateNum,
		PageNum:   pageNum,
		TotalPage: totalPageInt,
		BookNum:   len(books),
		Books:     books,
	}
	return
}

func Bqg24GetChapters(bookLink string) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Bqg24log.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if html == "" {
		global.Bqg24log.Errorf("GetHtmlcolly为空 %v", bookLink)
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var chaptersAll []*models.CollectChapterInfo
	document.Find(".mulu_list li").Each(func(i int, s *goquery.Selection) {
		chapterHref, _ := s.Find("a").Attr("href")
		if chapterHref == "" {
			return
		}
		chapterHref = fmt.Sprintf("%v%v", bookLink, chapterHref)
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

func Bqg24CollectChapterText(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Bqg24log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := Bqg24GetChapterText(chapterLink)
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

func Bqg24CollectChapterTextThread(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Bqg24log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := Bqg24GetChapterText(chapterLink)
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

func Bqg24GetChapterText(chapterLink string) (text string) {
	var tempHtml string
	var err error
	tempHtml, err = utils.GetHtmlcolly(chapterLink)
	if err != nil {
		global.Bqg24log.Errorf("GetHtmlcolly为空 %v %v", chapterLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Bqg24log.Errorf("GetHtmlcolly为空 %v", chapterLink)
		return
	}

	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	tempText, _ := document.Find("#htmlContent").Html()
	tempText = strings.ReplaceAll(tempText, "一秒记住【笔趣阁 www.bqg24.net】，精彩小说无弹窗免费阅读！", "")
	tempText = strings.ReplaceAll(tempText, "<br/><br/>", "<br/>")
	tempText = strings.ReplaceAll(tempText, "<br/>", "\r\n")
	tempText = strings.ReplaceAll(tempText, "&amp;8195；", "")
	text = strings.TrimSpace(tempText)
	return
}

func Bqg24GetBookDesc(bookLink string) (info *models.Bqg24BookDesc, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Bqg24log.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Bqg24log.Errorf("GetHtmlcolly为空 %v %v", bookLink, err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}

	bookName, _ := document.Find("meta[property='og:title']").Attr("content")
	author, _ := document.Find("meta[property='og:novel:author']").Attr("content")
	desc, _ := document.Find("meta[property='og:description']").Attr("content")
	pic, _ := document.Find("#fmimg img").Attr("src")
	//解析图片目录规则：日期（一周一个目录）/对应的图片信息路径
	mondayDate := utils.GetThisWeekFirstDate()     //获取每周的周一的日期
	savePicPath := "/data/pic/" + mondayDate + "/" //拼装存储的图片路径
	filePath, err := utils.DownImg(bookName, author, pic, savePicPath)
	pic = strings.TrimLeft(filePath, ".")
	categoryName, _ := document.Find("meta[property='og:novel:category']").Attr("content")

	serializeName, _ := document.Find("meta[property='og:novel:status']").Attr("content")
	if serializeName != "" {
		serializeName = strings.ReplaceAll(serializeName, "小说状态：", "")
	}
	lastChapterTitle, _ := document.Find("meta[property='og:novel:latest_chapter_name']").Attr("content")
	lastChapterTimeStr, _ := document.Find("meta[property='og:novel:update_time']").Attr("content")
	lastChapterTime := utils.DateToUnix(lastChapterTimeStr)

	categorys := []*models.CategoryReg{
		{"玄幻小说", 12},
		{"修真小说", 4},
		{"都市小说", 2},
		{"历史小说", 23},
		{"网游小说", 8},
		{"科幻小说", 63},
	}
	classId := utils.CategoryEquiv(categorys, categoryName)
	//log.Println("sourceUrl", bookLink)
	//log.Println("bookName", bookName)
	//log.Println("author", author)
	//log.Println("pic", pic)
	//log.Println("desc", desc)
	//log.Println("categoryName", categoryName)
	//log.Println("serializeName", serializeName)
	//log.Println("lastChapterTime", lastChapterTimeStr, lastChapterTime)
	//log.Println("lastChapterTitle", lastChapterTitle)
	//log.Println("classId", classId)
	var serialize = 1
	if strings.Contains(serializeName, "已完成") {
		serialize = 2
	}
	if bookName == "" || author == "" {
		global.Bqg24log.Errorf("小说名称或作者为空 %v %v %v", bookName, author, bookLink)
		isSkip = true
		return
	}
	info = &models.Bqg24BookDesc{
		BookName:         bookName,
		Author:           author,
		Pic:              pic,
		Desc:             desc,
		SourceUrl:        bookLink,
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

func Bqg24GetPageNum(pageUrl string) (cateNum, pageNum int) {
	urlObj, _ := url.Parse(pageUrl)
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(urlObj.Path, -1)
	if len(matches) > 0 {
		cateNum, _ = strconv.Atoi(matches[0])
	}
	if len(matches) > 1 {
		pageNum, _ = strconv.Atoi(matches[1])
	}
	return
}
