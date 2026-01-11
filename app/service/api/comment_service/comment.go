package comment_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/global"
	"go-novel/utils"
)

func GetTodayCommentCountByUserId(uid int64) (count int64) {
	err := global.DB.Model(models.McComment{}).Where("uid = ? and addtime >= ?", uid, utils.GetTomorrowUnix()).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentById(id int64) (comment *models.McComment, err error) {
	err = global.DB.Model(models.McComment{}).Where("id = ?", id).Find(&comment).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentCountById(commentId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McComment{}).Where("id = ?", commentId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentCountByUid(userId, agoTime int64) (count int64) {
	var err error
	db := global.DB.Model(models.McComment{})
	db = db.Where("uid = ?", userId)
	if agoTime > 0 {
		db = db.Where("addtime >= ?", agoTime)
	}
	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentTextById(id int64) (text string) {
	var err error
	err = global.DB.Model(models.McComment{}).Select("text").Where("id = ?", id).Scan(&text).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentCountByTime(userId, agoTime int64) (count int64) {
	var err error
	err = global.DB.Model(models.McComment{}).Order("addtime desc").Where("uid = ? and addtime >= ?", userId, agoTime).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentsByParentLink(parentLink string, commentId int64) (comments []*models.McComment, err error) {
	parentId := fmt.Sprintf("%v%d,", parentLink, commentId)
	err = global.DB.Model(models.McComment{}).Where("parent_link like ?", "%"+parentId+"%").Find(&comments).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCountByParentLink(commentId int64) (count int64) {
	var err error
	parentId := fmt.Sprintf(",%d,", commentId)
	err = global.DB.Model(models.McComment{}).Where("parent_link like ?", "%"+parentId+"%").Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateReplyNum(commentId, count int64) (err error) {
	err = global.DB.Model(models.McComment{}).Where("id", commentId).Update("reply_num", count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetComments() (comments []*models.McComment, err error) {
	err = global.DB.Model(models.McComment{}).Order("id desc").Find(&comments).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentsByBookId(bookId int64, limit int) (comments []*models.McComment, err error) {
	db := global.DB.Model(models.McComment{}).Order("id desc").Where("pid = 0 and bid = ?", bookId)
	if limit > 0 {
		db = db.Limit(limit)
	}
	err = db.Find(&comments).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCommentList(req *models.CommentListReq) (commentList []*models.CommentListRes, total int64, err error) {
	var list []*models.McComment
	db := global.DB.Model(&models.McComment{})
	sort := req.Sort
	if sort == utils.Hot {
		db = db.Order("praise_count desc,id desc")
	} else {
		db = db.Order("id desc")
	}

	bookId := req.BookId
	if bookId > 0 {
		db = db.Where("bid = ?", bookId)
	}

	db = db.Where("pid = 0")

	userId := req.UserId
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

	for _, comment := range list {
		user, _ := user_service.GetUserById(comment.Uid)
		com := models.CommentListRes{
			Id:          comment.Id,
			Ip:          comment.Ip,
			UserId:      user.Id,
			Nickname:    user.Nickname,
			Pic:         utils.GetFileUrl(user.Pic),
			Text:        comment.Text,
			Score:       comment.Score,
			ReplyNum:    comment.ReplyNum,
			Addtime:     comment.Addtime,
			PraiseCount: comment.PraiseCount,
		}
		if userId > 0 {
			count := GetIsPraiseByUserId(comment.Id, userId)
			if count > 0 {
				com.IsPraise = 1
			}
		}
		commentList = append(commentList, &com)
	}
	return
}

func GetCountByCommentIdAndUid(commentId, userId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McComment{}).Where("id = ? and uid = ?", commentId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteCommentByCommentId(commentId int64) (err error) {
	if commentId <= 0 {
		return
	}
	parentId := fmt.Sprintf(",%d,", commentId)
	err = global.DB.Where("id = ?", commentId).Or("parent_link like ?", "%"+parentId+"%").Delete(&models.McComment{}).Error
	return
}
