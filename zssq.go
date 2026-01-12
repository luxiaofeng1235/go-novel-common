package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-novel/app/lib/aJson"
	"go-novel/app/lib/zssq"
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
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/ants/v2"
	"github.com/thedevsaddam/gojsonq/v2"
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
	var err error
	task := zssq.NewTask()
	//登录 生成一个新的游客用户
	_, err = task.Login()
	if err != nil {
		global.Zssqlog.Errorf("登录失败 %v", err.Error())
		return
	}
	_, err = task.GetLoginByYW()
	if err != nil {
		global.Zssqlog.Errorf("起点登录失败 %v", err.Error())
		return
	}
	for {
		ZssqCollect(task)
	}
	//var books []*models.McBook
	//global.DB.Model(models.McBook{}).Debug().Order("id asc").Where("is_less= 1 and source_url not like ?", "%"+"http"+"%").Find(&books)
	//for _, book := range books {
	//	ZssqRepairChapter(book, task)
	//}
}

func ZssqRepairChapter(book *models.McBook, task *zssq.Task) {
	//章节表
	//var chapterFile string
	var err error
	//_, chapterFile, err = chapter_service.GetJsonqByBookName(book.BookName, book.Author)
	//if err != nil {
	//	global.Collectlog.Errorf("获取JSONQ对象失败 %v", err.Error())
	//	return
	//}
	//log.Println("bookid", book.Id)
	//err = utils.RemoveFile(chapterFile)
	//return
	bookId := book.SourceUrl
	bookName := book.BookName
	author := book.Author
	gender := "male"
	categoryName := "都市"
	var bookInfo = new(models.ZssqBookDesc)
	bookInfo, err = ZssqGetBookDesc(task, gender, categoryName, bookId)
	if err != nil {
		global.Zssqlog.Errorf("获取小说详情失败 %v", err.Error())
		return
	}
	bookInfo.LastChapterTime, err = ZssqGetBookChaptersById(task, bookId, bookName, author)
	return
}

