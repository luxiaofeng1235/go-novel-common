package chapter_service

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/utils"
	"io/ioutil"
	"os"
)

func GetBookJsonq(chapterFile string) (gq *gojsonq.JSONQ, err error) {
	if chapterFile == "" {
		err = fmt.Errorf("%v", "小说章节json文件地址为空")
		return
	}
	gq = gojsonq.New().File(chapterFile)
	return
}

func GetJsonqByBookName(bookName, author string) (gq *gojsonq.JSONQ, chapterFile string, err error) {
	chapterFile, err = GetChapterFile(bookName, author)
	if err != nil {
		return
	}
	if chapterFile == "" {
		err = fmt.Errorf("%v", "章节文件不存在")
		return
	}
	gq, err = GetBookJsonq(chapterFile)
	return
}

func GetChapterFile(bookName, author string) (chapterFile string, err error) {
	var uploadBookChapterPath string
	//uploadBookChapterPath, err = setting_service.GetValueByName(utils.UploadBookChapterPath)
	//if err != nil {
	//	err = fmt.Errorf("获取小说内容目录失败 uploadBookChapterPath=%v", uploadBookChapterPath)
	//	return
	//}
	//判断书名和作者
	if bookName == "" || author == "" {
		err = fmt.Errorf("%v", "小说标题或者作者为空")
		return
	}
	bookFile := utils.GetBookMd5(bookName, author) //获取小说和作者的加密值
	//目录组成结构：小说章节+MD5的前两个字符作为存储的路径信息
	//uploadBookChapterPath = "/data/chapter/"
	//按照指定的格式进行目录格式：/data/chapter/+md5的前两个字符+md5.json
	uploadBookChapterPath = "/data/chapter/" + bookFile[0:2] + "/"
	//log.Printf("当前存储章节目录的json 路径：%s\n", uploadBookChapterPath)
	err = utils.IsNotExistMkDir(uploadBookChapterPath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}
	chapterFile = fmt.Sprintf("%v%v.json", uploadBookChapterPath, bookFile)
	if utils.CheckNotExist(chapterFile) {
		_, err = os.Create(chapterFile)
		if err != nil {
			return
		}
	}
	return
}

func GetTxtDir(bookName, author string) (textDir string, err error) {
	var uploadBookTextPath string
	uploadBookTextPath, err = setting_service.GetValueByName(utils.UploadBookTextPath)
	if err != nil {
		err = fmt.Errorf("获取小说内容目录失败 uploadBookTextPath=%v", uploadBookTextPath)
		return
	}
	err = utils.IsNotExistMkDir(uploadBookTextPath)
	if err != nil {
		err = fmt.Errorf("%v", "创建目录失败")
		return
	}
	bookFile := utils.GetBookMd5(bookName, author)
	textDir = fmt.Sprintf("%v%v", uploadBookTextPath, bookFile)
	return
}

func GetChaptersByFile(chapterFile string) (chapterAll []*models.McBookChapter, err error) {
	if chapterFile == "" {
		err = fmt.Errorf("%v", "小说章节json文件地址为空")
		return
	}
	jsonData, err := ioutil.ReadFile(chapterFile)
	if err != nil {
		err = fmt.Errorf("读取文件错误 %v", err.Error())
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if len(jsonData) > 0 {
		err = json.Unmarshal(jsonData, &chapterAll)
		if err != nil {
			err = fmt.Errorf("解析小说章节数据错误 %v", err.Error())
			return
		}
	}
	return
}

func GetChapterNamesByFile(chapterFile string) (chapterAll []*models.McBookChapter, chapterNames []string, err error) {
	if chapterFile == "" {
		err = fmt.Errorf("%v", "小说章节json文件地址为空")
		return
	}
	jsonData, err := ioutil.ReadFile(chapterFile)
	if err != nil {
		err = fmt.Errorf("读取文件错误 %v", err.Error())
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if len(jsonData) > 0 {
		err = json.Unmarshal(jsonData, &chapterAll)
		if err != nil {
			err = fmt.Errorf("解析小说章节数据错误 %v", err.Error())
			return
		}
	}
	if len(chapterAll) <= 0 {
		return
	}
	for _, chapter := range chapterAll {
		chapterNames = append(chapterNames, chapter.ChapterName)
	}
	return
}
