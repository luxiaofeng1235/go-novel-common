package models

type McBuy struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Text    string `gorm:"column:text" json:"text"`       // 消费简介
	Mid     int64  `gorm:"column:mid" json:"mid"`         // 漫画ID
	Bid     int64  `gorm:"column:bid" json:"bid"`         // 小说ID
	Cid     int64  `gorm:"column:cid" json:"cid"`         // 章节ID
	Uid     int64  `gorm:"column:uid" json:"uid"`         // 消费会员ID
	Cion    int64  `gorm:"column:cion" json:"cion"`       // 消费积分
	Ip      string `gorm:"column:ip" json:"ip"`           // IP
	Addtime int64  `gorm:"column:addtime" json:"addtime"` // 消费时间
}

func (*McBuy) TableName() string {
	return "mc_buy"
}
