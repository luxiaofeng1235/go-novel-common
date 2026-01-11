package feedback_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func List(req *models.FeedBackListReq) (feeds []*models.FeedBackListRes, count int64, pageNum, pageSize int, err error) {
	var list []*models.McFeedback
	list, count, pageNum, pageSize, err = FeedBackListSearch(req)
	if len(list) <= 0 {
		return
	}
	for _, val := range list {
		var pics []string
		if val.Pics != "" {
			pics = utils.GetApiPic(val.Pics)
		}
		feed := &models.FeedBackListRes{
			Id:        val.Id,
			Text:      val.Text,
			Pics:      pics,
			Status:    val.Status,
			Reply:     val.Reply,
			Ip:        val.Ip,
			Addtime:   val.Addtime,
			Replytime: val.Replytime,
		}
		feeds = append(feeds, feed)
	}
	return
}

func FeedBackAdd(req *models.FeedBackAddReq) (err error) {
	text := req.Text
	phone := req.Phone
	email := req.Email
	pics := req.Pics
	ip := req.Ip
	userId := req.UserId

	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	if text == "" {
		err = fmt.Errorf("%v", "反馈内容不能为空")
		return
	}
	//if phone == "" {
	//	err = fmt.Errorf("%v", "联系方式不能为空")
	//	return
	//}
	//if email == "" {
	//	err = fmt.Errorf("%v", "联系方式不能为空")
	//	return
	//}

	agoTime := utils.AgoTime(7200)

	count := GetReadCountByBookId(userId, agoTime)
	if count > 10 {
		err = fmt.Errorf("%v", "反馈次数过多,请先休息一会再来反馈")
		return
	}

	feedback := models.McFeedback{
		Text:    text,
		Phone:   phone,
		Email:   email,
		Pics:    pics,
		Ip:      ip,
		Uid:     userId,
		Status:  0,
		Addtime: utils.GetUnix(),
	}
	if err = global.DB.Create(&feedback).Error; err != nil {
		return
	}
	return
}
