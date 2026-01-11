package task_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func UpdateTaskStatus(taskId, phaseId int64) (taskName string) {
	var err error
	err = global.DB.Model(models.McTask{}).Where("tid = ? and = phase_id = ?", taskId, phaseId).Update("status", 1).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 完成任务
func CompleteTask(task *models.McTask, userId int64) (err error) {
	taskId := task.Id
	taskType := task.TaskType
	var taskCount int64

	recored := models.McTaskList{
		Uid:         userId,
		Tid:         task.Id,
		WelfareType: task.WelfareType,
		TaskType:    task.TaskType,
		TaskName:    task.TaskName,
		Cion:        task.Cion,
		Vip:         task.Vip,
		IsReceive:   0,
		CompleteNum: task.CompleteNum,
		AlreadyNum:  1,
		LoopCount:   task.LoopCount,
		ReadMinute:  task.ReadMinute,
		Addtime:     utils.GetUnix(),
	}
	if taskType == 1 {
		global.DB.Model(models.McTaskList{}).Where("tid = ? and uid = ? ", taskId, userId).Count(&taskCount)
		if taskCount > 0 {
			err = fmt.Errorf("%v", "该任务只能完成一次")
			return
		} else {
			err = global.DB.Model(models.McTaskList{}).Create(&recored).Error
			if err != nil {
				global.Errlog.Errorf("任务完成失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
				err = nil
				return
			}
		}
	} else if taskType == 2 {
		global.DB.Model(models.McTaskList{}).Where("tid = ? and uid = ? and addtime >= ?", taskId, userId, utils.GetTodayUnix()).Count(&taskCount)
		if taskCount > 0 {
			alreadyNum := gorm.Expr("already_num + 1")
			global.DB.Model(models.McTaskList{}).Order("id desc").Where("tid = ? and uid = ? and addtime >= ?", taskId, userId, utils.GetTodayUnix()).Update("already_num", alreadyNum)
		} else {
			err = global.DB.Model(models.McTaskList{}).Create(&recored).Error
			if err != nil {
				global.Errlog.Errorf("任务完成失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
				err = nil
				return
			}
		}
	} else if taskType == 3 {
		err = global.DB.Model(models.McTaskList{}).Create(&recored).Error
		if err != nil {
			global.Errlog.Errorf("任务完成失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
			err = nil
			return
		}
		var user *models.McUser
		user, err = getUserById(userId)
		if err != nil {
			global.Errlog.Errorf("循环任务自动领取奖励 用户异常%v", err.Error())
			return
		}
		err = GiveUserReward(task, user)
		if err != nil {
			global.Errlog.Errorf("发送奖励失败 %v", err.Error())
			return
		}
	}
	return
}

func PushMsgCompleteTask(task *models.McTask, userId int64) (err error) {
	taskId := task.Id
	var count int64
	global.DB.Model(models.McTaskList{}).Where("tid = ? and uid = ? ", taskId, userId).Count(&count)
	if count > 0 {
		err = fmt.Errorf("%v", "推送任务只能完成一次")
		return
	}
	recored := models.McTaskList{
		Uid:         userId,
		Tid:         task.Id,
		TaskName:    task.TaskName,
		WelfareType: task.WelfareType,
		TaskType:    task.TaskType,
		Cion:        task.Cion,
		Vip:         task.Vip,
		IsReceive:   0,
		CompleteNum: task.CompleteNum,
		AlreadyNum:  1,
		Addtime:     utils.GetUnix(),
	}
	// 创建任务奖励
	err = global.DB.Model(models.McTaskList{}).Create(&recored).Error
	if err != nil {
		global.Errlog.Errorf("打开消息推送完成失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
		err = nil
		return
	}
	return
}

func PushVideoRewardTask(task *models.McTask, userId int64) (cion, vip int64, err error) {
	taskId := task.Id

	todayUnix := utils.GetTodayUnix()
	completeNum, alreadyNum, _, _ := GetIsReceiveByUid(taskId, userId, todayUnix)
	if task.TaskType == 3 {
		if alreadyNum > 0 && alreadyNum >= task.LoopCount && alreadyNum >= completeNum {
			err = fmt.Errorf("%v", "今个领取了太多奖励,明天再来吧")
			return
		}
	} else {
		if alreadyNum > 0 && alreadyNum >= completeNum {
			err = fmt.Errorf("%v", "今个领取了太多奖励,明天再来吧")
			return
		}
	}

	if alreadyNum > 0 {
		global.DB.Model(models.McTaskList{}).Order("id desc").Where("tid = ? and uid = ? and addtime >= ?", taskId, userId, todayUnix).Update("already_num", gorm.Expr("already_num + 1"))
	} else {
		var taskList models.McTaskList
		taskList.Uid = userId
		taskList.Tid = task.Id
		taskList.WelfareType = task.WelfareType
		taskList.TaskType = task.TaskType
		taskList.TaskName = task.TaskName
		taskList.Cion = task.Cion
		taskList.Vip = task.Vip
		taskList.IsReceive = 1
		taskList.CompleteNum = task.CompleteNum
		taskList.AlreadyNum = 1
		taskList.LoopCount = task.LoopCount
		taskList.ReadMinute = task.ReadMinute
		taskList.Addtime = utils.GetUnix()
		// 创建任务奖励
		err = global.DB.Model(models.McTaskList{}).Create(&taskList).Error
		if err != nil {
			global.Errlog.Errorf("视频激励任务完成失败 userId=%v taskId=%v  err=%v", userId, taskId, err.Error())
			err = nil
			return
		}
	}

	cion = task.Cion
	vip = task.Vip

	var user *models.McUser
	user, err = getUserById(userId)
	if err != nil {
		return
	}

	tx := global.DB.Begin()
	us := make(map[string]interface{})
	if cion > 0 {
		us["cion"] = gorm.Expr("cion +?", cion)

		var change models.McCionChange
		change.Uid = userId
		change.Tid = taskId
		change.Cion = cion
		change.ChangeType = 1
		change.OperatType = 3
		change.Addtime = utils.GetUnix()
		err = tx.Model(models.McCionChange{}).Create(&change).Error
		if err != nil {
			tx.Rollback()
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	}
	if vip > 0 {
		us["vip"] = 1
		timeStamp := utils.GetUnix()
		if user.Viptime > timeStamp {
			us["viptime"] = user.Viptime + 86400*vip
		} else {
			us["viptime"] = utils.GetUnix() + 86400*vip
		}
	}
	err = tx.Model(models.McUser{}).Where("id = ?", user.Id).Updates(us).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// 任务奖励
func AppTaskReward(taskId int64, user *models.McUser) (isReward bool, err error) {
	userId := user.Id
	task, err := GetTaskById(taskId)
	if err != nil {
		return
	}
	todayUnix := utils.GetTodayUnix()
	if task.TaskType == 1 {
		todayUnix = 0
	}
	_, _, taskStatus, isReceive := GetIsReceiveByUid(taskId, userId, todayUnix)
	if taskStatus <= 0 {
		err = fmt.Errorf("%v", "该任务未完成")
		return
	}
	if isReceive > 0 {
		err = fmt.Errorf("%v", "奖励已领取")
		return
	}
	err = GiveUserReward(task, user)
	if err != nil {
		global.Errlog.Errorf("发送奖励失败 %v", err.Error())
		return
	}
	isReward = true
	return
}

func GiveUserReward(task *models.McTask, user *models.McUser) (err error) {
	taskId := task.Id
	cion := task.Cion
	vip := task.Vip
	userId := user.Id
	todayUnix := utils.GetTodayUnix()
	if task.TaskType == 1 {
		todayUnix = 0
	}
	tx := global.DB.Begin()
	tx.Model(models.McTaskList{}).Order("id desc").Where("uid = ? and tid = ? and addtime >= ?", userId, taskId, todayUnix).Update("is_receive", 1)
	us := make(map[string]interface{})
	if cion > 0 {
		us["cion"] = gorm.Expr("cion +?", cion)

		var change models.McCionChange
		change.Uid = userId
		change.Tid = taskId
		change.Cion = cion
		change.ChangeType = 1
		change.OperatType = 3
		change.Addtime = utils.GetUnix()
		err = tx.Model(models.McCionChange{}).Create(&change).Error
		if err != nil {
			tx.Rollback()
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	}
	if vip > 0 {
		us["vip"] = 1
		timeStamp := utils.GetUnix()
		if user.Viptime > timeStamp {
			us["viptime"] = user.Viptime + 86400*vip
		} else {
			us["viptime"] = utils.GetUnix() + 86400*vip
		}
	}
	err = tx.Model(models.McUser{}).Where("id = ?", user.Id).Updates(us).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
