package models

type McCionChange struct {
	Id         int64 `gorm:"column:id" json:"id"`
	Tid        int64 `gorm:"column:tid" json:"tid"`                 // 任务ID
	Uid        int64 `gorm:"column:uid" json:"uid"`                 // 用户ID
	Cion       int64 `gorm:"column:cion" json:"cion"`               // 变动金额
	OperatType int   `gorm:"column:operat_type" json:"operat_type"` // 类型 1-每日签到 2-补签 3-任务 4-兑换人民币提现 5-兑换会员
	ChangeType int   `gorm:"column:change_type" json:"change_type"` // 资金类型 1-增加 2-减少
	Addtime    int64 `gorm:"column:addtime" json:"addtime"`         // 签到时间
}

func (*McCionChange) TableName() string {
	return "mc_cion_change"
}

type ChangeListReq struct {
	TaskId     string `form:"tid"  json:"tid"`
	UserId     string `form:"user_id"  json:"user_id"`
	ChangeType string `form:"change_type"  json:"change_type"`
	OperatType string `form:"operat_type"  json:"operat_type"`
	BeginTime  string `form:"beginTime" json:"beginTime"`
	EndTime    string `form:"endTime" json:"endTime"`
	PageNum    int    `form:"pageNum" json:"pageNum"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
}

type CionChangeListReq struct {
	Size   int   `form:"size" json:"size"`
	Page   int   `form:"page" json:"page"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type CionChangeListRes struct {
	Id         int64  `form:"id" json:"id"`
	Tid        int64  `form:"tid" json:"tid"`
	UserId     int64  `form:"uid" json:"uid"`
	Cion       int64  `form:"cion" json:"cion"`
	ChangeType int    `form:"change_type" json:"change_type"`
	ChangeName string `form:"change_name" json:"change_name"`
	OperatType int    `form:"operat_type" json:"operat_type"`
	OperatName string `form:"operat_name" json:"operat_name"`
	Addtime    int64  `form:"addtime" json:"addtime"`
}
