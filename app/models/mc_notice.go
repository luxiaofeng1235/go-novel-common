package models

type McNotice struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Title   string `gorm:"column:title" json:"title"`     // 公告标题
	Content string `gorm:"column:content" json:"content"` // 公告内容
	Status  int    `gorm:"column:status" json:"status"`   // 公告状态,0-不展示,1-展示
	Link    string `gorm:"column:link" json:"link"`       // 公告链接
	Addtime int64  `gorm:"column:addtime" json:"addtime"` // 创建时间
	Uptime  int64  `gorm:"column:uptime" json:"uptime"`   // 更新时间
}

func (*McNotice) TableName() string {
	return "mc_notice"
}

type NoticeListReq struct {
	Title     string `form:"title" json:"title"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type CreateNoticeReq struct {
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
	Link    string `form:"link" json:"link"`
	Status  int    `form:"status" json:"status"`
}

type UpdateNoticeReq struct {
	NoticeId int64 `form:"id" json:"id"`
	CreateNoticeReq
}

type DelNoticeReq struct {
	NoticeIds []int64 `json:"ids" form:"ids"`
}

type MessageListReq struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

type UpdateIsReadReq struct {
	NotifyId    int64  `form:"notify_id" json:"notify_id"`
	MessageType string `form:"message_type" json:"message_type"`
	UserId      int64  `form:"user_id" json:"user_id"`
}