func ZssqGetBookChaptersById(task *zssq.Task, bookId, bookName, author string) (lastChapterTime string, err error) {
	var chapters []*models.ZssqChapter
	chapters, lastChapterTime, err = ZssqGetBookChapters(task, bookId, bookName, author)
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		return
	}
	var chapterAll []*models.McBookChapter
	var id int64
	var sort int
	for _, val := range chapters {
		id++
		sort++
		var textNum, isLess int
		textNum, err = ZssqCollectChapterText(bookId, bookName, author, task, val)
		if err != nil && textNum > 0 {
			global.Zssqlog.Errorf("获取章节内容失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
			continue
		}
		chapterName := val.ChapterName
		chapterLink := val.ChapterLink
		if textNum <= 1000 {
			isLess = 1
		}
		if textNum <= 100 {
			continue
		}
		updatedChapter := &models.McBookChapter{
			Id:          id,
			Sort:        sort,
			ChapterLink: chapterLink,
			ChapterName: chapterName,
			IsLess:      isLess,
			Vip:         0,
			Cion:        0,
			TextNum:     textNum,
			Addtime:     utils.GetUnix(),
		}
		chapterAll = append(chapterAll, updatedChapter)
	}
	return
}
func ZssqCollect(task *zssq.Task) {
	categorys, err := ZssqGetCacheCategory(task)
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		return
	}
	if len(categorys) <= 0 {
		return
	}
	var categoryName, categoryAlias, gender string
	var bookCount int
	for _, category := range categorys {
		if category.Use == 0 {
			categoryName = category.CategoryName
			categoryAlias = category.CategoryAlias
			gender = category.Gender
			bookCount = category.BookCount
			break
		}
	}
	if categoryName == "" {
		err = redis_service.Del(utils.ZssqCategory)
		if err != nil {
			global.Zssqlog.Errorf("%v", err.Error())
			return
		}
	}

	pageBookKey := fmt.Sprintf("%v_%v_%v", utils.ZssqBooks, gender, categoryAlias)
	vals := redis_service.Get(pageBookKey)
	var bookDescs []*models.ZssqBookDesc
	if vals != "" {
		err = json.Unmarshal([]byte(vals), &bookDescs)
		if err != nil {
			global.Zssqlog.Errorf("%v", err.Error())
			return
		}
	}

	if len(vals) <= 0 {
		//通过分类模糊搜索小说
		var begin = 0
		var limit = 50
		var books []*aJson.Json
		totalpage := int(bookCount / limit)
		for begin*limit <= bookCount {
			global.Zssqlog.Infoln(gender, categoryName, begin, totalpage, limit, bookCount)
			//没有超过总数
			//通过多次模糊搜索遍历分类中的所有书本
			var fuzzResp *aJson.Json
			fuzzResp, err = task.FuzzSearch(gender, categoryAlias, categoryName, begin, limit)
			if err != nil {
				global.Zssqlog.Errorf("读取小说列表失败 %v", err.Error())
				return
			}
			books, err = fuzzResp.Get("books").TryJsonArray()
			if err != nil {
				global.Zssqlog.Errorf("", err.Error())
				return
			}
			for _, book := range books {
				bookDesc := ZssqGetBookDescByData(book, gender, categoryName)
				bookDescs = append(bookDescs, bookDesc)
			}
			begin += 1
		}
		err = redis_service.Set(pageBookKey, bookDescs, 0)
		if err != nil {
			global.Zssqlog.Errorf("缓存小说列表出错 %v", err.Error())
			return
		}
	}

	var bookDesc = new(models.ZssqBookDesc)
	var bookId, bookName, author string
	var isFind bool
	for _, val := range bookDescs {
		if val.Use == 0 {
			bookId = val.BookKey
			bookName = val.BookName
			author = val.Author
			bookDesc = val
			isFind = true
			break
		}
	}
	if !isFind {
		err = ZssqRemoveCategory(pageBookKey, categorys, categoryAlias)
		if err != nil {
			return
		}
		return
	}

	//通过小说id查询小说详情
	//var bookInfo = new(models.ZssqBookDesc)
	//bookInfo, err = ZssqGetBookDesc(task, gender, categoryName, bookId)
	//if err != nil {
	//	global.Zssqlog.Errorf("获取小说详情失败 %v", err.Error())
	//	return
	//}
	var chapters []*models.ZssqChapter
	chapters, bookDesc.LastChapterTime, err = ZssqGetBookChapters(task, bookId, bookName, author)
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		err = ZssqRemoveBook(bookId, bookName, author, pageBookKey, bookDescs)
		if err != nil {
			return
		}
		return
	}

	ZssqManyThreadChapter(bookId, bookName, author, task, chapters)
	var chapterAll []*models.McBookChapter
	var id int64
	var sort int
	for _, val := range chapters {
		id++
		sort++
		chapterName := val.ChapterName
		chapterLink := val.ChapterLink
		var textNum, isLess int
		textNum, err = ZssqCollectChapterText(bookId, bookName, author, task, val)
		is204 := ZssqRemoveBook204(bookDescs, pageBookKey, bookId, bookName, err)
		if is204 {
			return
		}
		if err != nil && textNum > 0 {
			global.Zssqlog.Errorf("获取章节内容失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
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
			Sort:        sort,
			ChapterLink: chapterLink,
			ChapterName: chapterName,
			IsLess:      isLess,
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
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		return
	}
	err = ZssqRemoveBook(bookId, bookName, author, pageBookKey, bookDescs)
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		return
	}

	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Zssqlog.Errorf("获取JSONQ对象失败 %v", err.Error())
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
		Pic:                bookDesc.Pic,
		ClassId:            bookDesc.ClassId,
		CategoryName:       bookDesc.CategoryName,
		Desc:               bookDesc.Desc,
		Tags:               bookDesc.Tags,
		ChapterNum:         bookDesc.ChapterNum,
		SourceId:           0,
		SourceUrl:          bookId,
		LastChapterTitle:   bookDesc.LastChapterTitle,
		LastChapterTime:    bookDesc.LastChapterTime,
		UpdateChapterId:    updateChapterId,
		UpdateChapterTitle: updateChapterTitle,
		UpdateChapterTime:  utils.GetUnix(),
		TextNum:            bookDesc.TextNum,
		Serialize:          bookDesc.Serialize,
		BookType:           bookDesc.BookType,
		IsClassic:          bookDesc.IsClassic,
	}

	err = collect_service.NsqCollectBookPush(msg)
	if err != nil {
		global.Zssqlog.Errorf("%v", err.Error())
		return
	}
	return
}

