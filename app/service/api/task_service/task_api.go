package task_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetTaskList(req *models.TaskListReq) (taskList *models.TaskListRes, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	todayUnix := utils.GetTodayUnix()

	taskList = new(models.TaskListRes)
	var newTasks []*models.McTask
	newTasks, err = GetTaskListByWelfareType(1)

	var dailyTasks []*models.McTask
	dailyTasks, err = GetTaskListByWelfareType(2)

	var readTasks []*models.McTask
	readTasks, err = GetTaskListByWelfareType(3)

	var news []*models.TaskInfoRes
	for _, task := range newTasks {
		completeNum, alreadyNum, taskStatus, isReceive := GetIsReceiveByUid(task.Id, userId, 0)
		taskRes := &models.TaskInfoRes{
			Id:          task.Id,
			WelfareType: task.WelfareType,
			TaskType:    task.TaskType,
			TaskName:    task.TaskName,
			Desc:        task.Desc,
			Cion:        task.Cion,
			Vip:         task.Vip,
			CompleteNum: completeNum,
			AlreadyNum:  alreadyNum,
			LoopCount:   task.LoopCount,
			TaskStatus:  taskStatus,
			IsReceive:   isReceive,
		}
		news = append(news, taskRes)
	}

	var dailys []*models.TaskInfoRes
	for _, task := range dailyTasks {
		completeNum, alreadyNum, taskStatus, isReceive := GetIsReceiveByUid(task.Id, userId, todayUnix)
		taskRes := &models.TaskInfoRes{
			Id:          task.Id,
			WelfareType: task.WelfareType,
			TaskType:    task.TaskType,
			TaskName:    task.TaskName,
			Desc:        task.Desc,
			Cion:        task.Cion,
			Vip:         task.Vip,
			CompleteNum: completeNum,
			AlreadyNum:  alreadyNum,
			LoopCount:   task.LoopCount,
			TaskStatus:  taskStatus,
			IsReceive:   isReceive,
		}
		dailys = append(dailys, taskRes)
	}

	var reads []*models.TaskInfoRes
	var readHighest int64
	for _, task := range readTasks {
		completeNum, alreadyNum, taskStatus, isReceive := GetIsReceiveByUid(task.Id, userId, todayUnix)
		taskRes := &models.TaskInfoRes{
			Id:          task.Id,
			WelfareType: task.WelfareType,
			TaskType:    task.TaskType,
			TaskName:    task.TaskName,
			Desc:        task.Desc,
			Cion:        task.Cion,
			Vip:         task.Vip,
			CompleteNum: completeNum,
			AlreadyNum:  alreadyNum,
			LoopCount:   task.LoopCount,
			TaskStatus:  taskStatus,
			IsReceive:   isReceive,
		}
		reads = append(reads, taskRes)
		readHighest += taskRes.Cion
	}

	taskList.News = news
	taskList.Dailys = dailys
	taskList.Reads = reads
	taskList.ReadText = "阅读得金币"
	taskList.ReadHighest = readHighest
	taskList.TodayReadSecond = getTodaySecondsByUserId(userId, todayUnix)
	return
}

func TaskReceive(req *models.TaskReceiveReq) (err error) {
	taskId := req.TaskId
	userId := req.UserId
	if taskId <= 0 {
		err = fmt.Errorf("%v", "任务ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	user, err := getUserById(userId)
	if err != nil {
		return
	}
	var isReward bool
	isReward, err = AppTaskReward(taskId, user)
	if !isReward {
		if err != nil {
			global.Errlog.Errorf("发放奖励失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
		}
	}
	return
}

func TaskShare(req *models.TaskShareReq) (data *models.TaskShareRes, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	var user *models.McUser
	user, err = getUserById(userId)
	if user.Id <= 0 {
		err = fmt.Errorf("%v", "用户不存在")
		return
	}
	data = new(models.TaskShareRes)
	data.Invitation = user.Invitation
	data.UserId = userId
	data.Pic = utils.GetFileUrl(user.Pic)
	data.Nickname = user.Nickname
	return
}
