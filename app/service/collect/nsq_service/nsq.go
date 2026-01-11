package nsq_service

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/collect_service"
	"go-novel/app/service/common/source_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"regexp"
)

// nsq订阅消息
type ConsumerT struct {
	Topic string
}

func (ct *ConsumerT) HandleMessage(msg *nsq.Message) error {
	//log.Println("receive", ct.Topic, fmt.Sprintf("%v", msg.Attempts), fmt.Sprintf("%x", msg.ID), "message:", string(msg.Body))
	_ = hanleTopic(ct.Topic, msg)
	return nil
}

func hanleTopic(topic string, msg *nsq.Message) (err error) {
	if topic == "utils.UpdateChapter" {
		bookSource := models.McBookSource{}
		err = json.Unmarshal(msg.Body, &bookSource)
		if err != nil {
			global.Updatelog.Errorf("%v", err.Error())
			return
		}
		bookSourceId := bookSource.Id
		bookId := bookSource.Bid
		sourceId := bookSource.Sid
		sourceUrl := bookSource.SourceUrl

		var book *models.McBook
		book, err = book_service.GetBookById(bookId)
		if book.Id <= 0 {
			err = delBookSourceById(bookSourceId)
			if err != nil {
				global.Updatelog.Errorf("%v", err.Error())
				return
			}
		}
		if err != nil {
			global.Updatelog.Errorf("%v", err.Error())
			return
		}
		log.Printf("bookSourceId=%v", bookSourceId)
		//msg.Finish()
		//return err
		//bookName := book.BookName
		var collect *models.McCollect
		collect, err = collect_service.GetCollectById(sourceId)
		if err != nil {
			global.Updatelog.Errorf("bookId=%v 获取源出错了", bookId)
			return
		}
		collectId := collect.Id
		chapterSection := collect.ChapterSectionReg
		chapterUrl := collect.ChapterUrlReg

		if chapterSection == "" {
			global.Updatelog.Errorf("bookId=%v collectId=%v 源获取章节区间正则不能为空", bookId, collectId)
			return
		}

		var html string
		sleepSecond := getSleepSecond()
		html, err = utils.GetHtml(sourceUrl, collect.Charset, collect.UrlComplete, sleepSecond)
		if html == "" {
			global.Updatelog.Errorf("获取小说详情页面失败 sourceUrl=%v", sourceUrl)
			return
		}

		var chapters []*models.CollectChapterInfo
		matchChapterSection := regexp.MustCompile(chapterSection).FindStringSubmatch(html)
		if len(matchChapterSection) > 0 {
			ulContent := matchChapterSection[0]
			liMatches := regexp.MustCompile(chapterUrl).FindAllStringSubmatch(ulContent, -1)
			for _, liMatch := range liMatches {
				if len(liMatch) > 2 {
					href := liMatch[1]
					chapterName := liMatch[2]
					chapterInfo := models.CollectChapterInfo{
						ChapterLink:  href,
						ChapterTitle: chapterName,
					}
					chapters = append(chapters, &chapterInfo)
				}
			}
		} else {
			global.Updatelog.Errorf("获取小说章节区间失败 %v sourceUrl=%v", sourceUrl)
			return
		}

		//章节表
		//var gq *gojsonq.JSONQ
		//gq, err = chapter_service.GetJsonqByBookId(bookId)
		//if err != nil {
		//	return
		//}
		var chapterFile string
		chapterFile, err = chapter_service.GetChapterFile(book.BookName, book.Author)
		if err != nil {
			return
		}
		var dbChapters []string
		_, dbChapters, err = chapter_service.GetChapterNamesByFile(chapterFile)
		if err != nil {
			return
		}
		lastSort := chapter_service.GetSortLast(book.BookName, book.Author)

		// 存储需要更新的章节列表
		var updatedChapters []*models.McBookChapter
		// 检查采集的章节是否需要更新
		for _, val := range chapters {
			chapterLink := val.ChapterLink
			chapterName := val.ChapterTitle
			lastSort += 1
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
				//chapter := &models.ChapterTextReq{
				//	Collect:      collect,
				//	BookName:     bookName,
				//	BookId:       bookId,
				//	ChapterTitle: chapterName,
				//	ChapterLink:  chapterLink,
				//}
				//if err = global.NsqPro.Publish(utils.UpdateChapterText, []byte(utils.JSONString(chapter))); err != nil { // 发布消息
				//	global.Nsqlog.Errorf("队列传输错误 %v", err.Error())
				//	continue
				//}
				updatedChapter := &models.McBookChapter{
					Sort:        lastSort,
					ChapterLink: chapterLink,
					ChapterName: chapterName,
					Vip:         0,
					Cion:        0,
					TextNum:     2000,
					Addtime:     utils.GetUnix(),
				}
				updatedChapters = append(updatedChapters, updatedChapter)
			}
		}

		if len(updatedChapters) > 0 {
			//batchSize := 1000 // 每批数据的大小
			//if err = global.DB.Table(chapterTable).Debug().CreateInBatches(updatedChapters, batchSize).Error; err != nil {
			//	global.Updatelog.Errorf("sql 书籍章节添加失败，稍后再试 err=%v ", err.Error())
			//	return
			//}
		}

		var updateTime, chapterTitle string
		_, updateTime, chapterTitle, err = source_service.GetSourceLastChapter(sourceUrl, collect)
		if err != nil {
			global.Updatelog.Errorf("%v", err.Error())
			return
		}
		//更新最新书信息和章节列表
		data := make(map[string]interface{})
		data["last_chapter_time"] = updateTime
		data["last_chapter_title"] = chapterTitle
		data["uptime"] = utils.GetUnix()
		err = updateBookSource(bookSourceId, data)
		global.Collectlog.Infof("%v", "更新章节成功")
		msg.Finish()
		return
	} else if topic == "utils.UpdateChapterText" {
		chapterText := models.ChapterTextReq{}
		err = json.Unmarshal(msg.Body, &chapterText)
		if err != nil {
			global.Nsqlog.Errorf("%v", err.Error())
			return
		}
		//log.Printf("chapterText %+v", chapterText)
		//msg.Finish()
		//return
		collect := chapterText.Collect
		charset := collect.Charset
		chapterPattern := collect.ChapterTextReg
		chapterTitle := chapterText.ChapterTitle
		chapterLink := chapterText.ChapterLink
		bookId := chapterText.BookId
		bookName := chapterText.BookName
		author := chapterText.Author

		var chapterHtml, text string
		chapterHtml, err = utils.GetHtml(chapterLink, charset, 0, 0)
		if err != nil {
			return
		}
		text, err = collect_service.GetChapterContent(chapterPattern, chapterHtml)
		if err != nil {
			global.Collectlog.Errorf("%v", err.Error())
			return
		}
		_, text, err = book_service.GetBookTxt(bookName, author, chapterTitle, text)
		if err != nil {
			return
		}
		log.Println(bookId, bookName, chapterTitle, chapterLink, "写入最新章节成功")
		msg.Finish()
	} else if topic == utils.ChapterText {
		chapterText := models.ChapterTextReq{}
		err = json.Unmarshal(msg.Body, &chapterText)
		if err != nil {
			global.Nsqlog.Errorf("%v", err.Error())
			return
		}
		//log.Printf("chapterText %+v", chapterText)
		//log.Printf("ChapterTextReg %+v", chapterText.Collect.ChapterTextReg)
		//msg.Finish()
		//return
		collect := chapterText.Collect
		charset := collect.Charset
		textReplace := collect.TextReplaceReg
		ChapterTextReg := collect.ChapterTextReg
		chapterTitle := chapterText.ChapterTitle
		chapterLink := chapterText.ChapterLink
		bookId := chapterText.BookId
		bookName := chapterText.BookName
		author := chapterText.Author

		var replaces []*models.TextReplace
		if textReplace != "" {
			err = json.Unmarshal([]byte(textReplace), &replaces)
			if err != nil {
				global.Collectlog.Errorf("解析collect小说替换出错 collectId=%v err=%v", collect.Id, err.Error())
				err = fmt.Errorf("解析collect小说内容替换出错 collectId=%v err=%v", collect.Id, err.Error())
				return
			}
		}

		var chapterHtml, text string
		chapterHtml, err = utils.GetHtml(chapterLink, charset, 0, 0)
		if err != nil {
			return
		}
		text, err = collect_service.GetChapterContent(ChapterTextReg, chapterHtml)
		if err != nil {
			global.Collectlog.Errorf("%v", err.Error())
			return
		}
		if len(replaces) > 0 {
			text = utils.ReplaceWords(text, replaces)
		}
		_, text, err = book_service.GetBookTxt(bookName, author, chapterTitle, text)
		if err != nil {
			return
		}
		log.Println(bookId, bookName, chapterTitle, chapterLink, "写入成功")
		msg.Finish()
	} else if topic == utils.SourceUpdateLastChapter {
		chapterMsg := models.NsqChapterInfoPush{}
		err = json.Unmarshal(msg.Body, &chapterMsg)
		if err != nil {
			global.Nsqlog.Errorf("%v", err.Error())
			return
		}
		bookName := chapterMsg.BookName
		author := chapterMsg.Author
		textNum := chapterMsg.TextNum
		chapterTitle := chapterMsg.ChapterTitle
		chapterLink := chapterMsg.ChapterLink
		chapterText := chapterMsg.ChapterText
		updatedChapter := &models.McBookChapter{
			ChapterName: chapterTitle,
			ChapterLink: chapterLink,
			TextNum:     textNum,
		}
		err = chapter_service.CreateChapter(bookName, author, updatedChapter)
		if err != nil {
			global.Jsonq.Errorf("%v", err.Error())
			return
		}
		_, _, err = book_service.GetBookTxt(bookName, author, chapterTitle, chapterText)
		if err != nil {
			return
		}
		msg.Finish()
	}
	msg.Finish()
	return
}

