package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/redis_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/db"
	"go-novel/global"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"os"
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
	db.InitKeyLock()
	//updatedChapter := &models.McBookChapter{
	//	ChapterLink: "www.baidu.com",
	//	ChapterName: "百度",
	//	Vip:         0,
	//	Cion:        0,
	//	TextNum:     2000,
	//	Addtime:     utils.GetUnix(),
	//}
	//err := chapter_service.CreateChapter(1, updatedChapter)
	//log.Println(err)
	//page := getPage()

	//page := getPage()
	//page := 7
	//minId := page * 10000
	//maxId := minId + 20000
	//minId = 0
	//maxId = 10000
	//collect, _ := collect_service.GetCollectById(2)
	var books []*models.McBook
	global.DB.Model(models.McBook{}).Debug().Order("id asc").Find(&books)
	//global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("source_url like ?", "%"+"paoshu8"+"%").Where("id > ? and id <= ?", minId, maxId).Find(&books)
	//global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("is_down != 1").Where("source_url like ?", "%"+"paoshu8"+"%").Where("id > ? and id <= ?", minId, maxId).Find(&books)
	//global.DB.Model(models.McBook{}).Debug().Order("id desc").Where("is_few = 1").Where("source_url like ?", "%"+"biquge34"+"%").Find(&books)
	//global.DB.Model(models.McBook{}).Debug().Where("source_url like ?", "%"+"paoshu8"+"%").Where("id > ? and id <= ?", minId, maxId).Find(&books)
	log.Println("len", len(books))
	if len(books) <= 0 {
		return
	}
	//for _, book := range books {
	//	writeJson(book)
	//}
	//collect, err := collect_service.GetCollectById(1)
	//if err != nil {
	//	return
	//}
	for _, info := range books {
		//downChapter(info, collect)
		resetChapterMd5(info)
	}
}

// 2
func resetBookMd5(book *models.McBook) {
	var err error
	bookName := book.BookName
	author := book.Author
	bookMd5NameOld := utils.GetBookMd5(bookName, author)
	oldFile := fmt.Sprintf("/data/chapter/%v.json", bookMd5NameOld)
	bookMd5NameNew := utils.GetBookMd5(book.BookName, "")
	newFile := fmt.Sprintf("/data/chapter/%v.json", bookMd5NameNew)
	err = renameFile(oldFile, newFile)
	//log.Println(book.Id, book.BookName, "小说名称", oldFile, newFile)

	oldFile1 := fmt.Sprintf("/data/txt/%v", bookMd5NameOld)
	newFile1 := fmt.Sprintf("/data/txt/%v", bookMd5NameNew)
	err = renameFile(oldFile1, newFile1)
	//log.Println(book.Id, book.BookName, "章节目录", oldFile1, newFile1)
	//log.Println("结束", err)
}

func resetChapterMd5(book *models.McBook) {
	log.Println("开始")
	var err error
	bookName := book.BookName
	author := book.Author
	bookMd5Name := utils.GetBookMd5(bookName, author)
	var chapterFile string
	_, chapterFile, err = chapter_service.GetJsonqByBookName(book.BookName, book.Author)
	if err != nil {
		global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
		return
	}
	log.Println("bookid", book.Id)
	chapterAll, _, err := chapter_service.GetChapterNamesByFile(chapterFile)
	if len(chapterAll) <= 0 {
		da := make(map[string]interface{})
		da["is_less"] = 1
		global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Updates(da)
		return
	}
	var isLess bool
	for _, val := range chapterAll {
		if strings.Contains(val.ChapterName, "\n") {
			chapterNameMd5Old := utils.GetChapterMd5(val.ChapterName)
			val.ChapterName = strings.ReplaceAll(val.ChapterName, "\n", "")
			chapterNameMd5New := utils.GetChapterMd5(val.ChapterName)
			oldFile := fmt.Sprintf("/data/chapter/%v/%v.json", bookMd5Name, chapterNameMd5Old)
			newFile := fmt.Sprintf("/data/chapter/%v/%v.json", bookMd5Name, chapterNameMd5New)
			log.Println(book.Id, book.BookName, val.ChapterName, "章节目录", oldFile, newFile)
			err = renameFile(oldFile, newFile)
			isLess = true
		}
	}
	log.Println("结束", err)
	if isLess {
		global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Update("is_less", 3)
	}
}

