package models

type McTaskList struct {
	Id          int64  `gorm:"column:id" json:"id"`
	Uid         int64  `gorm:"column:uid" json:"uid"`                   // 用户ID
	WelfareType int    `gorm:"column:welfare_type" json:"welfare_type"` // 福利类型 1-新人福利 2-日常福利 3-阅读福利
	Tid         int64  `gorm:"column:tid" json:"tid"`                   // 任务ID
	TaskName    string `gorm:"column:task_name" json:"task_name"`       // 任务名称
	TaskType    int    `gorm:"column:task_type" json:"task_type"`       // 任务类型 1-单次任务 2-每日任务 3-循环任务
	Cion        int64  `gorm:"column:cion" json:"cion"`                 // 获得金币
	Vip         int64  `gorm:"column:vip" json:"vip"`                   // 奖励VIP天数
	IsReceive   int    `gorm:"column:is_receive" json:"is_receive"`     // 是否领取 0-未领取 1-已领取
	CompleteNum int    `gorm:"column:complete_num" json:"complete_num"` // 当次需完成数量
	AlreadyNum  int    `gorm:"column:already_num" json:"already_num"`   // 当次已完成数量
	LoopCount   int    `gorm:"column:loop_count" json:"loop_count"`     // 循环任务最大循环次数
	ReadMinute  int64  `gorm:"column:read_minute" json:"read_minute"`   // 需要阅读分钟数
	Addtime     int64  `gorm:"column:addtime" json:"addtime"`           // 任务完成时间
	Uptime      int64  `gorm:"column:uptime" json:"uptime"`             // 更新时间
}

func (*McTaskList) TableName() string {
	return "mc_task_list"
}

type TaskReceiveReq struct {
	TaskId int64 `form:"tid" json:"tid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type TaskShareReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type CompleteMsgPushReq struct {
	TaskId int64 `form:"tid" json:"tid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type CompleteVideoRewardReq struct {
	TaskId int64 `form:"tid" json:"tid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type TaskShareRes struct {
	Invitation string `form:"invitation" json:"invitation"`
	UserId     int64  `form:"uid" json:"uid"`
	Nickname   string `form:"nickname" json:"nickname"`
	Pic        string `form:"pic" json:"pic"`
}

type TaskRecordReq struct {
	UserId    string `form:"user_id" json:"user_id"`
	TaskId    string `form:"task_id" json:"task_id"`
	TaskName  string `form:"task_name" json:"task_name"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}
