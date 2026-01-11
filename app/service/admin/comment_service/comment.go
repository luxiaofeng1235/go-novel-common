package comment_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/comment_service"
	"go-novel/app/service/api/user_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetCommentById(id int64) (comment *models.McComment, err error) {
	err = global.DB.Model(models.McComment{}).Where("id", id).First(&comment).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

//func GetChildCommentsByCommentId(commentId int64) (comments []*models.McComment, err error) {
//	parentId := fmt.Sprintf(",%d,", commentId)
//	err = global.DB.Model(models.McComment{}).Where("parent_link like ?", "%"+parentId+"%").Find(&comments).Error
//	if err != nil {
//		global.Sqllog.Errorf("%v", err.Error())
//		return
//	}
//	return
//}

func GetChildCommentsByCommentId(commentId int64) (replyList []*models.CommentReplyList, err error) {
	var comment *models.McComment
	comment, err = GetCommentById(commentId)
	if err != nil {
		return
	}
	if comment.Id <= 0 {
		err = fmt.Errorf("%v", "该评论不存在")
		return
	}

	var list []*models.McComment
	db := global.DB.Model(&models.McComment{}).Order("id DESC")

	parentLink := fmt.Sprintf("%v%d,", comment.ParentLink, commentId)
	db = db.Where("parent_link like ?", parentLink+"%")
	err = db.Find(&list).Error
	if len(list) <= 0 {
		return
	}
	for _, com := range list {
		user, _ := user_service.GetUserById(com.Uid)
		if user.Id <= 0 {
			continue
		}
		var parentUser *models.McUser
		parentUser, err = user_service.GetUserById(com.ReplyUid)
		if parentUser.Id <= 0 {
			continue
		}
		reply := &models.CommentReplyList{
			Id:             com.Id,
			PraiseCount:    com.PraiseCount,
			Text:           com.Text,
			Uid:            com.Uid,
			Nickname:       user.Nickname,
			Pic:            utils.GetFileUrl(user.Pic),
			Addtime:        com.Addtime,
			ParentId:       com.Pid,
			ParentNickname: parentUser.Nickname,
			ParentText:     comment_service.GetCommentTextById(com.Pid),
			Ip:             com.Ip,
			City:           com.City,
			Status:         com.Status,
		}
		replyList = append(replyList, reply)
	}
	return
}

func CommentListSearch(req *models.CommentListSearchReq) (list []*models.McComment, total int64, err error) {
	db := global.DB.Model(&models.McComment{}).Order("id desc")
	db = db.Where("pid = 0")

	bookId := strings.TrimSpace(req.BookId)
	if bookId != "" {
		db = db.Where("bid = ?", bookId)
	}

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	text := strings.TrimSpace(req.Text)
	if text != "" {
		db = db.Where("text LIKE ?", "%"+text+"%")
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

func UpdateComment(req *models.UpdateCommentReq) (res bool, err error) {
	id := req.CommentId
	praiseCount := req.PraiseCount
	status := req.Status
	text := strings.TrimSpace(req.Text)
	if text == "" {
		err = fmt.Errorf("%v", "评论内容不能为空")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["praise_count"] = praiseCount
	mapData["text"] = text
	mapData["status"] = status
	if err = global.DB.Model(models.McComment{}).Debug().Where("id", id).Updates(&mapData).Error; err != nil {
		return
	}
	return true, nil
}

func DeleteComment(req *models.DeleteCommentReq) (res bool, err error) {
	ids := req.CommentIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McComment{}).Error
		if err != nil {
			return
		}

	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
