package notify_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func SendNotify(notifyType string, parentText, sendPic string, sendUid, receiveId int64, notifyName, notifyText string, targetId int64) (err error) {
	notify := models.McNotify{
		SendPic:    sendPic,
		SendUid:    sendUid,
		ParentText: parentText,
		ReceiveUid: receiveId,
		NotifyType: notifyType,
		NotifyName: notifyName,
		NotifyText: notifyText,
		TargetId:   targetId,
		IsRead:     0,
		Addtime:    utils.GetUnix(),
	}
	if err = global.DB.Create(&notify).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
