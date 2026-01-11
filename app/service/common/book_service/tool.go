package book_service

import (
	"fmt"
	"go-novel/global"
	"go-novel/utils"
	"html"
	"io/ioutil"
	"strings"
)

func GetChapterTxtFile(bookName, author, chapterName string) (chapterNameMd5, txtFile string, err error) {
	bookName = strings.TrimSpace(bookName)
	var uploadBookTextPath string
	//uploadBookTextPath, err = setting_service.GetValueByName(utils.UploadBookTextPath)
	//if err != nil {
	//	err = fmt.Errorf("获取小说内容目录失败 uploadBookTextPath=%v", uploadBookTextPath)
	//	return
	//}
	uploadBookTextPath = "/data/txt/"
	txtDir := fmt.Sprintf("%v%v", uploadBookTextPath, utils.GetBookMd5(bookName, author))
	chapterNameMd5 = utils.GetChapterMd5(chapterName)
	txtFile = fmt.Sprintf("%v/%v", txtDir, fmt.Sprintf("%v.txt", chapterNameMd5))
	return
}

func GetTxtNum(txtFile string) (textNum int) {
	// 读取文件内容
	content, err := ioutil.ReadFile(txtFile)
	if err != nil {
		global.Collectlog.Errorf("无法读取文件章节文件 txtFile=%v err=%v", txtFile, err.Error())
		return
	}
	// 将文件内容转换为字符串
	text := string(content)
	// 统计字符数
	textNum = len([]rune(text))
	return
}

func GetBookTxt(bookName, author, chapterName, text string) (chapterNameMd5, content string, err error) {
	if chapterName == "" {
		err = fmt.Errorf("%v", "章节名称错误")
		return
	}
	var txtFile string
	chapterNameMd5, txtFile, err = GetChapterTxtFile(bookName, author, chapterName)
	if text != "" {
		// 写入文本内容到文件
		//text = html.EscapeString(strings.TrimSpace(text))
		text = utils.ReplaceText(text)
		err = utils.WriteFile(txtFile, text)
		content = text
		return
	} else {
		// 读取文件内容
		if utils.CheckNotExist(txtFile) {
			err = fmt.Errorf("%v", "内容不存在")
			return
		}
		var conByte []byte
		conByte, err = ioutil.ReadFile(txtFile)
		if err != nil {
			return
		}
		content = string(conByte)
		content = html.UnescapeString(content)
		content = utils.ReadTextReplace(content)
		return
	}
	return
}
