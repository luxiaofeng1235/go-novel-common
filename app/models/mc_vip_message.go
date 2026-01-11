package models

type McVipMessage struct {
	Id       int64  `gorm:"column:id" json:"id"`
	Uid      int64  `gorm:"column:uid" json:"uid"`             // 用户iD
	Nickname string `gorm:"column:nickname" json:"nickname"`   // 用户昵称
	Pic      string `gorm:"column:pic" json:"pic"`             // 用户头像
	CardId   int64  `gorm:"column:card_id" json:"card_id"`     // 会员卡ID
	CardName string `gorm:"column:card_name" json:"card_name"` // 会员卡名称
	PayType  int    `gorm:"column:pay_type" json:"pay_type"`   // 1-线上支付 2-金币兑换
	Addtime  int64  `gorm:"column:addtime" json:"addtime"`     // 添加时间
}

func (*McVipMessage) TableName() string {
	return "mc_vip_message"
}

type VipMessageReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}
