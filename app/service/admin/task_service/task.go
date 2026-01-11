package task_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetTaskList() (taskList []*models.McTask) {
	err := global.DB.Model(models.McTask{}).Where("status = 1").Find(&taskList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTaskById(id int64) (task *models.McTask, err error) {
	err = global.DB.Model(models.McTask{}).Where("id", id).First(&task).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func CheckTaskNameUnique(taskName string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McTask{}).Where("task_name = ?", taskName)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
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

func CionChangeList(req *models.ChangeListReq) (changes []*models.CionChangeListRes, total int64, err error) {
	var list []*models.McCionChange
	db := global.DB.Model(&models.McCionChange{}).Order("id desc")

	// 当pageNum > 0 且 pageSize > 0 才分页

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid =  ?", userId)
	}

	changeType := strings.TrimSpace(req.ChangeType)
	if changeType != "" {
		db = db.Where("change_type = ?", changeType)
	}

	operatType := strings.TrimSpace(req.OperatType)
	if operatType != "" {
		db = db.Where("operat_type = ?", operatType)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}

	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
	}

	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
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

func TaskListSearch(req *models.TaskListSearchReq) (list []*models.McTask, total int64, err error) {
	db := global.DB.Model(&models.McTask{}).Order("id desc")

	welfareType := strings.TrimSpace(req.WelfareType)
	if welfareType != "" {
		db = db.Where("welfare_type = ?", welfareType)
	}

	taskName := strings.TrimSpace(req.TaskName)
	if taskName != "" {
		db = db.Where("task_name = ?", taskName)
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

func CreateTask(req *models.CreateTaskReq) (InsertId int64, err error) {
	welfareType := req.WelfareType
	if welfareType <= 0 {
		err = fmt.Errorf("%v", "福利类型不能为空")
		return
	}

	taskType := req.TaskType
	if taskType <= 0 {
		err = fmt.Errorf("%v", "任务类型不能为空")
		return
	}

	taskName := strings.TrimSpace(req.TaskName)
	if taskName == "" {
		err = fmt.Errorf("%v", "任务名称不能为空")
		return
	}

	if !CheckTaskNameUnique(taskName, 0) {
		err = fmt.Errorf("%v", "任务名称已经存在")
		return
	}

	desc := strings.TrimSpace(req.Desc)
	if desc == "" {
		err = fmt.Errorf("%v", "任务简介不能为空")
		return
	}
	cion := req.Cion
	vip := req.Vip
	if cion <= 0 && vip <= 0 {
		err = fmt.Errorf("%v", "奖励不能为空")
		return
	}
	loopCount := req.LoopCount
	if loopCount <= 0 {
		loopCount = 1
	}
	completeNum := req.CompleteNum
	if completeNum <= 0 {
		completeNum = 1
	}
	sort := req.Sort
	status := req.Status
	readMinute := req.ReadMinute
	task := models.McTask{
		WelfareType: welfareType,
		TaskType:    taskType,
		TaskName:    taskName,
		Desc:        desc,
		Cion:        cion,
		Vip:         vip,
		LoopCount:   loopCount,
		CompleteNum: completeNum,
		Sort:        sort,
		Status:      status,
		ReadMinute:  readMinute,
		Addtime:     utils.GetUnix(),
	}

	if err = global.DB.Create(&task).Error; err != nil {
		return 0, err
	}

	return task.Id, nil
}

func UpdateTask(req *models.UpdateTaskReq) (res bool, err error) {
	id := req.TaskId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	welfareType := req.WelfareType
	taskType := req.TaskType
	taskName := strings.TrimSpace(req.TaskName)
	if taskName == "" {
		err = fmt.Errorf("%v", "任务名称不能为空")
		return
	}
	if !CheckTaskNameUnique(taskName, id) {
		err = fmt.Errorf("%v", "任务名称已经存在")
		return
	}
	desc := strings.TrimSpace(req.Desc)

	cion := req.Cion
	vip := req.Vip
	if cion <= 0 && vip <= 0 {
		err = fmt.Errorf("%v", "奖励不能为空")
		return
	}
	loopCount := req.LoopCount
	if loopCount <= 0 {
		loopCount = 1
	}
	completeNum := req.CompleteNum
	if completeNum <= 0 {
		completeNum = 1
	}
	readMinute := req.ReadMinute
	status := req.Status

	var mapData = make(map[string]interface{})

	if welfareType > 0 {
		mapData["welfare_type"] = welfareType
	}
	if taskType > 0 {
		mapData["task_type"] = taskType
	}
	if taskName != "" {
		mapData["task_name"] = taskName
	}
	if desc != "" {
		mapData["desc"] = desc
	}
	mapData["status"] = status
	mapData["cion"] = cion
	mapData["vip"] = vip
	mapData["sort"] = req.Sort
	mapData["loop_count"] = loopCount
	mapData["complete_num"] = completeNum
	mapData["read_minute"] = readMinute
	if err = global.DB.Model(models.McTask{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteTask(req *models.DeleteTaskReq) (res bool, err error) {
	id := req.TaskId
	if id <= 0 {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return

	}
	err = global.DB.Where("id = ?", id).Delete(&models.McTask{}).Error
	if err != nil {
		return
	}
	return true, nil
}

func TaskRecord(req *models.TaskRecordReq) (list []*models.McTaskList, total int64, err error) {
	db := global.DB.Model(&models.McTaskList{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	taskId := strings.TrimSpace(req.TaskId)
	if taskId != "" {
		db = db.Where("tid = ?", taskId)
	}

	taskName := strings.TrimSpace(req.TaskName)
	if taskName != "" {
		db = db.Where("task_name = ?", taskName)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}

	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
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