func ZssqManyThreadChapter(bookId, bookName, author string, task *zssq.Task, chapters []*models.ZssqChapter) {
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
			_, err = ZssqCollectChapterText(bookId, bookName, author, task, val)
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

func ZssqRemoveBookUrl(bookDescs []*models.ZssqBookDesc, bookId string) (newList []*models.ZssqBookDesc) {
	for _, info := range bookDescs {
		if info.BookKey != bookId {
			newList = append(newList, info)
		}
	}
	return
}

func ZssqRemoveCategory(pageBookKey string, categorys []*models.ZssqCategory, categoryAlias string) (err error) {
	for _, val := range categorys {
		if val.CategoryAlias == categoryAlias {
			val.Use = 1
			//break 里面有重复值不能用break
		}
	}
	err = redis_service.Set(utils.ZssqCategory, categorys, 0)
	if err != nil {
		err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
		return
	}
	err = redis_service.Del(pageBookKey)
	if err != nil {
		global.Zssqlog.Errorf("删除key=%v出错 %v", pageBookKey, err.Error())
		return
	}
	return
}

func ZssqRemoveBook(bookId, bookName, author, pageBookKey string, bookDescs []*models.ZssqBookDesc) (err error) {
	for _, val := range bookDescs {
		if val.BookKey == bookId {
			val.Use = 1
		}
	}
	err = redis_service.Set(pageBookKey, bookDescs, 0)
	if err != nil {
		global.Zssqlog.Errorf("缓存采集状态失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}
	return
}

func ZssqRemoveBook204(bookDescs []*models.ZssqBookDesc, pageBookKey, bookId, bookName string, err error) (is204 bool) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "get response code of 204") {
		for _, val := range bookDescs {
			if val.BookKey == bookId {
				val.Use = 2
			}
		}
		err = redis_service.Set(pageBookKey, bookDescs, 0)
		if err != nil {
			global.Zssqlog.Errorf("缓存采集状态失败 bookName=%v err=%v", bookName, err.Error())
			return
		}
		is204 = true
	}
	return
}

func ZssqGetBookDescByData(bookInfoResp *aJson.Json, gender, categoryName string) (bookInfo *models.ZssqBookDesc) {
	categorys := []*models.CategoryReg{
		{"都市", 2},
		{"玄幻", 12},
		{"历史", 23},
		{"军事", 15},
		{"悬疑灵异", 24},
		{"科幻", 63},
		{"游戏体育", 81},
		{"同人", 12},
		{"玄幻脑洞", 12},
		{"都市脑洞", 2},
		{"历史脑洞", 23},
		{"东方玄幻", 12},
		{"乡村暧昧", 82},
		{"战神兵王", 83},
		{"无限穿越", 21},
		{"生活爽文", 1},
		{"娱乐明星", 84},
		{"青春校园", 17},
		{"都市修仙", 31},
		{"虚拟网游", 8},
		{"轻小说", 36},
		{"抗战烽火", 85},
		{"都市异能", 86},
		{"洪荒神话", 13},
		{"末世进化", 14},
		{"同人衍生", 21},
		{"秦汉三国", 22},
		{"异界争霸", 21},
		{"两晋隋唐", 22},
		{"探险盗墓", 24},
		{"恐怖惊悚", 24},
		{"游戏异界", 81},
		{"灵异神怪", 81},
		{"谍战特工", 85},
		{"古武传承", 12},
		{"剑与魔法", 34},
		{"经典武侠", 87},
		{"异类兽族", 12},
		{"现代军旅", 15},
		{"幻想修仙", 31},
		{"未来机甲", 88},
		{"西方神话", 12},
		{"职场生活", 89},
		{"宋元明清", 23},
		{"篮球风云", 2},
		{"足球天下", 2},
		{"短篇", 2},
		{"重生", 5},
		{"权谋", 16},
		{"后宫", 12},
		{"无敌", 35},
		{"西游", 12},
		{"鉴宝", 12},
		{"种田", 27},
		{"官场", 16},
		{"炼药炼器", 12},
		{"装逼打脸", 12},
		{"剑道", 33},
		{"逆袭", 21},
		{"玄学", 24},
		{"灵气复苏", 12},
		{"穿越", 21},
		{"废材", 21},
		{"强者回归", 21},
		{"召唤", 29},
		{"生存", 2},
		{"开局", 91},
		{"兄弟情", 21},
		{"年代", 5},
		{"战争", 93},
		{"经营", 2},
		{"创业", 5},
		{"聊天群", 94},
		{"系统", 9},
		{"封神", 4},
		{"争霸", 7},
		{"科举", 76},
		{"迪化", 65},
		{"抽奖", 21},
		{"直播", 95},
		{"火影", 96},
		{"英雄联盟", 6},
		{"王者荣耀", 6},
		{"鉴宝", 6},
		{"全民", 6},
		{"女尊", 45},
		{"护花", 45},
		{"美女", 45},
		{"神医", 79},
		{"赘婿", 90},
		{"特种兵", 15},
		{"神豪", 97},
		{"农民", 9},
		{"学生", 9},
		{"无节操", 9},
		{"杀伐果断", 9},
		{"成熟", 9},
		{"智商在线", 9},

		{"皇帝", 76},
		{"道士", 4},
		{"至尊", 58},
		{"天才", 58},
		{"仙帝", 4},
		{"奶爸", 12},
		{"草根", 2},
		{"领主", 12},
		{"学霸", 2},
		{"高手", 2},
		{"反派", 3},
		{"杀手", 21},
		{"思路清奇", 2},
		{"凡人", 4},
		{"霸气", 4},
		{"职业选手", 42},
		{"配角", 2},
		{"宅男", 2},
		{"主播", 2},
		{"现代言情", 45},
		{"古代言情", 45},
		{"豪门总裁", 44},
		{"唯美纯爱", 45},
		{"玄幻仙侠", 12},
		{"青春校园", 17},
		{"女强女尊", 47},
		{"宫斗宅斗", 75},
		{"幻想言情", 45},
		{"现代纯爱", 45},
		{"烧脑悬疑", 24},
		{"现代异能", 86},
		{"灵异怪谈", 24},
		{"科幻空间", 63},
		{"末世危机", 14},
		{"星际科幻", 63},
		{"年代文", 5},
		{"都市情缘", 2},
		{"婚恋情感", 2},
		{"娱乐圈", 2},
		{"古代情缘", 21},
		{"悬疑推理", 24},
		{"游戏情缘", 81},
		{"业界精英", 2},
		{"民国情仇", 2},
		{"古代纯爱", 21},
		{"短篇", 2},
		{"重生", 5},
		{"穿越", 21},
		{"女强", 47},
		{"情有独钟", 2},
		{"打脸", 2},
		{"架空", 10},
		{"穿书", 69},
		{"逆袭", 11},
		{"复仇", 7},
		{"空间", 12},
		{"家长里短", 2},
		{"欢喜冤家", 2},
		{"养成", 2},
		{"虐恋", 68},
		{"系统", 9},
		{"治愈", 2},
		{"扮猪吃虎", 21},
		{"团宠", 45},
		{"双洁", 45},
		{"契约", 45},
		{"权谋", 45},
		{"日久生情", 45},
		{"强强", 45},
		{"虐渣", 45},
		{"马甲", 21},
		{"修仙", 31},
		{"美食", 99},
		{"青梅竹马", 2},
		{"金手指", 98},
		{"异能", 86},
		{"先婚后爱", 45},
		{"发家致富", 5},
		{"双男主", 45},
		{"破镜重圆", 45},
		{"一见钟情", 45},
		{"暗恋", 45},
		{"清穿", 69},
		{"女扮男装", 45},
		{"闪婚", 2},
		{"前世今生", 45},
		{"相爱相杀", 45},
		{"隐婚", 45},
		{"兽世", 21},
		{"追妻火葬场", 2},
		{"双向奔赴", 45},
		{"异世", 21},
		{"姐弟恋", 45},
		{"灵魂互换", 80},
		{"别后重逢", 45},
		{"替嫁", 45},
		{"黑化", 45},
		{"竞技", 81},
		{"恐怖", 24},
		{"替身", 21},
		{"双重生", 21},
		{"师徒", 21},
		{"基建", 9},
		{"联姻", 45},
		{"电竞", 81},
		{"生存", 2},
		{"带球跑", 2},
		{"群穿", 69},
		{"初恋", 45},
		{"倒追", 45},
		{"逃婚", 45},
		{"失忆", 45},
		{"离婚", 45},
		{"读心术", 45},
		{"错嫁", 45},
		{"黑道", 45},
		{"古穿今", 69},
		{"智斗", 45},
		{"驭兽", 21},
		{"无限流", 21},
		{"腹黑", 45},
		{"霸道", 45},
		{"萌宝", 45},
		{"王妃", 45},
		{"神医", 79},
		{"嫡女", 45},
		{"女配", 45},
		{"王爷", 45},
		{"大佬", 45},
		{"病娇", 45},
		{"傲娇", 45},
		{"医生", 79},
		{"明星", 84},
		{"帝王", 76},
		{"农女", 2},
		{"千金", 45},
		{"真假千金", 45},
		{"杀伐果断", 45},
		{"皇后", 45},
		{"大叔", 45},
		{"特工", 93},
		{"萌宠", 45},
		{"偏执", 45},
		{"反派", 3},
		{"将军", 76},
		{"公主", 45},
		{"炮灰", 45},
		{"锦鲤", 45},
		{"庶女", 45},
		{"忠犬", 2},
		{"丑女", 45},
		{"全能", 45},
		{"魔君", 21},
		{"可盐可甜", 45},
		{"冷酷", 45},
		{"校草", 45},
		{"超A", 45},
		{"吃货", 45},
		{"杀手", 45},
		{"影后", 45},
		{"极品", 45},
		{"师尊", 45},
		{"暴君", 25},
		{"首席", 21},
		{"丧尸", 6},
		{"反差萌", 45},
		{"万人迷", 45},
		{"网红", 45},
		{"萝莉", 45},
		{"弃妇", 45},
		{"白月光", 45},
		{"职业选手", 42},
		{"悍妻", 45},
		{"上神", 21},
		{"小狼狗", 45},
		{"御姐", 45},
		{"作精", 45},
		{"女帝", 73},
		{"影帝", 45},
		{"老师", 45},
		{"法医", 45},
		{"财迷", 45},
		{"厨娘", 45},
		{"小妾", 56},
		{"纨绔", 56},
		{"甜妻", 56},
		{"侦探", 45},
		{"丫鬟", 56},
		{"男配", 56},
		{"宫女", 56},
		{"小奶狗", 45},
		{"警察", 45},
		{"精分", 45},
		{"白莲花", 45},
		{"天师", 56},
		{"律师", 45},
		{"公版书", 100},
		{"人文社科", 100},
		{"出版小说", 100},
		{"文学艺术", 100},
	}

	classId := utils.CategoryEquiv(categorys, categoryName)
	bookKey, _ := bookInfoResp.Get("_id").TryString()
	bookName, _ := bookInfoResp.Get("title").TryString()
	author, _ := bookInfoResp.Get("author").TryString()
	pic, _ := bookInfoResp.Get("cover").TryString()
	pic, _ = url.QueryUnescape(pic)
	pic = strings.TrimRight(strings.TrimLeft(pic, "/agent/"), "/")
	desc, _ := bookInfoResp.Get("longIntro").TryString()
	shortIntro, _ := bookInfoResp.Get("shortIntro").TryString()
	textNum, _ := bookInfoResp.Get("wordCount").TryInt()
	//categoryName, _ := bookInfoResp.Get("cat").TryString()
	majorCate, _ := bookInfoResp.Get("majorCate").TryString()
	chaptersCount, _ := bookInfoResp.Get("chaptersCount").TryInt()
	updated, _ := bookInfoResp.Get("updated").TryString()
	if categoryName == "" {
		categoryName = majorCate
	}
	if desc == "" {
		desc = shortIntro
	}
	var filePath string
	mondayDate := utils.GetThisWeekFirstDate()
	savePicPath := "/data/pic/" + mondayDate + "/" //拼装存储的图片路径
	filePath, _ = utils.DownImg(bookName, author, pic, savePicPath)
	pic = strings.TrimLeft(filePath, ".")
	bookTags, _ := bookInfoResp.Get("tags").TryJsonArray()
	tags := ""
	for _, tag := range bookTags {
		tags = tags + tag.Data + ","
	}
	tags = strings.TrimRight(tags, ",")
	var serialize int
	lastChapterName, _ := bookInfoResp.Get("lastChapter").TryString()
	isSerial, _ := bookInfoResp.Get("isSerial").TryBool()
	if isSerial {
		serialize = 1
	} else {
		serialize = 2
	}

	var bookType, isClassic int
	if gender == "male" {
		bookType = 1
	} else if gender == "female" {
		bookType = 2
	} else {
		isClassic = 1
	}
	bookInfo = &models.ZssqBookDesc{
		BookKey:          bookKey,
		BookName:         bookName,
		Author:           author,
		Pic:              pic,
		Desc:             desc,
		Serialize:        serialize,
		TextNum:          int(textNum),
		CategoryName:     categoryName,
		ClassId:          classId,
		Tags:             tags,
		LastChapterTitle: lastChapterName,
		LastChapterTime:  utils.ParseDateTzTime(updated),
		ChapterNum:       int(chaptersCount),
		BookType:         bookType,
		IsClassic:        isClassic,
	}
	return
}

func ZssqGetCategory(task *zssq.Task) (categorys []*models.ZssqCategory) {
	//展示所有分类
	tagResp, err := task.SearchTag()
	if err != nil {
		return
	}
	firstTags, err := tagResp.Get("firstTags").TryJsonArray()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var gender = ""
	var alias = ""
	var name = ""
	var bookCount = int32(0)
	for _, firstTag := range firstTags {
		secondTags, _ := firstTag.Get("secondTags").TryJsonArray()
		for _, secondTag := range secondTags {
			thirdTags, _ := secondTag.Get("thirdTags").TryJsonArray()
			for _, thirdTag := range thirdTags {
				alias, _ = thirdTag.Get("alias").TryString()
				name, _ = thirdTag.Get("name").TryString()
				gender, _ = thirdTag.Get("gender").TryString()
				// 检查 originalString 是否包含 subStringToRemove
				subStringToRemove := "_v2"
				if strings.Contains(gender, subStringToRemove) {
					// 移除 subStringToRemove，替换为 ""
					gender = strings.ReplaceAll(gender, subStringToRemove, "")
					// fmt.Println(updatedString) // 输出: "female"
				}
				bookCount, _ = thirdTag.Get("bookCount").TryInt()
				category := &models.ZssqCategory{
					CategoryName:  name,
					CategoryAlias: alias,
					BookCount:     int(bookCount),
					Gender:        gender,
				}
				categorys = append(categorys, category)
			}
		}
	}
	return
}

func ZssqGetCacheCategory(task *zssq.Task) (categorys []*models.ZssqCategory, err error) {
	categoryVal := redis_service.Get(utils.ZssqCategory)
	if categoryVal == "" || categoryVal == "null" {
		categorys = ZssqGetCategory(task)
		err = redis_service.Set(utils.ZssqCategory, categorys, 0)
		if err != nil {
			err = fmt.Errorf("缓存分类失败 err=%v", err.Error())
			return
		}
		return
	} else {
		err = json.Unmarshal([]byte(categoryVal), &categorys)
		if err != nil {
			global.Zssqlog.Errorf("获取分类缓存失败 err=%v", err.Error())
			return
		}
	}
	return
}

func ZssqGetBookDesc(task *zssq.Task, gender, categoryName, bookId string) (bookInfo *models.ZssqBookDesc, err error) {
	//通过book id获取书本详情
	bookInfoResp, err := task.GetBookInfo(bookId)
	if err != nil {
		global.Zssqlog.Errorf("获取小说详情失败 %v", err.Error())
		return
	}
	bookInfo = ZssqGetBookDescByData(bookInfoResp, gender, categoryName)
	return
}

func ZssqGetBookChapters(task *zssq.Task, bookId, bookName, author string) (chaptersAll []*models.ZssqChapter, lastUpdateTime string, err error) {
	//通过book id获取书本目录
	bookDirectory, err := task.GetBookDirectory(bookId)
	if err != nil {
		global.Zssqlog.Errorf("获取小说章节目录失败 bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}
	lastUpdateTime, _ = bookDirectory.Get("updated").TryString()
	chapters, err := bookDirectory.Get("chapters").TryJsonArray()
	if err != nil {
		global.Zssqlog.Errorf("转换小说章节目录失败  bookName=%v author=%v err=%v", bookName, author, err.Error())
		return
	}
	if len(chapters) <= 0 {
		err = fmt.Errorf("获取小说章节目录失败 bookName=%v author=%v", bookName, author)
		return
	}

	for _, chapter := range chapters {
		//之前获取目录中带有的Link就是章节内容链接
		title, _ := chapter.Get("title").TryString()
		link, _ := chapter.Get("link").TryString()
		sortStr, _ := chapter.Get("order").TryString()
		sort, _ := strconv.Atoi(sortStr)
		chapterInfo := &models.ZssqChapter{
			ChapterName: title,
			ChapterLink: link,
			Sort:        sort,
		}
		chaptersAll = append(chaptersAll, chapterInfo)
	}
	return
}

func ZssqGetChapterText(task *zssq.Task, bookId, bookName, chapterName, chapterLink string, sort int) (text string, err error) {
	if strings.Contains(chapterLink, "yuewenhttp") {
		//起点SDK的链接，未写解密方法，直接跳过采集
		var chapterContentResp *aJson.Json
		chapterContentResp, err = task.GetChapterContentByYW(chapterLink)
		if err != nil {
			global.Zssqlog.Infof("%v", err.Error())
			return
		}
		encryptStr, _ := chapterContentResp.Get("data").Get("content").TryString()
		text = zssq.DecryptChapterContentByYW(encryptStr)
		return
	}
	bookName = strings.TrimSpace(bookName)
	chapterName = strings.TrimSpace(chapterName)
	orderId := fmt.Sprintf("%v", sort)
	//通过book id 和 顺序ID 获取章节txt解密密钥
	keyResp, err := task.GetChapterCryptoKey(bookId, orderId)
	if err != nil {
		global.Zssqlog.Errorf("获取小说章节解密秘钥失败 %v", err.Error())
		return
	}
	resOk, _ := keyResp.Get("ok").TryBool()
	if !resOk {
		resMsg, _ := keyResp.Get("msg").TryString()
		global.Zssqlog.Errorf("获取小说章节内容失败 bookId=%v bookName=%v chapterName=%v err=%v", bookId, bookName, chapterName, resMsg)
		return
	}
	// "noEncrypt":"okp2TyH/ZL8Ky8sVqmO3xg==" 是key
	keys, _ := keyResp.Get("data").TryJsonArray()
	key, _ := keys[0].Get("noEncrypt").TryString()

	//fmt.Println("keyResp:", keyResp.Data)
	//return
	//通过链接采集加密的章节内容（有时候有的章节是不加密的）
	chapterContentResp, err := task.GetChapterContent(chapterLink)
	if err != nil {
		global.Zssqlog.Errorf("获取章节解密key失败 %v", chapterLink)
		return
	}
	resOk, _ = chapterContentResp.Get("ok").TryBool()
	if !resOk {
		resMsg, _ := chapterContentResp.Get("msg").TryString()
		global.Zssqlog.Errorf("获取小说章节内容失败 bookId=%v bookName=%v chapterName=%v err=%v", bookId, bookName, chapterName, resMsg)
		err = fmt.Errorf("%v", resMsg)
		return
	}
	cp, _ := chapterContentResp.Get("chapter").Get("cpContent").TryString()
	cpBytes, err := base64.StdEncoding.DecodeString(cp)
	if err != nil {
		global.Zssqlog.Errorf("cpBytes Base64 Decode Error %v", err.Error())
		return
	}
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		global.Zssqlog.Errorf("keyBytes Base64 Decode Error %v", err.Error())
		return
	}
	//解密
	textBytes, err := zssq.AesCbcDecrypt(cpBytes, keyBytes)
	if err != nil {
		global.Zssqlog.Errorf("AesCbcDecrypt Error %v", err.Error())
		return
	}
	// 将字节切片转换为字符串
	text = string(textBytes)
	return
}

func ZssqCollectChapterText(bookId, bookName, author string, task *zssq.Task, chapter *models.ZssqChapter) (textNum int, err error) {
	bookName = strings.TrimSpace(bookName)
	author = strings.TrimSpace(author)
	chapterName := chapter.ChapterName
	chapterLink := chapter.ChapterLink
	sort := chapter.Sort
	var txtFile string
	_, txtFile, err = book_service.GetChapterTxtFile(bookName, author, chapterName)
	if !utils.CheckNotExist(txtFile) {
		textNum = book_service.GetTxtNum(txtFile)
		global.Zssqlog.Infof("章节内容文件已存在 bookName=%v chapterTitle=%v txtFile=%v textNum=%v", bookName, chapterName, txtFile, textNum)
		return
	}
	text, err := ZssqGetChapterText(task, bookId, bookName, chapterName, chapterLink, sort)
	if err != nil {
		return
	}
	textNum = len([]rune(text))
	var chapterNameMd5 string
	chapterNameMd5, _, err = book_service.GetBookTxt(bookName, author, chapterName, text)
	if err != nil {
		return
	}
	log.Println(utils.GetBookMd5(bookName, author), bookName, chapterName, chapterLink, textNum, chapterNameMd5, "写入成功")
	return
}
