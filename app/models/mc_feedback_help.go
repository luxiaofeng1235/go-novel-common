package models

type McFeedbackHelp struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Title   string `gorm:"column:title" json:"title"`     // 帮助标题
	Content string `gorm:"column:content" json:"content"` // 帮助内容
	Addtime int64  `gorm:"column:addtime" json:"addtime"` // 添加时间
	Uptime  int64  `gorm:"column:uptime" json:"uptime"`
}

func (*McFeedbackHelp) TableName() string {
	return "mc_feedback_help"
}

type HelpListReq struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

type HelpDetailReq struct {
	HelpId int64 `form:"help_id" json:"help_id"`
}

type FeedbackHelpListReq struct {
	Title     string `form:"title" json:"title"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type CreateHelpReq struct {
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
}

type UpdateHelpReq struct {
	HelpId int64 `form:"id" json:"id"`
	CreateHelpReq
}

type DelHelpReq struct {
	HelpIds []int64 `json:"ids" form:"ids"`
}
