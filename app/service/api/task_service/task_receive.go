package task_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

// 任务状态 0-未完成 1-已完成未领取 2-已完成已领取
func GetIsReceiveByUid(tid, uid, agoTime int64) (completeNum, alreadyNum, taskStatus, isReceive int) {
	var err error
	var task = new(models.McTaskList)
	db := global.DB.Model(models.McTaskList{})
	db = db.Where("tid = ? and uid = ?", tid, uid)
	if agoTime > 0 {
		db = db.Where("addtime > ?", agoTime)
	}
	err = db.Last(&task).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	completeNum = task.CompleteNum
	alreadyNum = task.AlreadyNum
	isReceive = task.IsReceive
	if alreadyNum >= completeNum {
		taskStatus = 1
		if task.IsReceive == 1 {
			taskStatus = 2
		}
		if task.LoopCount > alreadyNum {
			taskStatus = 1
			isReceive = 0
		}
	}
	return
}
