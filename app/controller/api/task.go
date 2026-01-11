package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/task_service"
	"go-novel/utils"
)

type Task struct{}

func (task *Task) List(c *gin.Context) {
	var req models.TaskListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	taskRes, err := task_service.GetTaskList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, taskRes, "ok")
}

func (task *Task) Receive(c *gin.Context) {
	var req models.TaskReceiveReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	err = task_service.TaskReceive(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "领取失败")
		return
	}

	utils.SuccessEncrypt(c, "", "领取成功")
}

func (task *Task) Share(c *gin.Context) {
	var req models.TaskShareReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	data, err := task_service.TaskShare(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, data, "ok")
}

func (task *Task) CionChangeList(c *gin.Context) {
	var req models.CionChangeListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, err := task_service.CionChangeList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list":  list,
		"total": total,
	}

	utils.SuccessEncrypt(c, res, "获取列表成功")
}

func (t *Task) CompleteMsgPush(c *gin.Context) {
	var req models.CompleteMsgPushReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	var taskId int64 = 2
	task, err := task_service.GetTaskById(taskId)
	if err != nil {
		return
	}
	err = task_service.PushMsgCompleteTask(task, req.UserId)
	if err != nil {
		utils.FailEncrypt(c, err, "完成消息推送任务失败")
		return
	}
	utils.SuccessEncrypt(c, "", "完成消息推送任务成功")
}

func (t *Task) CompleteVideoReward(c *gin.Context) {
	var req models.CompleteVideoRewardReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	var taskId int64 = 10
	task, err := task_service.GetTaskById(taskId)
	if err != nil {
		return
	}
	var cion, vip int64
	cion, vip, err = task_service.PushVideoRewardTask(task, req.UserId)
	if err != nil {
		utils.FailEncrypt(c, err, "完成视频奖励任务失败")
		return
	}
	res := gin.H{
		"cion": cion,
		"vip":  vip,
	}
	utils.SuccessEncrypt(c, res, "完成视频奖励任务成功")
}
