package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/comment_service"
	"go-novel/utils"
	"strconv"
)

type Comment struct{}

func (comment *Comment) CommentList(c *gin.Context) {
	var req models.CommentListSearchReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := comment_service.CommentListSearch(&req)
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

func (comment *Comment) SonComments(c *gin.Context) {
	commentId, _ := strconv.Atoi(c.Query("id"))
	if commentId <= 0 {
		utils.Fail(c, nil, "评论id不能为空")
		return
	}

	list, err := comment_service.GetChildCommentsByCommentId(int64(commentId))
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	res := gin.H{
		"list": list,
	}
	utils.Success(c, res, "ok")
}

func (comment *Comment) UpdateComment(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateCommentReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := comment_service.UpdateComment(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}
		utils.Success(c, "", "ok")
		return
	}

	commentId, _ := strconv.Atoi(c.Query("id"))
	commentInfo, err := comment_service.GetCommentById(int64(commentId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"commentInfo": commentInfo,
	}
	utils.Success(c, res, "ok")
}

func (comment *Comment) DelComment(c *gin.Context) {
	var req models.DeleteCommentReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := comment_service.DeleteComment(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}
