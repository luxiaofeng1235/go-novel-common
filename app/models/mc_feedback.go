package models

type McFeedback struct {
	Id        int64  `gorm:"column:id" json:"id"`
	Uid       int64  `gorm:"column:uid" json:"uid"`             // 反馈人UID
	Text      string `gorm:"column:text" json:"text"`           // 反馈内容
	Ip        string `gorm:"column:ip" json:"ip"`               // 反馈人IP
	Status    int    `gorm:"column:status" json:"status"`       // 反馈处理状态 0-未处理 1-处理中 2-已处理
	Phone     string `gorm:"column:phone" json:"phone"`         // 反馈人手机号
	Email     string `gorm:"column:email" json:"email"`         // 反馈人邮箱
	Pics      string `gorm:"column:pics" json:"pics"`           // 反馈人邮箱
	Reply     string `gorm:"column:reply" json:"reply"`         // 处理说明
	Replytime int64  `gorm:"column:replytime" json:"replytime"` // 处理时间
	Addtime   int64  `gorm:"column:addtime" json:"addtime"`     // 添加时间
}

func (*McFeedback) TableName() string {
	return "mc_feedback"
}

type FeedBackListReq struct {
	Page   int   `form:"page" json:"page"`
	Size   int   `form:"size" json:"size"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type FeedBackListRes struct {
	Id        int64    `form:"id" json:"id"`
	Text      string   `form:"text" json:"text"`
	Pics      []string `form:"pics" json:"pics"`
	Status    int      `form:"status" json:"status"`
	Reply     string   `form:"reply" json:"reply"`
	Ip        string   `form:"ip" json:"ip"`
	Addtime   int64    `form:"addtime" json:"addtime"`
	Replytime int64    `form:"replytime" json:"replytime"`
}

type FeedBackAddReq struct {
	Text   string `form:"text" json:"text"`
	Pics   string `form:"pics" json:"pics"`
	Phone  string `form:"phone" json:"phone"`
	Email  string `form:"email" json:"email"`
	Ip     string `form:"ip" json:"ip"`
	UserId int64  `form:"user_id" json:"user_id"`
}

type FeedBackListSearchReq struct {
	UserId    string `form:"user_id"  json:"user_id"`
	Status    string `form:"status"  json:"status"`
	Text      string `form:"text"  json:"text"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type FeedBackListSearchRes struct {
	Id        int64    `form:"id" json:"id"`
	Text      string   `form:"text" json:"text"`
	Pics      []string `form:"pics" json:"pics"`
	Status    int      `form:"status" json:"status"`
	Reply     string   `form:"reply" json:"reply"`
	Ip        string   `form:"ip" json:"ip"`
	UserId    int64    `form:"uid" json:"uid"`
	Phone     string   `form:"phone" json:"phone"`
	Email     string   `form:"email" json:"email"`
	Addtime   int64    `form:"addtime" json:"addtime"`
	Replytime int64    `form:"replytime" json:"replytime"`
}

type FeedBackReplyReq struct {
	FeedbackId int64  `form:"feedback_id" json:"feedback_id"`
	Status     int    `form:"status" json:"status"`
	Reply      string `form:"reply" json:"reply"`
}
