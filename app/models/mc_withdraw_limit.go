package models

type McWithdrawLimit struct {
	Id      int64   `gorm:"column:id" json:"id"`
	Money   float64 `gorm:"column:money" json:"money"`     // 可提现金额
	Read    int64   `gorm:"column:read" json:"read"`       // 阅读时间要求 默认0分钟 0 可以直接提现
	Sort    int     `gorm:"column:sort" json:"sort"`       // 排序
	Status  int     `gorm:"column:status" json:"status"`   // 是否开启：0-关闭，1-关闭
	Addtime int64   `gorm:"column:addtime" json:"addtime"` // 添加时间
	Uptime  int64   `gorm:"column:uptime" json:"uptime"`   // 更新时间
}

func (*McWithdrawLimit) TableName() string {
	return "mc_withdraw_limit"
}

type WithdrawLimitRes struct {
	*McWithdrawLimit
	Cion int64 `form:"cion" json:"cion"`
	Rate int64 `form:"rate" json:"rate"`
}

type WithdrawLimitListReq struct {
	Status   string `form:"status" json:"status"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateLimitReq struct {
	Money  float64 `form:"money"  json:"money"`
	Read   int64   `form:"read" json:"read"`
	Status int     `form:"status" json:"status"`
	Sort   int     `form:"sort" json:"sort"`
}

type UpdateLimitReq struct {
	LimitId int64 `form:"id"  json:"id"`
	CreateLimitReq
}

type DeleteLimitReq struct {
	LimitIds []int64 `json:"ids" form:"ids"`
}
