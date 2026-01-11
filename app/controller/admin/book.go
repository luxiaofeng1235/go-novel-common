package admin

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/book_service"
	"go-novel/app/service/admin/class_service"
	"go-novel/app/service/admin/tag_service"
	"go-novel/app/service/common/common_service"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"strconv"
)

type Book struct{}

// 设置推荐首页的列表数据信息
func (Book *Book) SetRecBookIndex(c *gin.Context) {
	if c.Request.Method == "POST" {
		//绑定列表数据
		var req models.CreateBookRecommandReq

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatalf("read body failed at Before,err:%s", err.Error())
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		log.Println("read body: ", string(body))

		if err := c.ShouldBind(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		res, err := book_service.SetBookRecData(&req)
		if err != nil || res == false {
			utils.Fail(c, err, "保存失败")
			return
		}
		utils.Success(c, res, "设置成功")
		return
	}
	utils.Success(c, "", "ok")
	return
}

func (Book *Book) GetBookRecList(c *gin.Context) {
	//var req models.SearchBookRecReq
	//if err := c.ShouldBind(&req); err != nil {
	//	utils.Fail(c, err, "参数绑定失败")
	//	return
	//}
	//获取推荐的活动列表数据信息
	list, err := book_service.GetBookRecList()
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.Success(c, res, "ok")
	return
}

// 获取书籍通过的数据列表信息
func (Book *Book) BookPassList(c *gin.Context) {
	var req models.BookListPassReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	//获取已经审核通过的书籍ID
	list, total, err := book_service.BookListPassSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}
	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (Book *Book) BookList(c *gin.Context) {
	var req models.BookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := book_service.BookListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	classList, err := class_service.GetClassBySex()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}

	tagList, err := tag_service.GetTagBySex()

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
		"classList":   classList,
		"tagList":     tagList,
	}
	utils.Success(c, res, "ok")
}

func (Book *Book) CreateBook(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateBookReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := book_service.CreateBook(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}
		common_service.AddSearchBookInfo(InsertId)
		utils.Success(c, InsertId, "ok")
		return
	}

	classList, err := class_service.GetClassBySex()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}

	res := gin.H{
		"classList": classList,
	}
	utils.Success(c, res, "ok")
}

func (Book *Book) UpdateBook(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateBookReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := book_service.UpdateBook(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}
		//删除索引中的数据
		if req.Status == 0 {
			common_service.XunDelBookById("delIndex", req.BookId) //删除索引
		} else {
			common_service.XunAddBookInfo("addIndex", req.BookId, req.BookName, req.Author) //添加索引
		}
		utils.Success(c, "", "ok")
		return
	}

	bookId, _ := strconv.Atoi(c.Query("id"))
	bookInfo, err := book_service.GetBookById(int64(bookId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}
	if bookInfo.Pic != "" {
		bookInfo.Pic = utils.GetAdminFileUrl(bookInfo.Pic)
	}
	classList, err := class_service.GetClassList()
	if err != nil {
		utils.Fail(c, err, "获取分类列表失败")
		return
	}

	if len(classList) > 0 {
		for _, val := range classList {
			text := val.ClassName
			if val.BookType == 1 {
				text = fmt.Sprintf("男生-%v", text)
			} else if val.BookType == 2 {
				text = fmt.Sprintf("女生-%v", text)
			}
			val.ClassName = text
		}
	}

	res := gin.H{
		"bookInfo":  bookInfo,
		"classList": classList,
	}
	utils.Success(c, res, "ok")
}

func (Book *Book) DelBook(c *gin.Context) {
	var req models.DeleteBookReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := book_service.DeleteBook(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}
	//删除对应的索引配置信息
	common_service.XunDelBookById("delIndex", req.BookId)
	utils.Success(c, "", "删除信息成功")
	return
}

func (Book *Book) DetailBook(c *gin.Context) {
	bookId, _ := strconv.Atoi(c.Query("id"))
	bookInfo, err := book_service.GetBookById(int64(bookId))
	if err != nil {
		utils.Fail(c, nil, "获取小说数据失败")
		return
	}

	res := gin.H{
		"bookInfo": bookInfo,
		"bookMd5":  utils.GetBookMd5(bookInfo.BookName, bookInfo.Author),
	}
	utils.Success(c, res, "ok")
}

func (Book *Book) GetRandHit(c *gin.Context) {
	hits, hitsDay, hitsWeek, hitsMonth, shits, score, readCount, searchCount := utils.GetRandNumBookHits()

	res := gin.H{
		"hits":        hits,
		"hitsDay":     hitsDay,
		"hitsWeek":    hitsWeek,
		"hitsMonth":   hitsMonth,
		"shits":       shits,
		"score":       score,
		"readCount":   readCount,
		"searchCount": searchCount,
	}
	utils.Success(c, res, "ok")
}
