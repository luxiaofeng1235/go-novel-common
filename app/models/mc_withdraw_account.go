package models

type McWithdrawAccount struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Uid        int64  `gorm:"column:uid" json:"uid"`                 // 用户ID
	PayType    int    `gorm:"column:pay_type" json:"pay_type"`       // 账号类型 1-支付宝 2-微信
	CardName   string `gorm:"column:card_name" json:"card_name"`     // 收款卡号姓名
	CardNumber string `gorm:"column:card_number" json:"card_number"` // 收款账号
	CardPic    string `gorm:"column:card_pic" json:"card_pic"`       // 收款二维码
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`         // 创建时间
	Uptime     int64  `gorm:"column:uptime" json:"uptime"`           // 更新时间
}

func (*McWithdrawAccount) TableName() string {
	return "mc_withdraw_account"
}

type AccountDetailRes struct {
	*McWithdrawAccount
	CardPicUrl string `form:"card_pic_url" json:"card_pic_url"`
}

type WithdrawAccountListReq struct {
	UserId     string `form:"user_id"  json:"user_id"`
	CardName   string `form:"card_name"  json:"card_name"`
	CardNumber string `form:"card_number"  json:"card_number"`
	PageNum    int    `form:"pageNum" json:"pageNum"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
}