func writeJson(book *models.McBook) {
	bookId := book.Id
	//章节表
	chapterTable, err := book_service.GetChapterTable(bookId)
	if err != nil {
		err = fmt.Errorf("%v", "获取章节失败")
		return
	}

	var uploadBookChapterPath string
	uploadBookChapterPath, err = setting_service.GetValueByName(utils.UploadBookChapterPath)
	if err != nil {
		err = fmt.Errorf("获取小说内容目录失败 uploadBookChapterPath=%v", uploadBookChapterPath)
		return
	}

	var chapterAll []*models.McBookChapter

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	filePath := uploadBookChapterPath

	err = utils.IsNotExistMkDir(filePath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}

	jsonFile := fmt.Sprintf("%v%v.json", filePath, bookId)
	if utils.CheckNotExist(jsonFile) {
		_, err = os.Create(jsonFile)
		if err != nil {
			return
		}
	}
	gq := gojsonq.New().File(jsonFile)

	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal("读取文件错误:", err)
	}

	if len(jsonData) > 0 {
		err = json.Unmarshal(jsonData, &chapterAll)
		if err != nil {
			log.Fatal("解析 JSON 数据错误:", err)
		}
	}

	var chapters []*models.McBookChapter
	global.DB.Table(chapterTable).Find(&chapters)
	if len(chapters) > 0 {
		for _, chapter := range chapters {
			gq.Reset()
			count := gq.Where("id", "=", chapter.Id).Count()
			if count <= 0 {
				chapterAll = append(chapterAll, chapter)
			}
		}
		//chapterAll = append(chapterAll, chapters...)
	}

	newJsonData, err := json.MarshalIndent(chapterAll, "", "  ")
	if err != nil {
		log.Fatal("转换为 JSON 数据错误:", err)
	}

	err = utils.WriteFile(jsonFile, string(newJsonData))
	log.Println(err)
}

func updatePic(book *models.McBook) (err error) {
	pic := book.Pic
	filePath := utils.GetFilePath(pic)
	ext := utils.GetExt(pic)
	filePath = "/data/pic"
	//fileBase := utils.GetFirstLetter(pic)
	//author := utils.GetFirstLetter(book.Author)
	//newPic := fmt.Sprintf("%v/%v.%v", filePath, fmt.Sprintf("%v-%v", fileBase, author), ext)
	newPic := fmt.Sprintf("%v/%v-%v.%v", filePath, utils.GetFirstLetter(book.BookName), utils.GetFirstLetter(book.Author), ext)
	//err = renameFile(pic, newPic)
	if err != nil {
		//log.Println(pic, newPic, err)
	} else {
		log.Println("处理图片成功", book.Id, pic, newPic)
		global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Update("pic", newPic)
	}
	//log.Fatalln(utils.GetFirstLetter(book.BookName), utils.GetFirstLetter(book.Author))

	//newPic := fmt.Sprintf("%v/%v.%v", filePath, strings.Join(pinyin.LazyPinyin(book.BookName, pinyin.NewArgs()), ""), ext)
	//log.Fatalln(pic, newPic)
	return
}

func downPic(book *models.McBook, collect *models.McCollect) (err error) {
	bookName := book.BookName
	author := book.Author
	bookUrl := book.SourceUrl
	picUrl := book.Pic
	bookName = utils.GetFirstLetter(bookName)
	author = utils.GetFirstLetter(author)
	pathDir := "/data/pic/"

	// 获取文件扩展名
	fileExt := strings.ToLower(utils.GetExt(picUrl))
	fileName := fmt.Sprintf("%v-%v.%v", bookName, author, fileExt)
	filePath1 := fmt.Sprintf("%s%s", pathDir, fileName)
	if utils.FileExist(filePath1) {
		global.Collectlog.Errorf("%v %v 文件已存在", bookUrl, filePath1)
		return
	}

	var html string
	html, err = utils.GetHtml(bookUrl, collect.Charset, collect.UrlComplete, 0)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var pic string
	if html != "" {
		matchPic := regexp.MustCompile(collect.PicReg).FindStringSubmatch(html)
		if len(matchPic) > 0 {
			pic = matchPic[1]
		}
	}
	var filePath string
	filePath, err = utils.DownImg(bookName, author, pic, pathDir)
	if err != nil {
		log.Println(pic, filePath, err)
	}
	log.Println(bookUrl, pic, filePath, err)
	return
}

