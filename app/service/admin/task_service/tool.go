package task_service

import (
	"fmt"
	"go-novel/app/models"
)

func GetTaskByWelfareType() (taskList []*models.McTask) {
	taskList = GetTaskList()
	if len(taskList) > 0 {
		for _, val := range taskList {
			text := val.TaskName
			if val.WelfareType == 1 {
				text = fmt.Sprintf("新人福利-%v", text)
			} else if val.WelfareType == 2 {
				text = fmt.Sprintf("日常福利-%v", text)
			} else if val.WelfareType == 3 {
				text = fmt.Sprintf("阅读福利-%v", text)
			}
			val.TaskName = text
		}
	}
	return
}
