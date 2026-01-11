package comment_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/task_service"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/common/notify_service"
	"go-novel/global"
	"go-novel/utils"
)

func GetCommentRes(req *models.CommentListReq) (commentRes *models.CommentRes, err error) {
	bookId := req.BookId
	var commentList []*models.CommentListRes
	var total int64
	commentList, total, err = GetCommentList(req)
	if err != nil {
		return
	}
	score := GetBookScoreById(bookId)
	commentRes = &models.CommentRes{
		Comments: commentList,
		Total:    total,
		Score:    score,
	}
	return
}

func CommentAdd(req *models.CommentAddReq) (commentAddRes *models.CommentAddRes, err error) {
	bookId := req.BookId
	text := req.Text
	parentId := req.Parentid
	userId := req.UserId
	replyUid := req.ReplyUid
	score := req.Score
	machine := req.Machine
	var parentLink string
	ip := req.Ip
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	if text == "" {
		err = fmt.Errorf("%v", "评论内容不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	var parentText string
	if parentId > 0 {
		pLink := getParentLinkByCommentId(parentId)
		parentLink = fmt.Sprintf("%v%v,", pLink, parentId)
		parentText = getTextByCommentId(parentId)
	} else {
		parentLink = ","
	}
	var count int64
	count = GetCommentCountByUid(userId, utils.GetAgoDayUnix(1))
	if count > utils.Pl_Add_Num {
		err = fmt.Errorf("%v", "您今天评论数已达上限，明天再来吧")
		return
	}
	cityName := utils.GetIpDbNameByIp(ip)
	comment := models.McComment{
		Text:        text,
		ParentLink:  parentLink,
		Pid:         parentId,
		Bid:         bookId,
		Uid:         userId,
		ReplyUid:    replyUid,
		Machine:     machine,
		Ip:          ip,
		PraiseCount: 0,
		ReplyNum:    0,
		Score:       score,
		Status:      1,
		City:        cityName,
		Addtime:     utils.GetUnix(),
	}
	if err = global.DB.Create(&comment).Error; err != nil {
		return
	}

	var user *models.McUser
	user, err = user_service.GetUserById(userId)
	if err != nil {
		err = fmt.Errorf("%v", "用户不存在")
		return
	}
	if parentId > 0 {
		_ = UpdateReplyNum(parentId, GetCountByParentLink(parentId))
	}
	var taskId int64 = 4
	task, err := task_service.GetTaskById(taskId)
	if err != nil {
		return
	}
	err = task_service.CompleteTask(task, userId)
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		return
	}

	//评论通知
	if userId != replyUid && replyUid > 0 {
		_ = notify_service.SendNotify(utils.Comment, parentText, user.Pic, userId, replyUid, fmt.Sprintf("%v 评论了您的回复", user.Nickname), fmt.Sprintf("%v", comment.Text), comment.Id)
	}

	commentAddRes = new(models.CommentAddRes)
	if user.Nickname != "" {
		commentAddRes.Nickname = user.Nickname
	} else {
		commentAddRes.Nickname = user.Username
	}
	if user.Pic != "" {
		commentAddRes.Pic = utils.GetFileUrl(user.Pic)
	}
	commentAddRes.Pid = parentId
	commentAddRes.Text = text
	commentAddRes.Uid = userId
	commentAddRes.ReplyUid = replyUid
	commentAddRes.ReplyNum = 0
	commentAddRes.PraiseCount = 0
	commentAddRes.IsPraise = 0
	commentAddRes.Addtime = user.Addtime
	return
}

func GetCommentReplyList(req *models.CommentReplyListReq) (commentRes *models.CommentReplyRes, total int64, err error) {
	commentId := req.CommentId
	userId := req.UserId
	var comment *models.McComment
	comment, err = GetCommentById(commentId)
	if err != nil {
		return
	}
	if comment.Id <= 0 {
		err = fmt.Errorf("%v", "该评论不存在")
		return
	}
	commentUid := comment.Uid
	bookId := comment.Bid
	var nickName string
	user, _ := getUserById(commentUid)
	if user.Nickname != "" {
		nickName = user.Nickname
	} else {
		nickName = "匿名"
	}

	comUser := &models.CommentReplyUserRes{
		Uid:      commentUid,
		Nickname: nickName,
		Pic:      utils.GetFileUrl(user.Pic),
	}

	book, _ := getBookById(bookId)
	comBook := &models.CommentReplyBookRes{
		BookId:   book.Id,
		BookName: book.BookName,
		Author:   book.Author,
		Pic:      utils.GetFileUrl(book.Pic),
	}
	var isCommentPraise int
	count := GetIsPraiseByUserId(commentId, userId)
	if count > 0 {
		isCommentPraise = 1
	}
	commentRes = &models.CommentReplyRes{
		Id:          comment.Id,
		Text:        comment.Text,
		City:        comment.City,
		ReplyNum:    comment.ReplyNum,
		IsPraise:    isCommentPraise,
		PraiseCount: comment.PraiseCount,
		Addtime:     comment.Addtime,
		User:        comUser,
		Book:        comBook,
	}

	if userId > 0 {
		count = GetFollowCountByUid(userId, commentUid)
		if count > 0 {
			commentRes.IsFollow = 1
		}
	}
	var list []*models.McComment
	db := global.DB.Model(&models.McComment{}).Order("id DESC")

	parentLink := fmt.Sprintf("%v%d,", comment.ParentLink, commentId)
	db = db.Where("parent_link like ?", parentLink+"%")

	commentIds := GetReportCommentIdsByUserId(userId)
	if len(commentIds) > 0 {
		db = db.Where("id not in ?", commentIds)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}
	var replyList []*models.CommentReplyListRes
	for _, com := range list {
		var isPraise int
		count = GetIsPraiseByUserId(com.Id, userId)
		if count > 0 {
			isPraise = 1
		}
		user, _ = getUserById(com.Uid)
		if user.Id <= 0 {
			continue
		}
		var parentUser *models.McUser
		parentUser, err = getUserById(com.ReplyUid)
		if parentUser.Id <= 0 {
			continue
		}
		reply := &models.CommentReplyListRes{
			Id:             com.Id,
			IsPraise:       isPraise,
			PraiseCount:    com.PraiseCount,
			Text:           com.Text,
			City:           com.City,
			Uid:            com.Uid,
			Nickname:       user.Nickname,
			Pic:            utils.GetFileUrl(user.Pic),
			Addtime:        com.Addtime,
			ParentId:       com.Pid,
			ParentNickname: parentUser.Nickname,
			ParentText:     GetCommentTextById(com.Pid),
		}
		replyList = append(replyList, reply)
	}
	commentRes.ReplyList = replyList
	return
}

func CommentDel(req *models.CommentDelReq) (err error) {
	commentId := req.CommentId
	userId := req.UserId
	if commentId <= 0 {
		err = fmt.Errorf("%v", "评论ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	var count int64
	count = GetCountByCommentIdAndUid(commentId, userId)
	if count <= 0 {
		err = fmt.Errorf("%v", "只可以删除自己的评论")
		return
	}
	err = DeleteCommentByCommentId(commentId)
	return
}

func PraiseUser(req *models.PraiseUserReq) (err error) {
	praiseType := req.PraiseType
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	commentId := req.Commentid
	if commentId <= 0 {
		err = fmt.Errorf("%v", "点赞的评论ID不能为空")
		return
	}
	var count int64
	count = GetPraiseCountByUid(userId, commentId)
	var comment *models.McComment
	comment, err = GetCommentById(commentId)
	if comment.Id <= 0 {
		err = fmt.Errorf("%v", "点赞评论不存在")
		return
	}
	commentText := comment.Text
	if praiseType <= 0 {
		if count <= 0 {
			return
		}
		err = DeletePraiseUid(userId, commentId)
		if err != nil {
			return
		}
	} else {
		if count > 0 {
			err = fmt.Errorf("%v", "已经点赞啦~")
			return
		}
		receiveUid := comment.Uid
		var commentType int = 1
		if comment.Pid > 0 {
			commentType = 2
		}

		var user *models.McUser
		user, err = user_service.GetUserById(userId)
		if err != nil {
			err = fmt.Errorf("%v", "用户不存在")
			return
		}

		var praise models.McCommentPraise
		praise.ByUid = receiveUid
		praise.Uid = userId
		praise.Cid = commentId
		praise.Type = commentType
		praise.Addtime = utils.GetUnix()
		if err = global.DB.Create(&praise).Error; err != nil {
			return
		}

		if userId != receiveUid {
			var typeText string = "书评"
			if comment.Pid > 0 {
				typeText = "回复"
			}
			//点赞通知
			_ = notify_service.SendNotify(utils.Praise, commentText, user.Pic, userId, receiveUid, fmt.Sprintf("%v 赞了您的%v", user.Nickname, typeText), comment.Text, commentId)
		}
	}

	err = UpdatePraiseByCommentId(commentId, praiseType)
	if err != nil {
		return
	}

	return
}

func StarGroup(req *models.StarGroupReq) (results []*models.StarGroupRes, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说ID不能为空")
		return
	}
	sql := `SELECT
			  FLOOR((score - 1) / 2) + 1 AS star,
			  COUNT(*) AS count
			FROM
			  mc_comment
			WHERE
			  score != 0
			GROUP BY
			  star
			ORDER BY
			  star`
	err = global.DB.Raw(sql).Scan(&results).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func CommentReport(req *models.CommentReportReq) (err error) {
	commentId := req.CommentId
	userId := req.UserId
	if commentId <= 0 {
		err = fmt.Errorf("%v", "举报评论id不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	var count int64
	err = global.DB.Model(models.McCommentReport{}).Where("cid = ? and uid = ?", commentId, userId).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		err = fmt.Errorf("%v", "已举报该评论,请勿重复举报")
		return
	}
	report := models.McCommentReport{
		Cid:     commentId,
		Uid:     userId,
		Addtime: utils.GetUnix(),
	}
	if err = global.DB.Create(&report).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
