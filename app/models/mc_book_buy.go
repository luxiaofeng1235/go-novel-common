package models

type McBookBuy struct {
	Id   int64 `gorm:"column:id" json:"id"`
	Bid  int64 `gorm:"column:bid" json:"bid"`   // 小说ID
	Cid  int64 `gorm:"column:cid" json:"cid"`   // 章节ID
	Uid  int64 `gorm:"column:uid" json:"uid"`   // 用户ID
	Auto int   `gorm:"column:auto" json:"auto"` // 1开启自动购买
}

func (*McBookBuy) TableName() string {
	return "mc_book_buy"
}
