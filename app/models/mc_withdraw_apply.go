package models

type McWithdrawApply struct {
	Id            int64   `gorm:"column:id" json:"id"`
	Uid           int64   `gorm:"column:uid" json:"uid"`                       // 用户ID
	PayType       int     `gorm:"column:pay_type" json:"pay_type"`             // 账号类型 1-支付宝 2-微信
	AccountId     int64   `gorm:"column:account_id" json:"account_id"`         // 取现账号ID
	WithdrawMoney float64 `gorm:"column:withdraw_money" json:"withdraw_money"` // 提现金额
	Cion          int64   `gorm:"column:cion" json:"cion"`                     // 扣除了多少个金币
	Status        int     `gorm:"column:status" json:"status"`                 // 取现状态 0-待处理 1-已处理 2-异常
	Reason        string  `gorm:"column:reason" json:"reason"`                 // 异常原因
	CardName      string  `gorm:"column:card_name" json:"card_name"`           // 收款卡号姓名
	CardNumber    string  `gorm:"column:card_number" json:"card_number"`       // 收款账号
	CardPic       string  `gorm:"column:card_pic" json:"card_pic"`             // 收款码
	CheckTime     int64   `gorm:"column:check_time" json:"check_time"`         // 审核时间
	Addtime       int64   `gorm:"column:addtime" json:"addtime"`               // 时间
	Uptime        int64   `gorm:"column:uptime" json:"uptime"`                 // 更新时间
}

func (*McWithdrawApply) TableName() string {
	return "mc_withdraw_apply"
}

type WithdrawAccountDetailReq struct {
	UserId int64 `form:"uid" json:"uid"`
}

type WithdrawAccountSaveReq struct {
	PayType    int    `form:"pay_type" json:"pay_type"`
	CardName   string `form:"card_name" json:"card_name"`
	CardNumber string `form:"card_number" json:"card_number"`
	CardPic    string `form:"card_pic" json:"card_pic"`
	UserId     int64  `form:"uid" json:"uid"`
}

type WithdrawAccountDelReq struct {
	AccountId int64 `form:"account_id" json:"account_id"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type WithdrawApplyReq struct {
	LimitId   int64 `form:"limit_id" json:"limit_id"`
	AccountId int64 `form:"account_id" json:"account_id"`
	UserId    int64 `form:"uid" json:"uid"`
}

type WithdrawApplyListReq struct {
	Size   int   `form:"size" json:"size"`
	Page   int   `form:"page" json:"page"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type WithdrawApplyListRes struct {
	Id            int64   `form:"id" json:"id"`
	UserId        int64   `form:"uid" json:"uid"`
	PayType       int     `form:"pay_type" json:"pay_type"`
	TypeName      string  `form:"type_name" json:"type_name"`
	WithdrawMoney float64 `form:"withdraw_money" json:"withdraw_money"`
	Cion          int64   `form:"cion" json:"cion"`
	Status        int     `form:"status" json:"status"`
	StatusName    string  `form:"status_name" json:"status_name"`
	Addtime       int64   `form:"addtime" json:"addtime"`
}

type WithdrawListReq struct {
	UserId    string `form:"user_id"  json:"user_id"`
	Status    string `form:"status"  json:"status"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type WithdrawCheckReq struct {
	CheckId int64  `form:"check_id" json:"check_id"`
	Status  int    `form:"status" json:"status"`
	Reason  string `form:"reason" json:"reason"`
}
