package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/collect/collect_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
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
	var siteUrl string = "https://www.xsw.tw/"
	var defaultCategoryName string = "周点击榜"
	for {
		CollectXsw(siteUrl, defaultCategoryName)
		//time.Sleep(time.Millisecond * 200)
	}
}
func CollectXsw(siteUrl, defaultCategoryName string) {
	categoryKey, categoryHref := XswGetCacheCategory(siteUrl, defaultCategoryName)
	pagebook, keyName, err := XswGetCacheCategoryBookList(categoryKey, categoryHref)
	if err != nil {
		return
	}
	books := pagebook.Books
	var bookInfo *models.XswBookDesc
	bookInfo, err = XswGetBookInfo(books)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			if bookInfo == nil {
				global.Xswlog.Errorf("%v", err.Error())
				return
			}
			err = XswSkipBook(pagebook, keyName, bookInfo.BookName, bookInfo.Author)
			if err != nil {
				global.Xswlog.Errorf("%v", err.Error())
				return
			}
		}
	}
	chapters, isSkip, err := XswGetFilterChapter(keyName, pagebook, bookInfo)
	if err != nil {
		return
	}
	if isSkip {
		return
	}
	err = XswSaveChapterJson(keyName, pagebook, bookInfo, chapters)
	if err != nil {
		return
	}
}

