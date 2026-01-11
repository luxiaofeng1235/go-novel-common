package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/task_service"
	"go-novel/app/service/admin/user_service"
	"go-novel/utils"
	"strconv"
)

type User struct{}

// 获取用户列表
func (user *User) UserList(c *gin.Context) {
	var req models.UserListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	// 获取
	list, total, err := user_service.UserListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取用户列表失败")
		return
	}
	utils.Success(c, gin.H{"list": list, "total": total, "currentPage": req.PageNum}, "ok")
}

func (user *User) DetailUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("id"))
	userInfo, err := user_service.GetUserById(int64(userId))
	if err != nil {
		utils.Fail(c, nil, "获取用户数据失败")
		return
	}

	res := gin.H{
		"userInfo": userInfo,
	}
	utils.Success(c, res, "ok")
}

func (user *User) UpdateUser(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateUserReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := user_service.UpdateUser(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	userId, _ := strconv.Atoi(c.Query("id"))
	userInfo, err := user_service.GetUserById(int64(userId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"userInfo": userInfo,
	}
	utils.Success(c, res, "ok")
}

func (user *User) DelUser(c *gin.Context) {
	var req models.DelUserReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := user_service.DelUser(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (user *User) CionChangeList(c *gin.Context) {
	var req models.ChangeListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := task_service.CionChangeList(&req)
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
