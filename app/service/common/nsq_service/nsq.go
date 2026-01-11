package nsq_service

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/collect/collect_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
	"time"
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
	if topic == utils.SourceUpdateLastChapter {
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
			global.Jsonq.Errorf("%v", err.Error())
			return
		}
		msg.Finish()
	} else if topic == utils.UpdateBook {
		bookMsg := models.NsqCollectBookPush{}
		err = json.Unmarshal(msg.Body, &bookMsg)
		if err != nil {
			global.Nsqlog.Errorf("%v", err.Error())
			msg.RequeueWithoutBackoff(time.Second * 3)
			return
		}
		bookName := strings.TrimSpace(bookMsg.BookName)
		author := strings.TrimSpace(bookMsg.Author)
		pic := strings.TrimSpace(bookMsg.Pic)
		classId := bookMsg.ClassId
		categoryName := strings.TrimSpace(bookMsg.CategoryName)
		serialize := bookMsg.Serialize
		tags := strings.TrimSpace(bookMsg.Tags)
		desc := strings.TrimSpace(bookMsg.Desc)
		textNum := bookMsg.TextNum
		chapterNum := bookMsg.ChapterNum
		sourceId := bookMsg.SourceId
		sourceUrl := strings.TrimSpace(bookMsg.SourceUrl)
		updateChapterId := bookMsg.UpdateChapterId
		updateChapterTitle := strings.TrimSpace(bookMsg.UpdateChapterTitle)
		updateChapterTime := bookMsg.UpdateChapterTime
		lastChapterTitle := strings.TrimSpace(bookMsg.LastChapterTitle)
		lastChapterTimeStr := strings.TrimSpace(bookMsg.LastChapterTime)
		bookType := bookMsg.BookType
		isClassic := bookMsg.IsClassic
		lastChapterTime := utils.DateToUnix(lastChapterTimeStr)

		if updateChapterTitle == "" {
			var gq *gojsonq.JSONQ
			gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
			if err != nil {
				global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
				msg.RequeueWithoutBackoff(time.Second * 3)
				return
			}

			var lastSortChapter *models.McBookChapter
			lastSortChapter, _ = chapter_service.GetLast(gq, "sort")
			if lastSortChapter != nil {
				updateChapterId = lastSortChapter.Id
				updateChapterTitle = lastSortChapter.ChapterName
			}
		}

		var book = new(models.McBook)
		book, err = book_service.GetBookByBookName(bookName, author)
		if err != nil {
			global.Collectlog.Errorf("%v", err.Error())
			return
		}
		var className string
		if classId > 0 {
			className = getClassName(classId)
		}

		if book.Id <= 0 {
			hits, hitsDay, hitsWeek, hitsMonth, shits, score, readCount, searchCount := utils.GetRandNumBookHits()
			if bookType <= 0 {
				bookType = getClassType(classId)
			}
			book = &models.McBook{
				BookType:           bookType,
				BookName:           bookName,
				Author:             author,
				Pic:                pic,
				Cid:                classId,
				CategoryName:       categoryName,
				ClassName:          className,
				Serialize:          serialize,
				Tags:               tags,
				Desc:               desc,
				TextNum:            textNum,
				ChapterNum:         chapterNum,
				SourceId:           sourceId,
				SourceUrl:          sourceUrl,
				UpdateChapterId:    updateChapterId,
				UpdateChapterTitle: updateChapterTitle,
				UpdateChapterTime:  updateChapterTime,
				LastChapterTitle:   lastChapterTitle,
				LastChapterTime:    lastChapterTime,
				IsClassic:          isClassic,
				ReadCount:          readCount,
				SearchCount:        searchCount,
				Status:             1,
				Score:              score,
				Hits:               hits,
				HitsDay:            hitsDay,
				HitsWeek:           hitsWeek,
				HitsMonth:          hitsMonth,
				Shits:              shits,
				Addtime:            utils.GetUnix(),
			}
			if err = global.DB.Create(book).Error; err != nil {
				global.Collectlog.Errorf("sql 书籍添加 记录失败，稍后再试 err=%v", err.Error())
				return
			}
		} else {
			//更新最新书信息和章节列表
			data := make(map[string]interface{})
			if serialize > 0 {
				data["serialize"] = bookMsg.Serialize
			}
			if tags != "" {
				data["tags"] = bookMsg.Tags
			}
			if desc != "" {
				data["desc"] = desc
			}
			if classId > 0 {
				if bookType <= 0 {
					bookType = getClassType(classId)
				}
				data["cid"] = classId
			}
			if bookType > 0 {
				data["book_type"] = bookType
			}
			if categoryName != "" {
				data["category_name"] = categoryName
			}
			if className != "" {
				data["class_name"] = className
			}
			if chapterNum > 0 {
				data["chapter_num"] = chapterNum
			}
			if sourceId > 0 {
				data["source_id"] = sourceId
			}
			if sourceUrl != "" {
				data["source_url"] = sourceUrl
			}
			if updateChapterId > 0 {
				data["update_chapter_id"] = updateChapterId
			}
			if updateChapterTitle != "" {
				data["update_chapter_title"] = updateChapterTitle
			}
			if updateChapterTime > 0 {
				data["update_chapter_time"] = updateChapterTime
			}
			if lastChapterTitle != "" {
				data["last_chapter_title"] = lastChapterTitle
			}
			if lastChapterTime > 0 {
				data["last_chapter_time"] = lastChapterTime
			}
			if serialize > 0 {
				data["serialize"] = serialize
			}
			data["uptime"] = utils.GetUnix()
			err = book_service.UpdateBookInfoByName(bookName, author, data)
			if err != nil {
				global.Sqllog.Errorf("%v", err.Error())
				return
			} else {
				global.Sqllog.Infof("bookName=%v 更新小说详情成功", bookName)
				return
			}
		}
		if book == nil {
			msg.RequeueWithoutBackoff(time.Second * 3)
			return
		}
		bookId := book.Id

		var count int64
		if sourceId > 0 {
			count = collect_service.GetSourceCountById(bookId, sourceId)
			if count <= 0 {
				bookSource := &models.McBookSource{
					Bid:              bookId,
					BookName:         bookName,
					Author:           author,
					Sid:              sourceId,
					SourceUrl:        sourceUrl,
					LastChapterTitle: lastChapterTitle,
					LastChapterTime:  lastChapterTimeStr,
					Addtime:          utils.GetUnix(),
				}
				if err = global.DB.Create(bookSource).Error; err != nil {
					global.Sqllog.Errorf("sql 书源添加 记录失败，稍后再试 err=%v", err.Error())
					return
				}
			}
		}

		msg.Finish()
		return
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
	if err = c.ConnectToNSQDs([]string{"127.0.0.1:4150", "localhost:4150"}); err != nil {
		global.Nsqlog.Errorf("ConnectToNSQDs %s", err.Error())
		return
	}
	//建立NSQLookupd连接
	if err = c.ConnectToNSQLookupd(address); err != nil {
		global.Nsqlog.Errorf("ConnectToNSQLookupd %s", err.Error())
		return
	}
}

func getClassType(classId int64) (classType int) {
	var err error
	err = global.DB.Model(models.McBookClass{}).Select("book_type").Where("id = ?", classId).Scan(&classType).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
