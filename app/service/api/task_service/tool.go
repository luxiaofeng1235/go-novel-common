package task_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

// 判断任务完成状态
//func getNums(task *models.McTask, userId int64) (nums, init int64, err error) {
//	if task.TaskType == 1 {
//		todayUnix := utils.GetTodayUnix()
//		nums  = GetTaskCountByUid(task.Id, 0, userId, todayUnix)
//	} else {
//		nums  = GetTaskCountByUid(task.Id, 0,userId, 0)
//	}
//	if task.DayNum == 1 {
//		nums = 0
//	} else {
//		if task.DayNum == 0 || nums < task.DayNum {
//			init = 0
//		} else {
//			init = 1
//		}
//	}
//	return
//}

func getUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getInviteCountByUid(inviteId, agoUnix int64) (count int64) {
	db := global.DB.Model(models.McUserInvite{}).Where("inviteid = ?", inviteId)
	if agoUnix > 0 {
		db = db.Where("inviteid = ? and addtime > ?", inviteId, agoUnix)
	}
	err := db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
