package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/chapter_service"
	"go-novel/app/service/common/book_service"
	"go-novel/utils"
	"strconv"
)

type Chapter struct{}

func (chapter *Chapter) ChapterList(c *gin.Context) {
	var req models.ChapterListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := chapter_service.ChapterListSearch(&req)
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

func (chapter *Chapter) CreateChapter(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateChapterReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := chapter_service.CreateChapter(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}
	bookId := utils.FormatInt64(c.Query("book_id"))
	if bookId <= 0 {
		utils.Fail(c, nil, "小说id为空")
		return
	}
	bookName := c.Query("book_name")
	author := c.Query("author")
	lastSort := chapter_service.GetSortLast(bookName, author)
	res := gin.H{
		"lastSort": lastSort,
	}
	utils.Success(c, res, "ok")
}

func (chapter *Chapter) UpdateChapter(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateChapterReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := chapter_service.UpdateChapter(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	chapterId, _ := strconv.Atoi(c.Query("chapter_id"))
	bookName := c.Query("book_name")
	author := c.Query("author")
	chapterInfo, err := chapter_service.GetChapterById(bookName, author, int64(chapterId))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	var chapterText string
	_, chapterText, _ = book_service.GetBookTxt(bookName, author, chapterInfo.ChapterName, "")
	chapterInfo.TextNum = len([]rune(chapterText))
	res := gin.H{
		"chapterInfo":    chapterInfo,
		"chapterText":    chapterText,
		"ChapterNameMd5": utils.GetChapterMd5(chapterInfo.ChapterName),
		"bookMd5":        utils.GetBookMd5(bookName, author),
	}
	utils.Success(c, res, "ok")
}

func (chapter *Chapter) DelChapter(c *gin.Context) {
	var req models.DeleteChapterReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := chapter_service.DeleteChapter(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
