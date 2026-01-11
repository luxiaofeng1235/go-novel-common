package main

import (
	"go-novel/app/service/collect_service"
	"go-novel/db"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	host, name, user, passwd := db.GetDB()
	db.InitDB(host, name, user, passwd)

	addr, passwd, defaultdb := db.GetRedis()
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	//redis_service.Set("test", "test", time.Second*100)
	//collect_service.CollectContinuation(1)
	//collect_service.Collect(1, 1)
	collect_service.CollectThread(1)
	//collect_service.CollectChapterThread()

	//collect_service.Alone(1, "http://www.biquw.la/book/64442/")
	//collect_service.Alone(1, "http://www.biquw.la/book/6085/")
	//collect_service.Alone(1, "http://www.biquw.la/book/139802/")

	//getSource()

}

//func getSource() (err error) {
//	var sources []*models.McCollect
//	sources, err = collect_service.GetCollectList()
//	if err != nil {
//		return
//	}
//	if len(sources) <= 0 {
//		return
//	}
//	for _, source := range sources {
//		collectListSet(source, 0)
//	}
//	return
//}
//
//func collectListSet(source *models.McCollect, current int) (err error) {
//	patternListSection := source.ListSectionReg
//	if patternListSection == "" {
//		global.Collectlog.Errorf("列表区间正则不能为空 sourceId=%v", source.Id)
//		return
//	}
//
//	patternListUrl := source.ListUrlReg
//	if patternListUrl == "" {
//		global.Collectlog.Errorf("获取小说详情链接正则不能为空 sourceId=%v", source.Id)
//		return
//	}
//	urls, err := utils.GetListSourceURL(source.ListPageReg)
//	if err != nil {
//		return
//	}
//	//log.Println(urls)
//	//return
//	var html string
//	html, err = getHtml(urls[current], source.Charset, source.UrlComplete)
//	re := regexp.MustCompile(patternListSection)
//	match := re.FindStringSubmatch(html)
//	if len(match) <= 0 {
//		global.Collectlog.Errorf("采集获取列表失败 %v", err.Error())
//		return
//	}
//	linkPattern := source.ListUrlReg
//	linksRe := regexp.MustCompile(linkPattern)
//	links := linksRe.FindAllStringSubmatch(match[0], -1)
//	var bookLinks []string
//	for _, link := range links {
//		if len(link) > 1 {
//			href := link[1]
//			bookLinks = append(bookLinks, href)
//		}
//	}
//	if source.UrlReverse > 0 {
//		bookLinks = utils.ArrayReverse(bookLinks)
//	}
//	if len(bookLinks) <= 0 {
//		return
//	}
//	for _, link := range bookLinks {
//		getBookDetail(source, link)
//		return
//	}
//	return
//}
//
//func getHtml(url, encode string, urlComplete int) (html string, err error) {
//	html, err = utils.DoGet(url)
//	if urlComplete > 0 {
//		html = utils.UrlComplete(html, url)
//	}
//	html = utils.AutoConvertToUTF8(html, encode)
//	return
//}
//
//func getBookDetail(source *models.McCollect, bookUrl string) (err error) {
//	var html string
//	html, err = getHtml(bookUrl, source.Charset, source.UrlComplete)
//	if html == "" {
//		global.Collectlog.Errorf("获取小说详情页面失败 bookUrl=%v", bookUrl)
//		return
//	}
//	//log.Println(html)
//
//	var categoryName, bookName, author, desc, serialize, tag, pic string
//	matchCate := regexp.MustCompile(source.CategoryNameReg).FindStringSubmatch(html)
//	if len(matchCate) > 0 {
//		categoryName = matchCate[1]
//	} else {
//		global.Collectlog.Errorf("获取小说分类失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchBookName := regexp.MustCompile(source.BookNameReg).FindStringSubmatch(html)
//	if len(matchBookName) > 0 {
//		bookName = matchBookName[1]
//	} else {
//		global.Collectlog.Errorf("获取小说名称失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchAuthor := regexp.MustCompile(source.AuthorReg).FindStringSubmatch(html)
//	if len(matchAuthor) > 0 {
//		author = matchAuthor[1]
//	} else {
//		global.Collectlog.Errorf("获取小说作者失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchPic := regexp.MustCompile(source.PicReg).FindStringSubmatch(html)
//	if len(matchPic) > 0 {
//		pic = matchPic[1]
//	} else {
//		global.Collectlog.Errorf("获取小说图片失败 %v bookUrl=%v", bookUrl)
//	}
//
//	if source.PicLocal > 0 {
//		var uploadBookPath string
//		uploadBookPath, err = setting_service.GetValueByName("uploadBookPath")
//		if err != nil {
//			global.Collectlog.Errorf("获取小说上传目录失败 %v", err.Error())
//			return
//		}
//		var picPath string
//		//http://www.biquw.la/files/article/image/4/4046/4046s.jpg
//		picPath, err = utils.DownImg(bookName, pic, uploadBookPath)
//		if err != nil {
//			global.Collectlog.Errorf("下载小说图片失败 %v", err.Error())
//			return
//		}
//		log.Println(picPath, err)
//	}
//
//	matchDesc := regexp.MustCompile(source.DescReg).FindStringSubmatch(html)
//	if len(matchDesc) > 0 {
//		desc = matchDesc[1]
//	} else {
//		global.Collectlog.Errorf("获取小说简介失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchSerialize := regexp.MustCompile(source.SerializeReg).FindStringSubmatch(html)
//	if len(matchSerialize) > 0 {
//		serialize = matchSerialize[1]
//	} else {
//		global.Collectlog.Errorf("获取小说简介失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchTagName := regexp.MustCompile(source.TagName).FindStringSubmatch(html)
//	if len(matchTagName) > 0 {
//		tag = matchTagName[1]
//	} else {
//		global.Collectlog.Errorf("获取小说标签失败 %v bookUrl=%v", bookUrl)
//	}
//
//	matchChapterSection := regexp.MustCompile(source.ChapterSectionReg).FindStringSubmatch(html)
//	if len(matchTagName) > 0 {
//		ulContent := matchChapterSection[0]
//		liMatches := regexp.MustCompile(source.ChapterUrlReg).FindAllStringSubmatch(ulContent, -1)
//		for _, liMatch := range liMatches {
//			if len(liMatch) > 2 {
//				href := liMatch[1]
//				chapterName := liMatch[2]
//				fmt.Printf("Href: %s\nChapter Name: %s\n\n", href, chapterName)
//			}
//		}
//	} else {
//		global.Collectlog.Errorf("获取小说章节区间失败 %v bookUrl=%v", bookUrl)
//	}
//	log.Println("categoryName", categoryName)
//	log.Println("bookName", bookName)
//	log.Println("author", author)
//	log.Println("desc", desc)
//	log.Println("serialize", serialize)
//	log.Println("tag", tag)
//	log.Println("pic", pic)
//	return
//}
