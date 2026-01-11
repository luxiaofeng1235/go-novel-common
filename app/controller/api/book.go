package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/book_service"
	"go-novel/utils"
)

type Book struct{}

// 书籍评论的列表
func (book *Book) BookCommentUserList(c *gin.Context) {
	var req models.BookCommnetUserReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取用户列表信息
	userList, err := book_service.GetCommentUserList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "数据异常")
		return
	}
	res := gin.H{
		"userList": userList,
	}
	utils.Success(c, res, "用户列表")
}

// 书籍排行列表
func (book *Book) BookCommentRankList(c *gin.Context) {
	var req models.BookCommentListReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取书评列表
	bookCommentList, err := book_service.GetBookCommnetLists(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "获取失败")
		return
	}
	//返回值设置
	res := gin.H{
		"rankBookList": bookCommentList,
	}
	utils.SuccessEncrypt(c, res, "书评排行列表")

}

// 精彩书评列表
func (book *Book) BookWonderfulList(c *gin.Context) {
	var req models.BookCommnetUserReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	rankCommentList, err := book_service.GetWonderfulCommentList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	//获取排行列表
	rankList := gin.H{
		"rankCommentList": rankCommentList,
	}
	utils.SuccessEncrypt(c, rankList, "精彩书评排行")
}

// 根据用户名获取书评信息
func (book *Book) BookCommentByUserId(c *gin.Context) {
	var req models.BookCommentSingeUidReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId := int64(req.UserId)
	if userId == 0 {
		utils.FailEncrypt(c, nil, "参数错误")
		return
	}
	//获取用户关联的书评信息
	commentRes := book_service.GetBookCommentByUid(userId)
	res := gin.H{
		"comment_list": commentRes,
	}
	utils.SuccessEncrypt(c, res, "用户书评信息")
}

// 书评基本信息-包含评论+基础信息
func (book *Book) BookCommentInfo(c *gin.Context) {
	var req models.BookCommentInfoReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取书籍详情
	bookDetail, err := book_service.GetBookDetailInfo(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "数据异常")
		return
	}
	res := gin.H{
		"bookInfo": bookDetail,
	}
	utils.SuccessEncrypt(c, res, "书评详情")
}

// 根据分类查书籍
func (book *Book) GetCateBookList(c *gin.Context) {
	var req models.ApiCateBookReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//查询分类对应的书籍
	list, err := book_service.BookCateListRes(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) List(c *gin.Context) {
	var req models.ApiBookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, nums, page, size, err := book_service.BookListSearch(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list":  list,
		"count": nums,
		"page":  page,
		"size":  size,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) Info(c *gin.Context) {
	var req models.BookInfoReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	rbook, err := book_service.Info(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, rbook, "ok")
}

func (book *Book) GetHighScoreBook(c *gin.Context) {
	var req models.BookHighScoreReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	if req.Bid > 0 {
		info, _ := book_service.GetBookById(req.Bid)
		if info.Cid > 0 {
			req.ClassId = info.Cid
		}
		if info.Tid > 0 {
			req.TagId = info.Tid
		}
	}
	books, err := book_service.GetHighScoreBook(req.Bid, req.ClassId, req.TagId, req.UserId, req.Size, req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, books, "ok")
}

func (book *Book) Chapter(c *gin.Context) {
	var req models.ChapterReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	chapterRes, bookUrl, chapterTextReg, err := book_service.Chapter(req.BookId, req.SourceId, req.Sort)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"chapters":       chapterRes,
		"bookUrl":        bookUrl,
		"chapterTextReg": chapterTextReg,
	}

	utils.SuccessEncrypt(c, res, "获取章节列表成功")
}

func (book *Book) Read(c *gin.Context) {
	var req models.ChapterReadReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	readRes, err := book_service.ChapterRead(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, readRes, "ok")
}

func (book *Book) RankList(c *gin.Context) {
	var req models.ApiRankListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark

	list, err := book_service.GetRankList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetSectionForYouRec(c *gin.Context) {
	var req models.SectionForYouRecReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}

	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark

	req.Ip = utils.RemoteIp(c)
	//fmt.Println(req.Ip, req.DeviceType, req.PackageName)
	list, err := book_service.GetSectionForYouRec(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetSectionHighScore(c *gin.Context) {
	var req models.SectionHighScoreReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark

	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, err := book_service.GetSectionHighScore(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetSectionEnd(c *gin.Context) {
	var req models.SectionEndReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	req.Ip = utils.RemoteIp(c)
	list, err := book_service.GetSectionEnd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetSectionNew(c *gin.Context) {
	var req models.SectionNewReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark

	fmt.Println(req.Ip, req.DeviceType, req.PackageName)
	list, err := book_service.GetSectionNew(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetTags(c *gin.Context) {
	var req models.GetTagsReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	tags, err := book_service.GetTags(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, tags, "ok")
}

func (book *Book) TeenZoneList(c *gin.Context) {
	var req models.TeenZoneListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	list, err := book_service.GetTeenZoneList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetNewBookRec(c *gin.Context) {
	var req models.GetTagsReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	res, err := book_service.GetNewBookRec(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetNewBookList(c *gin.Context) {
	var req models.GetNewBookListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取客户端IP
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息(渠道信息)
	mark := utils.GetDeviceQdhInfo(c)

	req.Mark = mark
	list, total, err := book_service.GetNewBookList(c, &req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list":  list,
		"total": total,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) TodayUpdateBooks(c *gin.Context) {
	var req models.TodayUpdateBooksReq

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	fmt.Println(req.PackageName, req.DeviceType)
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	list, total, err := book_service.TodayUpdateBooks(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list":  list,
		"total": total,
	}
	//fmt.Println(rediskey.GetTodayBookKey())
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) GetHotCount(c *gin.Context) {
	var req models.HotBookCountReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	endCount, rankCount, newCount := book_service.GetHotCount(&req)
	res := gin.H{
		"endCount":  endCount,
		"rankCount": rankCount,
		"newCount":  newCount,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (book *Book) Test(c *gin.Context) {
	//bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil {
	//	global.Paylog.Infof("read request body failed,err =%s", err)
	//	return
	//}
	//
	//log.Printf("bodyBytes=%+v", string(bodyBytes))
	var req models.TodayUpdateBooksReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	utils.SuccessEncrypt(c, req, "ok")
}
