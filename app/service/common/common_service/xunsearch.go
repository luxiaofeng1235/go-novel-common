package common_service

import (
	"fmt"
	"go-novel/config"
	"go-novel/utils"
	"log"
)

//迅搜的API接口实现

var oneLineUrl = "http://192.168.10.16:9205/api/search.php" //内网接口地址信息
var offlineUrl = "http://103.36.91.36:9205/api/search.php"  //外网接口地址

type XunResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Book `json:"data"`
}

type Book struct {
	BookID   int64  `json:"book_id"`
	BookName string `json:"book_name"`
	Author   string `json:"author"`
	Chrono   int64  `json:"chrono"`
}

// 选择线上的地址信息
func SwitchOnlineUrl() (url string) {
	env := config.GetString("server.env")
	if env == utils.Prod {
		url = oneLineUrl //线上使用内网地址
	} else {
		url = offlineUrl //使用外网地址
	}
	return
}

// 迅搜的搜索接口集成
func XunSearchByBookName(req string, book_name string) (content string) {
	if book_name == "" {
		return
	}
	apiUrl := SwitchOnlineUrl()
	mapData := make(map[string]interface{})
	mapData["req"] = req
	mapData["book_name"] = book_name
	log.Printf("param = %v", mapData)
	fmt.Println(apiUrl)
	//获取请求关联的字段内容信息
	result := utils.GetPostData(apiUrl, mapData)
	return result
}

// 删除数据
func XunDelBookById(req string, bookId int64) (str string) {
	if bookId <= 0 {
		return
	}
	apiUrl := SwitchOnlineUrl()
	mapData := make(map[string]interface{})
	mapData["req"] = req
	mapData["id"] = bookId
	log.Printf("param = %v", mapData)
	result := utils.GetPostData(apiUrl, mapData)
	return result
}

// 自动同步迅搜的接口数据
func XunAddBookInfo(req string, bookId int64, bookName, Author string) (content string) {
	if bookId == 0 {
		return
	}
	apiUrl := SwitchOnlineUrl()
	mapData := make(map[string]interface{})
	mapData["req"] = req
	mapData["book_id"] = bookId
	mapData["book_name"] = bookName
	mapData["author"] = Author
	log.Printf("param = %v", mapData)
	result := utils.GetPostData(apiUrl, mapData)
	return result
}
