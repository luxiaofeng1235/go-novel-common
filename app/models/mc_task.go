package models

type McTask struct {
	Id          int64  `gorm:"column:id" json:"id"`
	WelfareType int    `gorm:"column:welfare_type" json:"welfare_type"` // 福利类型 1-新人福利 2-日常福利 3-阅读福利
	TaskType    int    `gorm:"column:task_type" json:"task_type"`       // 任务类型 1-单次任务 2-每日任务 3-无限制任务
	TaskName    string `gorm:"column:task_name" json:"task_name"`       // 任务标题
	Desc        string `gorm:"column:desc" json:"desc"`                 // 任务介绍
	Status      int    `gorm:"column:status" json:"status"`             // 是否开启：0-关闭，1-关闭
	Cion        int64  `gorm:"column:cion" json:"cion"`                 // 奖励金币
	Vip         int64  `gorm:"column:vip" json:"vip"`                   // 奖励VIP天数
	Sort        int    `gorm:"column:sort" json:"sort"`                 // 排序
	ReadMinute  int64  `gorm:"column:read_minute" json:"read_minute"`   // 需要阅读分钟数
	DayNum      int64  `gorm:"column:day_num" json:"day_num"`           // 当日任务要求完成次数
	LoopCount   int    `gorm:"column:loop_count" json:"loop_count"`     // 循环任务最大循环次数
	CompleteNum int    `gorm:"column:complete_num" json:"complete_num"` // 每次需完成数量
	Addtime     int64  `gorm:"column:addtime" json:"addtime"`           // 添加时间
	Uptime      int64  `gorm:"column:uptime" json:"uptime"`             // 更新时间
}

func (*McTask) TableName() string {
	return "mc_task"
}

type TaskListReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type TaskListRes struct {
	News            []*TaskInfoRes `form:"news" json:"news"`
	Dailys          []*TaskInfoRes `form:"dailys" json:"dailys"`
	Reads           []*TaskInfoRes `form:"reads" json:"reads"`
	ReadText        string         `form:"read_text" json:"read_text"`
	TodayReadSecond int64          `form:"today_read_second" json:"today_read_second"`
	ReadHighest     int64          `form:"read_highest" json:"read_highest"`
}

type TaskInfoRes struct {
	Id          int64  `form:"id" json:"id"`
	WelfareType int    `form:"welfare_type" json:"welfare_type"`
	TaskType    int    `form:"task_type" json:"task_type"`
	TaskName    string `form:"task_name" json:"task_name"`
	Desc        string `form:"desc" json:"desc"`
	Cion        int64  `form:"cion" json:"cion"`
	Vip         int64  `form:"vip" json:"vip"`
	CompleteNum int    `form:"complete_num" json:"complete_num"`
	AlreadyNum  int    `form:"already_num" json:"already_num"`
	TaskStatus  int    `form:"task_status" json:"task_status"`
	IsReceive   int    `form:"is_receive" json:"is_receive"`
	LoopCount   int    `form:"loop_count" json:"loop_count"`
}

type TaskListSearchReq struct {
	WelfareType string `form:"welfare_type" json:"welfare_type"`
	TaskName    string `form:"task_name" json:"task_name"`
	PageNum     int    `form:"pageNum" json:"pageNum"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
}

type CreateTaskReq struct {
	WelfareType int    `form:"welfare_type"  json:"welfare_type"`
	TaskType    int    `form:"task_type"  json:"task_type"`
	TaskName    string `form:"task_name" json:"task_name"`
	Desc        string `form:"desc" json:"desc"`
	Cion        int64  `form:"cion" json:"cion"`
	Vip         int64  `form:"vip" json:"vip"`
	Sort        int    `form:"sort" json:"sort"`
	LoopCount   int    `form:"loop_count" json:"loop_count"`
	CompleteNum int    `form:"complete_num" json:"complete_num"`
	ReadMinute  int64  `form:"read_minute" json:"read_minute"`
	Status      int    `form:"status" json:"status"`
}

type UpdateTaskReq struct {
	TaskId int64 `form:"id"  json:"id"`
	CreateTaskReq
}

type DeleteTaskReq struct {
	TaskId int64 `json:"id" form:"id"`
}
