package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/task_service"
	"go-novel/utils"
	"strconv"
)

type Task struct{}

func (task *Task) TaskList(c *gin.Context) {
	var req models.TaskListSearchReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := task_service.TaskListSearch(&req)
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

func (task *Task) CreateTask(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateTaskReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := task_service.CreateTask(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (task *Task) UpdateTask(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateTaskReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := task_service.UpdateTask(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	taskId, _ := strconv.Atoi(c.Query("id"))
	taskInfo, err := task_service.GetTaskById(int64(taskId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"taskInfo": taskInfo,
	}
	utils.Success(c, res, "ok")
}

func (task *Task) DelTask(c *gin.Context) {
	var req models.DeleteTaskReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := task_service.DeleteTask(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (task *Task) RecordList(c *gin.Context) {
	var req models.TaskRecordReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := task_service.TaskRecord(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}
	taskList := task_service.GetTaskByWelfareType()
	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
		"taskList":    taskList,
	}
	utils.Success(c, res, "ok")
}
