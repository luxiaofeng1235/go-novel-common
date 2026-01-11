package common_service

import (
	"fmt"
	"go-novel/app/service/common/book_service"
	"go-novel/config"
	"go-novel/global"
	"go-novel/utils"
	"log"
)

var hostNameOnline = "http://192.168.10.17:5678" //内网接口地址信息
var hostNameOutline = "http://103.36.91.35:5678" //外网接口地址
var dbName = "novel"                             //搜索数据库

// 选择线上的地址信息
func SwitchGetUrl() (url string) {
	env := config.GetString("server.env")
	if env == utils.Prod {
		url = hostNameOnline //线上使用内网地址
	} else {
		url = hostNameOutline //使用外网地址
	}
	return
}

// 定义书籍搜素的字段类型
type BookSearch struct {
	Query string `json:"query"`
	Page  int64  `json:"page"`
	Limit int64  `json:"limit"`
	Order string `json:"order"`
}

type Document struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Document struct {
		Number   int64  `json:"number"`    //索引位置
		Title    string `json:"title"`     //标题
		BookId   int64  `json:"book_id"`   //小说id
		BookName string `json:"book_name"` //书名
		Author   string `json:"author"`    //作者
	} `json:"document"`
	Score int      `json:"score"`
	Keys  []string `json:"keys"`
}

type Response struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Data    struct {
		Time      float64    `json:"time"`
		Total     int        `json:"total"`
		PageCount int        `json:"pageCount"`
		Page      int        `json:"page"`
		Limit     int        `json:"limit"`
		Documents []Document `json:"documents"`
	} `json:"data"`
}

// 获取搜索引擎里的数据方法
func GoqueryByBookName(title string, page, limit int64, scoreExp, order string) (content string) {
	if title == "" {
		return
	}
	url := SwitchGetUrl()
	//获取需要组装的url
	apiUrl := fmt.Sprintf("%v/api/query?database=%v", url, dbName)
	mapData := make(map[string]interface{})
	mapData["query"] = title
	mapData["page"] = page
	mapData["limit"] = limit
	mapData["order"] = order
	mapData["scoreExp"] = scoreExp
	log.Printf("param = %v", mapData)
	//获取请求关联的字段内容信息
	result := utils.GetPostData(apiUrl, mapData)
	return result
}

// 添加索引数据
func AddSearchBookInfo(bookId int64) (str string) {
	if bookId <= 0 {
		return
	}
	bookInfo, err := book_service.GetBookById(bookId)
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//生成添加的切片数据进行索引添加同步
	nestedData := make(map[string]interface{})
	nestedData["id"] = bookInfo.Id
	var text string
	if bookInfo.BookName != "" && bookInfo.Author != "" {
		text = fmt.Sprintf("%v@%v", bookInfo.BookName, bookInfo.Author)
	} else {
		text = bookInfo.BookName
	}
	nestedData["text"] = text
	document := make(map[string]interface{})
	document["book_id"] = bookInfo.Id
	document["book_name"] = bookInfo.BookName
	document["author"] = bookInfo.Author
	nestedData["document"] = document
	//jsonbyte, err := json.Marshal(nestedData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(jsonbyte))
	log.Printf("add search param = %v", nestedData)
	url := SwitchGetUrl()
	apiUrl := fmt.Sprintf("%v/api/index?database=%v", url, dbName)
	result := utils.GetPostData(apiUrl, nestedData)
	return result
}

// 删除的对应引擎数据
func DelSearchDataById(bookId int64) (str string) {
	if bookId == 0 {
		return
	}
	url := SwitchGetUrl()
	apiUrl := fmt.Sprintf("%v/api/index/remove?database=%v", url, dbName)
	mapData := make(map[string]interface{})
	mapData["id"] = bookId
	log.Printf("delete search param = %v", mapData)
	result := utils.GetPostData(apiUrl, mapData)
	return result
}
