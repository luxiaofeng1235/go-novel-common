package task_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/notify_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetTaskById(id int64) (task *models.McTask, err error) {
	err = global.DB.Model(models.McTask{}).Where("id", id).First(&task).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTaskListByWelfareType(welfareType int) (tasks []*models.McTask, err error) {
	err = global.DB.Model(models.McTask{}).Order("id asc").Where("status = 1 and welfare_type = ?", welfareType).Find(&tasks).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTaskCountByUid(tid, uid, agoTime int64) (count int64) {
	db := global.DB.Model(models.McTaskList{}).Order("id desc")
	db = db.Where("tid = ? and uid = ?", tid, uid)
	if agoTime > 0 {
		db = db.Where("addtime > ?", agoTime)
	}
	var err error
	err = db.Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func CionChangeList(req *models.CionChangeListReq) (changes []*models.CionChangeListRes, total int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}

	var list []*models.McCionChange
	db := global.DB.Model(&models.McCionChange{}).Order("id desc")

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	if len(list) > 0 {
		for _, val := range list {
			var changeName, operatName string
			if val.ChangeType == 1 {
				changeName = "增加"
			} else if val.ChangeType == 2 {
				changeName = "减少"
			}
			if val.Tid > 0 {
				operatName = getTaskNameById(val.Tid)
			} else {
				if val.OperatType == 1 {
					operatName = "每日签到"
				} else if val.OperatType == 2 {
					operatName = "补签"
				} else if val.OperatType == 4 {
					operatName = "邀请"
				} else if val.OperatType == 5 {
					operatName = "兑换人民币提现"
				} else if val.OperatType == 6 {
					operatName = "兑换会员"
				}
			}
			change := &models.CionChangeListRes{
				Id:         val.Id,
				Tid:        val.Tid,
				UserId:     val.Uid,
				Cion:       val.Cion,
				OperatType: val.OperatType,
				OperatName: operatName,
				ChangeType: val.ChangeType,
				ChangeName: changeName,
				Addtime:    val.Addtime,
			}
			changes = append(changes, change)
		}
	}
	return
}

func getTaskNameById(id int64) (taskName string) {
	var err error
	err = global.DB.Model(models.McTask{}).Select("task_name").Where("id", id).First(&taskName).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getTodaySecondsByUserId(userId, todayUnix int64) (seconds int64) {
	err := global.DB.Model(models.McBookTime{}).Select("coalesce(sum(second), 0)").Where("uid = ? and addtime >= ?", userId, todayUnix).Scan(&seconds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func InviteGiveReward(parentId, userId int64, nickName, userPic string) (err error) {
	if parentId <= 0 {
		return
	}
	inviteGiveCion := setting_service.GetInviteGive()
	if err != nil {
		return
	}
	if inviteGiveCion <= 0 {
		return
	}
	tx := global.DB.Begin()
	change := models.McCionChange{
		Tid:        0,
		Uid:        parentId,
		Cion:       inviteGiveCion,
		ChangeType: 1,
		OperatType: 4,
		Addtime:    utils.GetUnix(),
	}
	err = tx.Create(&change).Error
	if err != nil {
		tx.Rollback()
		return
	}
	data := make(map[string]interface{})
	data["cion"] = gorm.Expr("cion + ?", inviteGiveCion)
	err = tx.Model(models.McUser{}).Where("id = ?", parentId).Updates(data).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	//邀请通知
	_ = notify_service.SendNotify(utils.Invite, "", userPic, userId, parentId, fmt.Sprintf("新用户%v邀请成功通知", nickName), fmt.Sprintf("您已经成功邀请1人,奖励金币%v", inviteGiveCion), userId)
	return
}