func XswSkipBook(pagebook *models.XswPageBooks, keyName, bookName, author string) (err error) {
	books := pagebook.Books
	for _, val := range books {
		if val.BookName == bookName && val.Author == author {
			val.Use = 1
			break
		}
	}
	pagebook.Books = books
	err = redis_service.Set(keyName, pagebook, 0)
	if err != nil {
		global.Xswlog.Errorf("缓存采集状态失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}
	return
}

func XswGetBookInfo(books []*models.XswPageBook) (bookInfo *models.XswBookDesc, err error) {
	if len(books) <= 0 {
		return
	}
	var bookLink, bookName, author string
	for _, val := range books {
		if val.Use == 0 {
			bookLink = val.Link
			bookName = val.BookName
			author = val.Author
			break
		}
	}
	bookInfo, err = GetBookDesc(bookLink)
	if bookInfo == nil {
		bookInfo = new(models.XswBookDesc)
		bookInfo.BookName = bookName
		bookInfo.Author = author
	}
	return
}

func XswGetFilterChapter(keyName string, pagebook *models.XswPageBooks, bookInfo *models.XswBookDesc) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if pagebook == nil || bookInfo == nil {
		global.Xswlog.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	sourceUrl := bookInfo.SourceUrl
	chapters, isSkip, err = XswGetChapters(sourceUrl)
	if err != nil {
		global.Xswlog.Errorf("%v", err.Error())
		if bookInfo == nil {
			global.Xswlog.Errorf("%v", err.Error())
			return
		}
		err = XswSkipBook(pagebook, keyName, bookInfo.BookName, bookInfo.Author)
		if err != nil {
			global.Xswlog.Errorf("%v", err.Error())
			return
		}
		return
	}
	if isSkip {
		err = XswSkipBook(pagebook, keyName, bookInfo.BookName, bookInfo.Author)
		if err != nil {
			global.Xswlog.Errorf("%v", err.Error())
			return
		}
		return
	}
	return
}

func XswSaveChapterJson(keyName string, pagebook *models.XswPageBooks, bookInfo *models.XswBookDesc, chapters []*models.CollectChapterInfo) (err error) {
	if pagebook == nil || bookInfo == nil {
		global.Xswlog.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	books := pagebook.Books
	bookName := bookInfo.BookName
	author := bookInfo.Author
	var chapterAll []*models.McBookChapter
	var id int64
	var sort int
	for _, val := range chapters {
		id++
		sort++
		chapterName := val.ChapterTitle
		chapterLink := val.ChapterLink
		var textNum, isLess int
		textNum, err = XswCollectChapterText(bookName, author, chapterName, chapterLink)
		if err != nil {
			global.Xswlog.Errorf("获取章节内容失败 bookName=%v author=%v chapterName=%v chapterLink=%v err=%v", bookName, author, chapterName, chapterLink, err.Error())
			continue
		}
		if textNum <= 1000 {
			isLess = 1
		}
		if textNum <= 100 {
			continue
		}
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
	for _, val := range books {
		if val.BookName == bookName && val.Author == author {
			val.Use = 1
			break
		}
	}
	pagebook.Books = books
	err = redis_service.Set(keyName, pagebook, 0)
	if err != nil {
		global.Xswlog.Errorf("缓存采集状态失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}

	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Xswlog.Errorf("获取JSONQ对象失败 %v", err.Error())
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
		global.Xswlog.Errorf("%v", err.Error())
		return
	}
	return
}

func XswGetCacheCategoryBookList(categoryKey, categoryHref string) (pagebook *models.XswPageBooks, keyName string, err error) {
	keyName = fmt.Sprintf("%v_%v", utils.XswBooks, categoryKey)
	pagebookVal := redis_service.Get(keyName)
	if pagebookVal == "" || pagebookVal == "null" {
		pagebook, err = XswGetPageBookList(categoryHref)
		if err != nil {
			return
		}
		err = redis_service.Set(keyName, pagebook, time.Hour*24*30)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
	} else {
		err = json.Unmarshal([]byte(pagebookVal), &pagebook)
		if err != nil {
			global.Xswlog.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
		book := pagebook.Books
		nextPageLink := pagebook.NextPageLink
		for _, val := range book {
			if val.Use == 0 {
				return
			}
		}
		pagebook, err = XswGetPageBookList(nextPageLink)
		if err != nil {
			return
		}
		err = redis_service.Set(keyName, pagebook, time.Hour*24*30)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
	}
	return
}

func XswGetCacheCategory(siteUrl, defaultCategoryName string) (categoryKey, categoryHref string) {
	categoryVal := redis_service.Get(utils.XswCategory)
	var categorys []*models.XswCategory
	var err error
	if categoryVal == "" || categoryVal == "null" {
		categorys = XswGetCategory(siteUrl)
		err = redis_service.Set(utils.XswCategory, categorys, time.Hour*24*30)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
	} else {
		err = json.Unmarshal([]byte(categoryVal), &categorys)
		if err != nil {
			global.Xswlog.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
	}
	for _, val := range categorys {
		if val.CategoryName == defaultCategoryName {
			categoryKey = val.CategoryKey
			categoryHref = val.CategoryHref
			return
		}
	}
	return
}

func XswGetCategory(siteUrl string) (categorys []*models.XswCategory) {
	var html string
	var err error
	html, err = utils.GetHtmlcolly(siteUrl)
	if err != nil {
		global.Xswlog.Errorf("XswGetCategory err:%v", err.Error())
		return
	}
	html = utils.GetSimpleHtml(html)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		global.Xswlog.Errorf("GetTotalPage err:%v", err.Error())
		return
	}
	domain := utils.GetUrlDomain(siteUrl)
	document.Find("#nav li").Each(func(i int, s *goquery.Selection) {
		categoryName := s.Find("a").Text()
		categoryHref, _ := s.Find("a").Attr("href")
		categoryKey, _ := s.Attr("id")
		if categoryName == "" || categoryKey == "" {
			return
		}
		if categoryHref != "" {
			categoryHref = fmt.Sprintf("%v%v", domain, categoryHref)
		}
		category := &models.XswCategory{
			CategoryName: categoryName,
			CategoryKey:  categoryKey,
			CategoryHref: categoryHref,
		}
		categorys = append(categorys, category)
		return
	})
	return
}

func XswGetTotalPage(categoryHref string) (totalPage int) {
	var err error
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(categoryHref)
	if err != nil {
		return
	}
	if tempHtml == "" {
		global.Xswlog.Errorln("GetTotalPage html为空")
		return
	}
	tempHtml = utils.GetSimpleHtml(tempHtml)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		global.Xswlog.Errorf("GetTotalPage err:%v", err.Error())
		return
	}

	pageString := document.Find(".articlepage em").Text()
	pageInfo := strings.Split(pageString, "/")
	if len(pageInfo) > 1 {
		totalPageStr := utils.GetUrlBookNum(pageInfo[1])
		// 转换为整数
		totalPage, err = strconv.Atoi(totalPageStr)
		if err != nil {
			global.Xswlog.Errorf("无法解析页数:%v", err.Error())
			return
		}
	}
	return
}

func XswGetPageBookList(pageLink string) (pagebook *models.XswPageBooks, err error) {
	if pageLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(pageLink)
	if html == "" {
		global.Xswlog.Errorf("GetHtmlcolly为空 %v", err.Error())
		return
	}
	html = utils.GetSimpleHtml(html)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	domain := utils.GetUrlDomain(pageLink)
	var books []*models.XswPageBook

	nextPageLink, _ := document.Find("#main .articlepage .next").Attr("href")
	if nextPageLink != "" {
		nextPageLink = fmt.Sprintf("%v%v", domain, nextPageLink)
	}
	pageString := document.Find(".articlepage em").Text()
	document.Find("#alist #alistbox").Each(func(i int, s *goquery.Selection) {
		pic, _ := s.Find("img").Attr("src")
		bookLink, _ := s.Find("h2 a").Attr("href")
		bookName := s.Find("h2 a").Text()
		author := s.Find("span").Text()
		author = strings.TrimLeft(author, "作者：")
		//desc := s.Find(".intro").Text()
		if bookLink != "" {
			bookLink = fmt.Sprintf("%v%v", domain, bookLink)
		}
		book := &models.XswPageBook{
			BookName: bookName,
			Author:   author,
			Pic:      pic,
			Link:     bookLink,
		}
		books = append(books, book)
		return
	})
	pagebook = &models.XswPageBooks{
		PageLink:     pageLink,
		NextPageLink: nextPageLink,
		TotalPage:    pageString,
		Books:        books,
	}
	return
}

func XswGetChapters(bookLink string) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Xswlog.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if html == "" {
		global.Xswlog.Errorf("GetHtmlcolly为空 %v", bookLink)
		return
	}
	html = utils.GetSimpleHtml(html)
	domain := utils.GetUrlDomain(bookLink)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var chaptersAll []*models.CollectChapterInfo
	document.Find(".liebiao ul li").Each(func(i int, s *goquery.Selection) {
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

func XswCollectChapterText(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Biquge34log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := XswGetChapterText(chapterLink)
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

func XswGetChapterText(chapterLink string) (text string) {
	var tempHtml string
	var err error
	tempHtml, err = utils.GetHtmlcolly(chapterLink)
	if err != nil {
		global.Xswlog.Errorf("GetHtmlcolly为空 %v %v", chapterLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Xswlog.Errorf("GetHtmlcolly为空 %v", chapterLink)
		return
	}
	tempHtml = utils.GetSimpleHtml(tempHtml)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	num1, num2 := XswGetChapterLinkNum(chapterLink)
	tempText, _ := document.Find("#content").Html()
	tempText = html.UnescapeString(tempText)
	tempText = strings.ReplaceAll(tempText, "<p>", "")
	tempText = strings.ReplaceAll(tempText, "</p>", "")
	tempText = strings.ReplaceAll(tempText, "<center>", "")
	tempText = strings.ReplaceAll(tempText, "</center>", "")
	tempText = strings.ReplaceAll(tempText, fmt.Sprintf("javascript:baocuo('%v','%v');", num1, num2), "")
	tempText = strings.ReplaceAll(tempText, "style=\"color:red\"", "")
	tempText = strings.ReplaceAll(tempText, ">>章节报错<<", "")
	tempText = strings.ReplaceAll(tempText, "<br/>", "\r\n")
	tempText = strings.ReplaceAll(tempText, "<a href=\"\" ></a>", "")
	tempText = strings.ReplaceAll(tempText, "精华书阁", "")
	tempText = strings.ReplaceAll(tempText, "该章节缺失txt文件", "")
	tempText = strings.Replace(tempText, "精华书阁", "", 1)
	tempText = strings.Replace(tempText, "，", "", 1)
	text = strings.TrimLeft(tempText, "，")
	//text = strings.TrimSpace(tempText)
	return
}

func GetBookDesc(bookDetailLink string) (info *models.XswBookDesc, err error) {
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(bookDetailLink)
	if err != nil {
		global.Xswlog.Errorf("%v %v", bookDetailLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Xswlog.Errorf("GetHtmlcolly为空 %v %v", bookDetailLink, err.Error())
		return
	}
	tempHtml = utils.GetSimpleHtml(tempHtml)
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	domain := utils.GetUrlDomain(bookDetailLink)
	sourceUrl, _ := document.Find(".option .btopt a").Attr("href")
	if sourceUrl != "" {
		sourceUrl = fmt.Sprintf("%v%v", domain, sourceUrl)
	}
	bookinfo := document.Find(".box_info tr").First().Find("td").First()
	bookName := bookinfo.Find("h1.f20h").Contents().Not("em").Text()
	author := bookinfo.Find("h1.f20h em").Text()
	if author != "" {
		author = strings.ReplaceAll(author, "作者：", "")
	}

	pic, _ := document.Find(".box_intro .pic img").Attr("src")
	filePath, _ := utils.DownImg(bookName, author, pic, "/data/pic/")
	pic = strings.TrimLeft(filePath, ".")

	tr := document.Find(".box_info tr").Eq(4)
	categoryName := tr.Find("td").First().Text()
	if categoryName != "" {
		categoryName = strings.ReplaceAll(categoryName, "小说分类：", "")
	}
	serializeName := tr.Find("td").Eq(2).Text()
	if serializeName != "" {
		serializeName = strings.ReplaceAll(serializeName, "小说状态：", "")
	}

	lastChapterTitle := document.Find(".book_newchap .con .ti").First().Find("a").Text()

	desc := document.Find(".intro").Text()
	tr = document.Find(".box_info tr").Last()
	lastChapterTimeStr := tr.Find("td").Last().Text()
	lastChapterTimeStr = strings.ReplaceAll(lastChapterTimeStr, "更新时间：", "")
	lastChapterTimeStr = strings.TrimSpace(strings.ReplaceAll(lastChapterTimeStr, "0:00:00", ""))
	lastChapterTime := XswDateToUnix(lastChapterTimeStr)

	categorys := []*models.CategoryReg{
		{"玄幻魔法", 12},
		{"武侠修真", 87},
		{"都市言情", 45},
		{"历史军事", 23},
		{"游戏动漫", 81},
		{"恐怖灵异", 24},
		{"其他类型", 56},
		{"同人小说", 56},
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
	if strings.Contains(serializeName, "连载") {
		serialize = 1
	} else if strings.Contains(serializeName, "完本") {
		serialize = 2
	}
	if bookName == "" || author == "" {
		global.Xswlog.Errorf("小说名称或作者为空 %v %v", bookName, author)
		return
	}
	info = &models.XswBookDesc{
		BookName:         bookName,
		Author:           author,
		Pic:              pic,
		Desc:             desc,
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

func XswDateToUnix(date string) (timestamp int64) {
	layout := "2006/1/2"
	t, err := time.Parse(layout, date)
	if err != nil {
		return
	}
	timestamp = t.Unix()
	return
}

func XswGetChapterLinkNum(chapterLink string) (num1, num2 int) {
	// 使用正则表达式提取数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(chapterLink, -1)
	if len(matches) > 0 {
		num1, _ = strconv.Atoi(matches[0])
	}
	if len(matches) > 1 {
		num2, _ = strconv.Atoi(matches[1])
	}
	return
}
