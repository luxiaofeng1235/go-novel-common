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
	"html"
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
	var siteUrl string = "https://www.siluke520.net"
	for {
		Siluke520Collect(siteUrl)
	}
}

func Siluke520Collect(siteUrl string) {
	categorys := Siluke520GetCacheCategory(siteUrl)
	if len(categorys) <= 0 {
		return
	}
	var category = new(models.Siluke520Category)
	for _, val := range categorys {
		if val.Use == 0 {
			category = val
			break
		}
	}
	Siluke520GetBooksByCategory(category)
	return
}

func Siluke520GetCacheCategory(siteUrl string) (categorys []*models.Siluke520Category) {
	categoryVal := redis_service.Get(utils.Siluke520Category)
	var err error
	if categoryVal == "" || categoryVal == "null" {
		categorys = Siluke520GetCategory(siteUrl)
		err = redis_service.Set(utils.Siluke520Category, categorys, 0)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
		return
	} else {
		err = json.Unmarshal([]byte(categoryVal), &categorys)
		if err != nil {
			global.Siluke520log.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func Siluke520GetCategory(siteUrl string) (categorys []*models.Siluke520Category) {
	var html string
	var err error
	html, err = utils.GetHtmlcolly(siteUrl)
	if err != nil {
		global.Siluke520log.Errorf("Siluke520GetCategory err:%v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		global.Siluke520log.Errorf("err:%v", err.Error())
		return
	}
	domain := utils.GetUrlDomain(siteUrl)
	document.Find(".nav_cont ul li").Each(func(i int, s *goquery.Selection) {
		categoryName := s.Find("a").Text()
		categoryHref, _ := s.Find("a").Attr("href")
		if categoryName == "" {
			return
		}
		if categoryName == "首页" || categoryName == "热门小说" {
			return
		}
		categoryKey := strings.Join(pinyin.LazyPinyin(categoryName, pinyin.NewArgs()), "")
		if categoryHref != "" {
			categoryHref = fmt.Sprintf("%v%v", domain, categoryHref)
		}
		category := &models.Siluke520Category{
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

func Siluke520GetBooksByCategory(category *models.Siluke520Category) {
	if category.CategoryHref == "" {
		err := redis_service.Del(utils.Siluke520Category)
		if err != nil {
			global.Siluke520log.Errorf("%v", err.Error())
			return
		}
	}
	pagebook, pageBookKey, err := Siluke520GetCacheCategoryBookList(category.CategoryKey, category.CategoryHref)
	if err != nil {
		global.Siluke520log.Errorf("%v", err.Error())
		return
	}
	if pagebook == nil {
		return
	}
	books := pagebook.Books
	var bookInfo *models.Siluke520BookDesc
	var isSkip bool
	bookInfo, isSkip, err = Siluke520GetBookInfo(books)
	if isSkip {
		err = Siluke520RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Siluke520log.Errorf("%v", err.Error())
			return
		}
	}
	if bookInfo == nil {
		global.Siluke520log.Errorln("bookInfo为nil")
		return
	}
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			err = Siluke520RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
			if err != nil {
				global.Siluke520log.Errorf("%v", err.Error())
				return
			}
		}
	}
	chapters, isSkip, err := Siluke520GetFilterChapter(pageBookKey, pagebook, bookInfo)
	if err != nil {
		global.Siluke520log.Errorf("Siluke520GetFilterChapter err=%v", err.Error())
		return
	}
	if isSkip {
		return
	}
	err = Siluke520SaveChapterJson(pageBookKey, pagebook, bookInfo, chapters)
	if err != nil {
		return
	}
}

func Siluke520GetCacheCategoryBookList(categoryKey, categoryHref string) (pagebook *models.Siluke520PageBooks, pageBookKey string, err error) {
	pageBookKey = fmt.Sprintf("%v_%v", utils.Siluke520Books, categoryKey)
	pagebookVal := redis_service.Get(pageBookKey)
	if pagebookVal == "" || pagebookVal == "null" {
		pagebook, err = Siluke520GetPageBookList(categoryHref)
		if err != nil {
			global.Siluke520log.Errorf("%v", err.Error())
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
			global.Siluke520log.Errorf("获取分类缓存失败 err=%v", err.Error())
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
			pagebook, err = Siluke520GetPageBookList(nextPageLink)
			if err != nil {
				return
			}
			err = redis_service.Set(pageBookKey, pagebook, 0)
			if err != nil {
				err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
				return
			}
			if len(pagebook.Books) <= 0 {
				err = Siluke520RemoveCategory(pageBookKey, categoryHref)
				if err != nil {
					return
				}
			}
		} else {
			err = Siluke520RemoveCategory(pageBookKey, categoryHref)
			if err != nil {
				return
			}
		}
	}
	return
}

func Siluke520RemoveCategory(pageBookKey, categoryHref string) (err error) {
	err = redis_service.Del(pageBookKey)
	if err != nil {
		global.Siluke520log.Errorf("删除pageBookKey=%v 出错 %v", pageBookKey, err.Error())
		return
	}
	categoryVal := redis_service.Get(utils.Siluke520Category)
	var categorys []*models.Siluke520Category
	err = json.Unmarshal([]byte(categoryVal), &categorys)
	if err != nil {
		global.Siluke520log.Errorf("获取分类缓存失败 err=%v", err.Error())
		return
	}
	for _, category := range categorys {
		if category.CategoryHref == categoryHref {
			category.Use = 1
			break
		}
	}
	err = redis_service.Set(utils.Siluke520Category, categorys, 0)
	if err != nil {
		err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
		return
	}
	return
}

func Siluke520RemoveBook(pagebook *models.Siluke520PageBooks, pageBookKey, bookName, author, sourceUrl string, isSkip bool) (err error) {
	books := pagebook.Books
	var isDel bool = true
	for _, val := range books {
		if val.Use == 0 {
			isDel = false
		}
		if val.BookLink == sourceUrl {
			val.Use = 1
			if isSkip {
				val.Use = 2
			}
			break
		}
	}
	if isDel {
		return
	}
	pagebook.Books = books
	err = redis_service.Set(pageBookKey, pagebook, 0)
	if err != nil {
		global.Siluke520log.Errorf("Siluke520RemoveBook bookName=%v author=%v sourceUrl=%v err=%v", bookName, author, sourceUrl, err.Error())
		return
	}
	return
}

func Siluke520GetBookInfo(books []*models.Siluke520PageBook) (bookInfo *models.Siluke520BookDesc, isSkip bool, err error) {
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
	bookInfo, isSkip, err = Siluke520GetBookDesc(bookLink)
	if bookInfo == nil {
		bookInfo = new(models.Siluke520BookDesc)
		bookInfo.BookName = bookName
		bookInfo.Author = author
		bookInfo.SourceUrl = bookLink
	}
	return
}

func Siluke520GetFilterChapter(pageBookKey string, pagebook *models.Siluke520PageBooks, bookInfo *models.Siluke520BookDesc) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if pagebook == nil || bookInfo == nil {
		global.Siluke520log.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	sourceUrl := bookInfo.SourceUrl
	chapters, isSkip, err = Siluke520GetChapters(sourceUrl)
	if err != nil {
		global.Siluke520log.Errorf("%v", err.Error())
		if bookInfo == nil {
			global.Siluke520log.Errorf("%v", err.Error())
			return
		}
		err = Siluke520RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Siluke520log.Errorf("%v", err.Error())
			return
		}
		return
	}
	if isSkip {
		err = Siluke520RemoveBook(pagebook, pageBookKey, bookInfo.BookName, bookInfo.Author, bookInfo.SourceUrl, true)
		if err != nil {
			global.Siluke520log.Errorf("%v", err.Error())
			return
		}
		return
	}
	return
}

func Siluke520SaveChapterJson(pageBookKey string, pagebook *models.Siluke520PageBooks, bookInfo *models.Siluke520BookDesc, chapters []*models.CollectChapterInfo) (err error) {
	if pagebook == nil || bookInfo == nil {
		global.Siluke520log.Errorln("pagebook 或 bookInfo 为nil")
		return
	}
	bookName := bookInfo.BookName
	author := bookInfo.Author
	sourceUrl := bookInfo.SourceUrl
	Siluke520ManyThreadChapter(bookName, author, chapters)
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
		textNum, err = Siluke520CollectChapterText(bookName, author, chapterName, chapterLink)
		if err != nil {
			global.Siluke520log.Errorf("获取章节内容失败 bookName=%v author=%v chapterName=%v chapterLink=%v err=%v", bookName, author, chapterName, chapterLink, err.Error())
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

	err = Siluke520RemoveBook(pagebook, pageBookKey, bookName, author, sourceUrl, false)
	if err != nil {
		return
	}

	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Siluke520log.Errorf("获取JSONQ对象失败 %v", err.Error())
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
		global.Siluke520log.Errorf("%v", err.Error())
		return
	}
	return
}

func Siluke520ManyThreadChapter(bookName, author string, chapters []*models.CollectChapterInfo) {
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
			_, err = Siluke520CollectChapterTextThread(bookName, author, chapterName, chapterLink)
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

func Siluke520GetTotalPage(categoryHref string) (totalPage int) {
	var err error
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(categoryHref)
	if err != nil {
		return
	}
	if tempHtml == "" {
		global.Siluke520log.Errorln("GetTotalPage html为空")
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		global.Siluke520log.Errorf("GetTotalPage err:%v", err.Error())
		return
	}

	pageString := document.Find(".articlepage em").Text()
	pageInfo := strings.Split(pageString, "/")
	if len(pageInfo) > 1 {
		totalPageStr := utils.GetUrlBookNum(pageInfo[1])
		// 转换为整数
		totalPage, err = strconv.Atoi(totalPageStr)
		if err != nil {
			global.Siluke520log.Errorf("无法解析页数:%v", err.Error())
			return
		}
	}
	return
}

func Siluke520GetPageBookList(pageLink string) (pagebook *models.Siluke520PageBooks, err error) {
	if pageLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(pageLink)
	if html == "" {
		global.Siluke520log.Errorf("GetHtmlcolly为空 %v", err.Error())
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var books []*models.Siluke520PageBook
	cateNum, pageNum := Siluke520GetPageNum(pageLink)
	nextLink, _ := document.Find("#pagelink .next").Attr("href")
	if nextLink != "" {
		if !strings.Contains(nextLink, "http") {
			domain := utils.GetUrlDomain(pageLink)
			nextLink = fmt.Sprintf("%v%v", domain, nextLink)
		}
	}
	totalPage := document.Find("#pagelink .last").Text()
	totalPageInt, _ := strconv.Atoi(totalPage)
	document.Find(".news ul li").Each(func(i int, s *goquery.Selection) {
		bookLink, _ := s.Find(".s2 a").Attr("href")
		bookName := s.Find(".s2 a").Text()
		author := s.Find(".s4").Text()
		if bookName == "" || author == "" {
			return
		}
		book := &models.Siluke520PageBook{
			BookName: bookName,
			Author:   author,
			BookLink: bookLink,
		}
		books = append(books, book)
		return
	})
	pagebook = &models.Siluke520PageBooks{
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

func Siluke520GetChapters(bookLink string) (chapters []*models.CollectChapterInfo, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var html string
	html, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Siluke520log.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if html == "" {
		global.Siluke520log.Errorf("GetHtmlcolly为空 %v", bookLink)
		return
	}
	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return
	}
	var chaptersAll []*models.CollectChapterInfo
	document.Find(".book_list ul").Last().Find("li").Each(func(i int, s *goquery.Selection) {
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

func Siluke520CollectChapterText(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Siluke520log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := Siluke520GetChapterText(bookName, author, chapterName, chapterLink)
	textNum = len([]rune(text))
	if textNum <= 10 {
		err = fmt.Errorf("获取内容失败 %v", text)
		return
	}
	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func Siluke520CollectChapterTextThread(bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Siluke520log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text := Siluke520GetChapterText(bookName, author, chapterName, chapterLink)
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

func Siluke520GetChapterText(bookName, author, chapterName, chapterLink string) (text string) {
	var tempHtml string
	var err error
	tempHtml, err = utils.GetHtmlcolly(chapterLink)
	if err != nil {
		global.Siluke520log.Errorf("GetHtmlcolly为空 %v %v", chapterLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Siluke520log.Errorf("GetHtmlcolly为空 %v", chapterLink)
		return
	}

	var document *goquery.Document
	document, err = goquery.NewDocumentFromReader(strings.NewReader(tempHtml))
	if err != nil {
		return
	}
	tempText, _ := document.Find("#htmlContent").Html()
	tempText = html.UnescapeString(tempText)
	bookLink, _ := document.Find(".srcbox a").Last().Attr("href")
	tempText = strings.ReplaceAll(tempText, "<br/><br/>", "<br/>")
	tempText = strings.ReplaceAll(tempText, "<br/>", "\r\n")
	tempText = strings.ReplaceAll(tempText, fmt.Sprintf("<a href=\"%v\">%v</a>", bookLink, bookName), "")
	tempText = strings.ReplaceAll(tempText, bookName, "")
	tempText = strings.ReplaceAll(tempText, author, "")
	tempText = strings.ReplaceAll(tempText, chapterName, "")
	tempText = strings.ReplaceAll(tempText, "思路客小说网", "")
	tempText = strings.ReplaceAll(tempText, "www.siluke520.net", "")
	tempText = strings.ReplaceAll(tempText, "，", "")
	tempText = strings.ReplaceAll(tempText, "·", "")
	tempText = strings.ReplaceAll(tempText, "最快更新", "")
	tempText = strings.ReplaceAll(tempText, "最新章节！", "")
	tempText = strings.ReplaceAll(tempText, "本书由红薯网授权掌阅科技电子版制作与发行", "")
	tempText = strings.ReplaceAll(tempText, "版权所有", "")
	tempText = strings.ReplaceAll(tempText, "侵权必究", "")
	tempText = strings.ReplaceAll(tempText, "<table class=\"zhangyue-tablebody\">", "")
	tempText = strings.ReplaceAll(tempText, "<tbody>", "")
	tempText = strings.ReplaceAll(tempText, "<tr style=\"height: 78%;vertical-align: middle;\">", "")
	tempText = strings.ReplaceAll(tempText, "<td class=\"biaoti\">", "")
	tempText = strings.ReplaceAll(tempText, "<span class=\"kaiti\">", "")
	tempText = strings.ReplaceAll(tempText, "</span>", "")
	tempText = strings.ReplaceAll(tempText, "</td>", "")
	tempText = strings.ReplaceAll(tempText, "</tr>", "")
	tempText = strings.ReplaceAll(tempText, "<tr style=\"height: 17%;vertical-align: bottom;\">", "")
	tempText = strings.ReplaceAll(tempText, "<td class=\"copyright\">", "")
	tempText = strings.ReplaceAll(tempText, "<span class=\"lantinghei\">", "")
	tempText = strings.ReplaceAll(tempText, "</span>", "")
	tempText = strings.ReplaceAll(tempText, "<span class=\"dotStyle2\">", "")
	tempText = strings.ReplaceAll(tempText, "<span class=\"lantinghei\">", "")
	tempText = strings.ReplaceAll(tempText, "</tbody>", "")
	tempText = strings.ReplaceAll(tempText, "</table>", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;table class=&#34;zhangyue-tablebody&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;tbody&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;tr style=&#34;height: 78%;vertical-align: middle;&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;td class=&#34;biaoti&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;span class=&#34;kaiti&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/span&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/td&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/tr&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;tr style=&#34;height: 17%;vertical-align: bottom;&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;td class=&#34;copyright&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;span class=&#34;lantinghei&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/span&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;span class=&#34;dotStyle2&#34;&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/tbody&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/table&gt;", "")
	//tempText = strings.ReplaceAll(tempText, "&lt;/span&gt;", "")
	text = strings.TrimSpace(tempText)
	return
}

func Siluke520GetBookDesc(bookLink string) (info *models.Siluke520BookDesc, isSkip bool, err error) {
	if bookLink == "" {
		return
	}
	var tempHtml string
	tempHtml, err = utils.GetHtmlcolly(bookLink)
	if err != nil {
		global.Siluke520log.Errorf("%v %v", bookLink, err.Error())
		return
	}
	if tempHtml == "" {
		global.Siluke520log.Errorf("GetHtmlcolly为空 %v %v", bookLink, err.Error())
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
	pic, _ := document.Find("meta[property='og:image']").Attr("content")
	//修改存储的图片路径
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
		global.Siluke520log.Errorf("小说名称或作者为空 %v %v %v", bookName, author, bookLink)
		isSkip = true
		return
	}
	info = &models.Siluke520BookDesc{
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

func Siluke520GetPageNum(pageUrl string) (cateNum, pageNum int) {
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