func updateChapter(book *models.McBook) (err error) {
	newFile, err := GetChapterFile(book.BookName, book.Author)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var uploadBookChapterPath string = "/data/json/"
	//uploadBookChapterPath, err = setting_service.GetValueByName(utils.UploadBookChapterPath)
	//if err != nil {
	//	err = fmt.Errorf("获取小说内容目录失败 uploadBookChapterPath=%v", uploadBookChapterPath)
	//	return
	//}
	oldFile := fmt.Sprintf("%v%v.json", uploadBookChapterPath, book.Id)
	err = renameFile(newFile, oldFile)
	if err != nil {
		log.Println(oldFile, newFile, err)
	} else {
		log.Println(book.Id, "处理完", oldFile, newFile)
	}
	return
}

func updateChapterTxt(book *models.McBook) (err error) {
	//var uploadBookTextPath string
	//uploadBookTextPath, err = setting_service.GetValueByName(utils.UploadBookTextPath)
	//if err != nil {
	//	err = fmt.Errorf("获取小说内容目录失败 uploadBookTextPath=%v", uploadBookTextPath)
	//	return
	//}
	oldFile := fmt.Sprintf("%v%v", "/data/novel/", book.Id)
	bookName := utils.GetBookMd5(book.BookName, book.Author)
	newFile := fmt.Sprintf("%v%v", "/data/txt/", bookName)
	err = renameFile(oldFile, newFile)
	if err != nil {
		log.Println(oldFile, newFile, err)
	} else {
		log.Println(book.Id, "处理完", oldFile, newFile)
	}
	return
}

func updateRandom(book *models.McBook) (err error) {
	_, _, _, _, _, _, _, searchCount := utils.GetRandNumBookHits()
	data := make(map[string]interface{})
	//data["hits"] = hits
	//data["hits_day"] = hitsDay
	//data["hits_week"] = hitsWeek
	//data["hits_month"] = hitsMonth
	//data["shits"] = shits
	//data["score"] = score
	//data["read_count"] = readCount
	data["search_count"] = searchCount
	global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Updates(data)
	return
}

func GetChapterFile(bookName, author string) (chapterFile string, err error) {
	var uploadBookChapterPath string
	uploadBookChapterPath, err = setting_service.GetValueByName(utils.UploadBookChapterPath)
	if err != nil {
		err = fmt.Errorf("获取小说内容目录失败 uploadBookChapterPath=%v", uploadBookChapterPath)
		return
	}
	err = utils.IsNotExistMkDir(uploadBookChapterPath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}
	bookName = utils.GetBookMd5(bookName, author)
	chapterFile = fmt.Sprintf("%v%v.json", uploadBookChapterPath, bookName)
	return
}

func renameFile(oldPath, newPath string) error {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return fmt.Errorf("failed to rename file: %s", err)
	}
	return nil
}

func downChapter(book *models.McBook, collect *models.McCollect) (err error) {
	bookName := book.BookName
	author := book.Author
	collectId := collect.Id
	chaptersData, err := getChapters(book.SourceUrl)
	if err != nil {
		time.Sleep(time.Second * 3)
		utils.GetS5()
		global.Collectlog.Errorf("采集章节报错 len=%v err=%v", len(chaptersData), err.Error())
		//downChapter(book, collect)
		return
	}

	// 存储需要更新的章节列表
	var chapters []*models.CollectChapterInfo
	// 检查采集的章节是否需要更新
	for _, val := range chaptersData {
		chapterName := val.ChapterTitle
		chapterLink := val.ChapterLink

		// 如果章节不存在于数据库中，则需要更新
		chapter := &models.CollectChapterInfo{
			ChapterTitle: strings.TrimSpace(chapterName),
			ChapterLink:  strings.TrimSpace(chapterLink),
		}
		chapters = append(chapters, chapter)
	}

	collectChapter := models.CollectPageBookChapter{
		Collect:  collect,
		BookName: bookName,
		Author:   author,
		Chapters: chapters,
	}

	err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectChapterUrlTemp, collectId), collectChapter, 0)
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
		log.Println(threadMaxCount)
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

	//for {
	//	isChapterEnd := getNextChapterLink(collectId)
	//	if isChapterEnd {
	//		break
	//	}
	//}
	global.DB.Model(models.McBook{}).Where("id = ?", book.Id).Update("is_less", 99)
	return
}

