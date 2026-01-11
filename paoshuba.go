package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"log"
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
	//bookUrl := "http://www.paoshu8.info/27_27766/"
	//bookName := "农门医妻"
	//author := "林十五"
	//PaoshubaCollectBook(bookUrl, bookName, author)
	var books []*models.McBook
	global.DB.Model(models.McBook{}).Debug().Order("id asc").Where("is_less = 1 and source_url like ?", "%"+"paoshu8"+"%").Find(&books)
	for _, book := range books {
		PaoshubaCollectBook(book.SourceUrl, book.BookName, book.Author)
	}
	//for {
	//
	//}
}

func PaoshubaCollectBook(bookUrl, bookName, author string) {
	chapters, err := PaoshubaGetCacheChapters(bookUrl, bookName, author)
	if err != nil {
		global.Paoshu8log.Errorf("%v", err.Error())
		return
	}
	for _, val := range chapters {
		chapterName := val.ChapterTitle
		chapterLink := val.ChapterLink
		var textNum int
		textNum, err = PaoshubaCollectChapterText(bookUrl, bookName, author, chapterName, chapterLink)
		if err != nil && textNum > 0 {
			global.Paoshu8log.Errorf("获取章节内容失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
			continue
		}
		updatedChapter := &models.McBookChapter{
			ChapterLink: chapterLink,
			ChapterName: chapterName,
			Vip:         0,
			Cion:        0,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
		}
		err = chapter_service.CreateChapter(bookName, author, updatedChapter)
		if err != nil {
			global.Paoshu8log.Errorf("%v", err.Error())
			return
		}
	}
}
func PaoshubaGetCacheChapters(bookUrl, bookName, author string) (chapters []*models.CollectChapterInfo, err error) {
	bookNum := utils.GetUrlBookNum(bookUrl)
	bookUrlKey := fmt.Sprintf("%v_%v", utils.PaoshubaChapters, bookNum)
	chaptersVal := redis_service.Get(bookUrlKey)
	if chaptersVal == "" || chaptersVal == "null" {
		chapters, err = PaoshubaGetChapters(bookUrl, bookName, author)
		err = redis_service.Set(bookUrlKey, chapters, 0)
		if err != nil {
			err = fmt.Errorf("缓存泡书吧章节列表失败 err=%v", err.Error())
			return
		}
		return
	} else {
		err = json.Unmarshal([]byte(chaptersVal), &chapters)
		if err != nil {
			global.Paoshu8log.Errorf("获取泡书吧章节列表缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func PaoshubaGetChapters(bookUrl, bookName, author string) (chapters []*models.CollectChapterInfo, err error) {
	var html string
	html, err = utils.GetHtmlcolly(bookUrl)
	log.Println(html)
	if html == "" {
		err = fmt.Errorf("获取泡书吧小说详情页面失败 bookUrl=%v bookName=%v author=%v", bookUrl, bookName, author)
		return
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		err = fmt.Errorf("goquery err=%v", err.Error())
		return
	}
	// 选择包含章节信息的 dl 元素
	dl := document.Find("#list dl")
	// 查找包含 "正文" 字样的 dt 标签
	dt := dl.Find("dt:contains('正文')")
	// 获取该 dt 标签下面的所有 dd 标签
	dd := dt.NextAllFiltered("dd")
	// 遍历 dd 标签并获取链接和标题信息
	dd.Each(func(i int, selection *goquery.Selection) {
		link := selection.Find("a").AttrOr("href", "")
		link = utils.GetUrlSuffix(link)
		if link != "" {
			if strings.HasPrefix(link, "/") {
				link = fmt.Sprintf("%v%v", bookUrl, link)
			} else {
				link = fmt.Sprintf("%v/%v", bookUrl, link)
			}
		}
		title := selection.Find("a").Text()
		chapter := &models.CollectChapterInfo{
			ChapterTitle: title,
			ChapterLink:  link,
		}
		chapters = append(chapters, chapter)
	})
	return
}

func PaoshubaCollectChapterText(bookUrl, bookName, author, chapterName, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		if textNum > 10 {
			global.Paoshu8log.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
			return
		}
	}
	text, err := PaoshubaGetChapterText(bookName, chapterName, chapterLink)
	if err != nil {
		global.Paoshu8log.Infof("采集章节内容失败 bookName=%v chapterTitle=%v chapterLink=%v ", bookName, chapterName, chapterLink)
		return
	}
	textNum = len([]rune(strings.TrimSpace(text)))
	if textNum <= 500 {
		err = fmt.Errorf("获取内容失败 link=%v text=%v textNum=%v", chapterLink, text, textNum)
		domain := utils.GetUrlDomain(chapterLink)
		chapterLink = strings.ReplaceAll(chapterLink, domain, "http://www.biquge5200.net")
		text, err = PaoshubaGetChapterText(bookName, chapterName, chapterLink)
		if err != nil {
			global.Paoshu8log.Infof("采集章节内容失败 bookName=%v chapterTitle=%v chapterLink=%v ", bookName, chapterName, chapterLink)
			return
		}
		textNum = len([]rune(strings.TrimSpace(text)))
	}
	//if textNum <= 10 {
	//	err = fmt.Errorf("获取内容失败 link=%v text=%v textNum=%v", chapterLink, text, textNum)
	//	return
	//}
	log.Println("text", bookName, author, chapterName, chapterLink, text, textNum)

	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}

func PaoshubaGetChapterText(bookName, chapterTitle, chapterLink string) (text string, err error) {
	var chapterHtml string
	chapterHtml, err = utils.GetHtmlcolly(chapterLink)
	if err != nil {
		return
	}
	if strings.Contains(chapterHtml, "503 Service Temporarily Unavailable") {
		global.Paoshu8log.Errorf("%v 503 Service Temporarily Unavailable", chapterLink)
		time.Sleep(time.Second)
		return
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(chapterHtml))
	if err != nil {
		err = fmt.Errorf("goquery err=%v", err.Error())
		return
	}
	html, _ := document.Find("#content").Html()
	// 将 <p> 标签替换为换行符 \r\n
	text = strings.ReplaceAll(html, "<p>", "\r\n")
	text = strings.ReplaceAll(text, "</p>", "")
	text = strings.ReplaceAll(text, "【】", "")
	text = strings.ReplaceAll(text, "本文由。。首发", "")
	text = strings.ReplaceAll(text, "樂文小說", "")
	text = strings.ReplaceAll(text, fmt.Sprintf("《<b>%v</b>》", bookName), "")
	text = strings.ReplaceAll(text, fmt.Sprintf("正在手打中，请稍等片刻，内容更新后，请重新刷新页面，即可获取最新更新！", bookName), "")
	text = strings.ReplaceAll(text, fmt.Sprintf("《%v》%v", bookName, chapterTitle), "")
	text = strings.ReplaceAll(text, "try{", "")
	text = strings.ReplaceAll(text, "mad1('gad2');}", "")
	text = strings.ReplaceAll(text, "catch(ex){}", "")
	text = strings.ReplaceAll(text, "mad1(", "")
	text = strings.ReplaceAll(text, "gad2", "")
	text = strings.ReplaceAll(text, "')", "")
	text = strings.ReplaceAll(text, ";}", "")
	text = strings.ReplaceAll(text, "'')", "")
	text = strings.ReplaceAll(text, "56书库新网址：", "")
	text = strings.ReplaceAll(text, "＠樂＠文＠小＠说|", "")
	text = strings.ReplaceAll(text, "泡书吧全文字更新,牢记网址:www.paoshu8.info", "")
	text = strings.ReplaceAll(text, "www.paoshu8.info", "")
	text = strings.ReplaceAll(text, "泡书吧全文字更新,牢记网址:", "")
	text = strings.ReplaceAll(text, "正在手打中，请稍等片刻，内容更新后，请重新刷新页面，即可获取最新更新！", "")
	text = strings.ReplaceAll(text, "支持（狂沙文学网）把本站分享那些需要的小伙伴！找不到书请留言！", "")
	text = strings.ReplaceAll(text, "剩余内容请前往纵横小说继续阅读.百度或各大应用市场搜索“纵横小说”，仙侠玄幻雪中,一剑奇幻武侠,青鸾,土豆脑洞剑道第一仙为生活添点料。或直接访问www.zongheng.com", "")
	text = strings.ReplaceAll(text, "☆★☆★☆", "")
	text = strings.ReplaceAll(text, "<div id=\"gc1\" class=\"gcontent1\"><script type=\"text/javascript\">ggauto() </script></div>", "")
	text = strings.ReplaceAll(text, "www.shukeba", "")
	return
}
