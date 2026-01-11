package models

type McIncome struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Text    string `gorm:"column:text" json:"text"`       // 收入简介
	Mid     int64  `gorm:"column:mid" json:"mid"`         // 漫画ID
	Bid     int64  `gorm:"column:bid" json:"bid"`         // 小说ID
	Uid     int64  `gorm:"column:uid" json:"uid"`         // 收入会员ID
	Cion    int    `gorm:"column:cion" json:"cion"`       // 分成金额
	Zcion   int64  `gorm:"column:zcion" json:"zcion"`     // 总金额
	Addtime int64  `gorm:"column:addtime" json:"addtime"` // 收入时间
}

func (*McIncome) TableName() string {
	return "mc_income"
}