func replaceTitle(title string) (newTitle string) {
	index := strings.Index(title, "第")
	if index != -1 {
		newTitle = title[index:]
	} else {
		newTitle = title
	}
	return
}
func getChapters(bookUrl string) (chapters []*models.CollectChapterInfo, err error) {
	var html string
	html, err = utils.GetHtml(bookUrl, "utf-8", 1, 0)
	if html == "" {
		err = fmt.Errorf("采集内容为空%v", err.Error())
		return
	}
	if html == "" {
		time.Sleep(time.Second)
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
		title := selection.Find("a").Text()
		chapter := &models.CollectChapterInfo{
			ChapterTitle: title,
			ChapterLink:  link,
		}
		chapters = append(chapters, chapter)
	})
	return
}

func GetChapterContent(bookName, chapterTitle, chapterLink string) (text string, err error) {
	var chapterHtml string
	chapterHtml, err = utils.GetHtml(chapterLink, "utf-8", 0, 0)
	if err != nil {
		return
	}
	log.Println(chapterHtml)
	if strings.Contains(chapterHtml, "503 Service Temporarily Unavailable") {
		global.Collectlog.Errorf("%v 503 Service Temporarily Unavailable", chapterLink)
		time.Sleep(time.Second)
		GetChapterContent(bookName, chapterTitle, chapterLink)
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
	text = strings.ReplaceAll(text, "＠樂＠文＠小＠说|", "")
	return
}

func processChapter(collect *models.McCollect, bookName, author string, chapters []*models.CollectChapterInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(chapters) <= 0 {
		return
	}
	// 生成100毫秒到2000毫秒之间的随机时间间隔
	//minDuration := 100 * time.Millisecond
	//maxDuration := 2000 * time.Millisecond
	//randomDuration := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
	var err error
	for _, chapter := range chapters {
		_, err = CollectChapterText(bookName, author, chapter.ChapterTitle, chapter.ChapterLink)
		if err != nil {
			global.Collectlog.Errorf("解析collect小说替换出错 ChapterTitle=%v ChapterLink=%v err=%v", chapter.ChapterTitle, chapter.ChapterLink, err.Error())
			time.Sleep(time.Second * 3)
			utils.GetS5()
			CollectChapterText(bookName, author, chapter.ChapterTitle, chapter.ChapterLink)
			continue
		}
	}
	return
}

func CollectChapterText(bookName, author, chapterTitle, chapterLink string) (textNum int, err error) {
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterTitle)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		global.Collectlog.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterTitle, txtFile, textNum)
		//err = fmt.Errorf("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterTitle, txtFile, textNum)
		return
	}
	var text string
	text, err = GetChapterContent(bookName, chapterTitle, chapterLink)
	if err != nil {
		time.Sleep(time.Second * 3)
		utils.GetS5()
		global.Collectlog.Errorf("bookName=%v chapterTitle=%v chapterLink = %v err=%v", bookName, chapterTitle, chapterLink, err.Error())
		return
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

func getNextChapterLink(collectId int64) (isEnd bool) {
	var err error
	var collectPageChapter = new(models.CollectPageBookChapter)
	chapterVal := redis_service.Get(fmt.Sprintf("%v_%v", utils.CollectChapterUrlTemp, collectId))
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
		err = redis_service.Del(fmt.Sprintf("%v_%v", utils.CollectChapterUrlTemp, collectId))
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
		textNum, err = CollectChapterText(bookName, author, chapterTitle, chapterLink)
		if err != nil && textNum > 0 {
			global.Collectlog.Errorf("解析collect小说替换出错 collectId=%v err=%v", collect.Id, err.Error())
			return
		}
		chapterLinks := removeChapterUrl(chapters, chapterLink)
		collectPageChapter.Chapters = chapterLinks
		err = redis_service.Set(fmt.Sprintf("%v_%v", utils.CollectChapterUrlTemp, collectId), collectPageChapter, 0)
		if err != nil {
			global.Collectlog.Errorf("缓存采集当前章节链接失败 %v", err.Error())
			return
		}
		return
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

func getPage() (page int) {
	var pageStr, ip, port, username, passwd string
	flag.StringVar(&pageStr, "page", "1", "default :page")
	flag.StringVar(&ip, "ip", "", "default :ip")
	flag.StringVar(&port, "port", "", "default :port")
	flag.StringVar(&username, "username", "", "default :username")
	flag.StringVar(&passwd, "passwd", "", "default :passwd")
	flag.Parse()
	utils.S5Domain = ip
	utils.S5Port = port
	utils.S5Username = username
	utils.S5Passwd = passwd
	page, _ = strconv.Atoi(pageStr)
	return
}