// 初始化消费者
func InitConsumer(topic string, channel string, address string) {
	cfg := nsq.NewConfig()
	//cfg.DefaultRequeueDelay = time.Second * 1
	cfg.MaxAttempts = 1000
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		global.Nsqlog.Errorf("InitConsumer %s", err.Error())
		return
	}
	c.SetLogger(nil, 3) //屏蔽系统日志

	ct := &ConsumerT{}
	ct.Topic = topic
	c.AddHandler(ct) // 添加消费者接口

	//建立NSQLookupd连接
	if err = c.ConnectToNSQLookupd(address); err != nil {
		global.Nsqlog.Errorf("ConnectToNSQLookupd %s", err.Error())
		return
	}
}

//func SaveData(collect *models.McCollect, info *models.CollecBookInfoRes) (err error) {
//	book, err := book_service.GetBookByBookName(info.BookName)
//	if err != nil {
//		global.Collectlog.Errorf("%v", err.Error())
//		return
//	}
//	bookName := info.BookName
//	var serialize int
//	if info.Serialize == "连载中" {
//		serialize = 1
//	} else if info.Serialize == "完本" {
//		serialize = 2
//	}
//	chapterNum := int64(len(info.Chapters))
//	textNum := chapterNum * 2000
//
//	var chapterTitle string
//	if len(info.Chapters) > 0 {
//		chapterTitle = info.Chapters[len(info.Chapters)-1].ChapterTitle
//	}
//
//	if book.Id <= 0 {
//		hits, hitsDay, hitsWeek, hitsMonth, shits, score, readCount := utils.GetRandNumBookHits()
//		book = &models.McBook{
//			BookName:     bookName,
//			Pic:          info.Pic,
//			Cid:          0,
//			IsRec:        0,
//			IsChoice:     0,
//			IsClassic:    0,
//			SearchNum:    0,
//			Serialize:    serialize,
//			Author:       info.Author,
//			Tags:         info.TagName,
//			Desc:         info.Desc,
//			TextNum:      textNum,
//			Hits:         hits,
//			HitsMonth:    hitsMonth,
//			HitsWeek:     hitsWeek,
//			HitsDay:      hitsDay,
//			Shits:        shits,
//			IsPay:        0,
//			ChapterNum:   chapterNum,
//			Score:        float64(score),
//			SourceId:     collect.Id,
//			SourceUrl:    info.BookUrl,
//			ChapterId:    0,
//			ChapterTitle: chapterTitle,
//			ReadCount:    readCount,
//			Addtime:      utils.GetUnix(),
//		}
//		if err = global.DB.Debug().Create(book).Error; err != nil {
//			global.Collectlog.Errorf("sql 书籍添加 记录失败，稍后再试 err=%v", err.Error())
//			return
//		}
//	}
//
//	//章节表
//	chapterTable, err := book_service.GetChapterTable(book.Id)
//	if err != nil {
//		global.Collectlog.Errorf("%v", "生成章节表失败")
//		return
//	}
//	if len(info.Chapters) <= 0 {
//		global.Collectlog.Errorf("采集章节为空 bookName=%v", bookName)
//		return
//	}
//
//	dbChapters := book_service.GetChapterNames(chapterTable)
//	lastSort := book_service.GetSortLast(book.Id)
//
//	// 存储需要更新的章节列表
//	var updatedChapters []*models.McBookChapter
//	var collectChapters []*models.CollectChapterInfo
//	// 检查采集的章节是否需要更新
//	for _, val := range info.Chapters {
//		chapterLink := val.ChapterLink
//		chapterName := val.ChapterTitle
//		lastSort += 1
//		// 查找章节是否存在于数据库中
//		found := false
//		for _, dbChapter := range dbChapters {
//			if chapterName == dbChapter {
//				found = true
//				break
//			}
//		}
//		// 如果章节不存在于数据库中，则需要更新
//		if !found {
//			chapter := &models.CollectChapterInfo{
//				ChapterTitle: chapterName,
//				ChapterLink:  chapterLink,
//			}
//			collectChapters = append(collectChapters, chapter)
//
//			updatedChapter := &models.McBookChapter{
//				Sort:        lastSort,
//				ChapterName: chapterName,
//				Vip:         0,
//				Cion:        0,
//				TextNum:     2000,
//				Addtime:     utils.GetUnix(),
//			}
//			updatedChapters = append(updatedChapters, updatedChapter)
//		}
//	}
//
//	if len(updatedChapters) > 0 {
//		batchSize := 1000 // 每批数据的大小
//		if err = global.DB.Table(chapterTable).CreateInBatches(updatedChapters, batchSize).Error; err != nil {
//			global.Collectlog.Errorf("sql 书籍章节添加失败，稍后再试 err=%v ", err.Error())
//			return
//		}
//	}
//
//	// 输出需要更新的章节列表
//	chapterLast := book_service.GetChapterIdLast(chapterTable)
//
//	//更新最新书信息和章节列表
//	data := make(map[string]interface{})
//	data["serialize"] = serialize
//	data["tags"] = info.TagName
//	data["desc"] = info.Desc
//	data["book_type"] = getClassType(info.CategoryId)
//	data["cid"] = info.CategoryId
//	data["chapter_num"] = chapterNum
//	data["chapter_id"] = chapterLast.Id
//	data["chapter_title"] = chapterLast.ChapterName
//	err = book_service.UpdateBookInfoByName(bookName, data)
//	global.Collectlog.Infof("%v", "更新章节成功")
//
//	chapterText := &models.ChapterTextReq{
//		Collect:  collect,
//		BookName: book.BookName,
//		BookId:   book.Id,
//		//Chapters: collectChapters,
//	}
//	if err = global.NsqPro.Publish(utils.ChapterText, []byte(utils.JSONString(chapterText))); err != nil { // 发布消息
//		global.Nsqlog.Errorf("队列传输错误 %v", err.Error())
//		return
//	}
//	return
//}

func getClassType(classId int64) (classType int) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("book_type").Where("id = ?", classId).Scan(&classType).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
