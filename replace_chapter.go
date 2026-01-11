package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
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
	db.InitKeyLock()
	var books []*models.McBook
	//.Where("source_url like ?", "%"+"paoshu8"+"%")
	var start = time.Now()
	global.DB.Model(models.McBook{}).Debug().Order("id asc").Where("source_url like ?", "%"+"paoshu8"+"%").Find(&books)
	for _, book := range books {
		replaceText(book)
	}
	log.Println(time.Since(start))
}
func replaceText(book *models.McBook) {
	log.Println(book.Id, "开始")
	bookId := book.Id
	bookName := book.BookName
	author := book.Author
	//章节表
	var chapterFile string
	var err error
	_, chapterFile, err = chapter_service.GetJsonqByBookName(book.BookName, book.Author)
	if err != nil {
		global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
		return
	}
	chapterAll, _, err := chapter_service.GetChapterNamesByFile(chapterFile)
	if err != nil {
		global.Collectlog.Errorf("%v", err.Error())
		return
	}
	var isLess bool
	for _, chapter := range chapterAll {
		var text string
		var chapterName = chapter.ChapterName
		var chapterLink = chapter.ChapterLink
		_, text, err = book_service.GetBookTxt(bookName, author, chapterName, "")
		if err != nil {
			isLess = true
			log.Println("err:", err.Error())
			continue
		}
		if text != "" {
			textNum := len([]rune(strings.TrimSpace(text)))
			if textNum <= 1500 {
				//log.Println("尝试访问 text", bookId, bookName, author, chapterName, textNum, text)
				text, err = PaoshubaGetChapterText1(bookName, chapterName, chapterLink)
				if err != nil {
					global.Paoshu8log.Infof("采集章节内容失败 bookName=%v chapterTitle=%v chapterLink=%v ", bookName, chapterName, chapterLink)
					continue
				}
				textNum = len([]rune(strings.TrimSpace(text)))
				if textNum <= 1500 {
					domain := utils.GetUrlDomain(chapterLink)
					chapterLink = strings.ReplaceAll(chapterLink, domain, "http://www.biquge5200.net")
					text1, err := PaoshubaGetChapterText1(bookName, chapterName, chapterLink)
					log.Println("再次尝试访问 text", bookId, bookName, author, chapterName, textNum, text1)
					if text1 != "" && err == nil {
						text = text1
					}
					//if err != nil {
					//	global.Paoshu8log.Infof("采集章节内容失败 bookName=%v chapterTitle=%v chapterLink=%v ", bookName, chapterName, chapterLink)
					//	continue
					//}
					textNum = len([]rune(strings.TrimSpace(text)))
					if text != "" {
						_, text, err = book_service.GetBookTxt(bookName, author, chapterName, text)
						if err != nil {
							continue
						}
					}

				}
			}

			//var replaceText1 = "从远端拉取内容失败，有可能是对方服务器响应超时，后续待更新"
			//if !strings.Contains(text, replaceText1) {
			//	continue
			//}
			//isLess = true
			//global.Errlog.Errorf("%v %v %v %v %v", book.Id, bookName, author, chapterName, text)
			//if text == replaceText1 {
			//	var txtFile string
			//	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
			//	isFew = true
			//	//log.Println("txtFile", book.BookName, book.Author, chapterName, chapterNameMd5, txtFile)
			//	continue
			//}
			//log.Println("text", bookName, author, chapterName, chapterNameMd5, len(text))

			//if strings.Contains(text, "try{mad1('gad2');} catch(ex){}") {
			//	text = strings.ReplaceAll(text, "try{mad1('gad2');} catch(ex){}", "")
			//}
			//if strings.Contains(text, "try{") {
			//	text = strings.ReplaceAll(text, "try{", "")
			//}
			//if strings.Contains(text, "mad1('gad2');}") {
			//	text = strings.ReplaceAll(text, "mad1('gad2');}", "")
			//}
			//if strings.Contains(text, "catch(ex){}") {
			//	text = strings.ReplaceAll(text, "catch(ex){}", "")
			//}
			//if strings.Contains(text, "mad1(") {
			//	text = strings.ReplaceAll(text, "mad1(", "")
			//}
			//if strings.Contains(text, "gad2") {
			//	text = strings.ReplaceAll(text, "gad2", "")
			//}
			//if strings.Contains(text, "')") {
			//	text = strings.ReplaceAll(text, "')", "")
			//}
			//if strings.Contains(text, ";}") {
			//	text = strings.ReplaceAll(text, ";}", "")
			//}
			//if strings.Contains(text, "'')") {
			//	text = strings.ReplaceAll(text, "'')", "")
			//}
			//if strings.Contains(text, "content1()") {
			//	text = strings.ReplaceAll(text, "content1()", "")
			//}
			//if strings.Contains(text, "<br>") {
			//	text = strings.ReplaceAll(text, "<br>", "/r/n")
			//}

			//log.Println("len章节内容", bookName, author, chapterName, len(text))
		}
	}
	if isLess {
		global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Update("is_less", 1)
	}
	log.Println(book.Id, "结束")
}

func PaoshubaGetChapterText1(bookName, chapterTitle, chapterLink string) (text string, err error) {
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
	text = strings.ReplaceAll(text, "www.paoshu8.info", "")
	text = strings.ReplaceAll(text, "正在手打中，请稍等片刻，内容更新后，请重新刷新页面，即可获取最新更新！", "")
	text = strings.ReplaceAll(text, "支持（狂沙文学网）把本站分享那些需要的小伙伴！找不到书请留言！", "")
	text = strings.ReplaceAll(text, "剩余内容请前往纵横小说继续阅读.百度或各大应用市场搜索“纵横小说”，仙侠玄幻雪中,一剑奇幻武侠,青鸾,土豆脑洞剑道第一仙为生活添点料。或直接访问www.zongheng.com", "")
	text = strings.ReplaceAll(text, "☆★☆★☆", "")
	return
}
