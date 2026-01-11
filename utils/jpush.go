package utils

import (
	"fmt"
	"github.com/zwczou/jpush"
)

func JpushMsg(msg string, rids []string) (msgId string, err error) {
	if len(rids) <= 0 {
		err = fmt.Errorf("%v", "推送设备ID不能为空")
		return
	}
	client := jpush.New(JKey, Jsecret)
	//cidList, err := client.PushCid(1, "push")
	payload := &jpush.Payload{
		Platform: jpush.NewPlatform().All(),
		//Audience: jpush.NewAudience().All().SetRegistrationId("140fe1da9e2573338e7"),
		//Audience: jpush.NewAudience().All().SetTag("abc", "ef").SetTagAnd("filmtest"),
		//Audience: jpush.NewAudience().All(),
		Notification: &jpush.Notification{
			Alert: msg, // 通知内容
		},
		Options: &jpush.Options{
			TimeLive:       60,
			ApnsProduction: false,
		},
	}
	payload.Audience = jpush.NewAudience().All().SetRegistrationId(rids...)
	msgId, err = client.Push(payload)
	// msgId, err = client.PushValidate(payload)
	return
}
