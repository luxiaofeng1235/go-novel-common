package models

type McBookTime struct {
	Id      int64  `gorm:"column:id" json:"id"`
	Uid     int64  `gorm:"column:uid" json:"uid"`         // 用户ID
	Bid     int64  `gorm:"column:bid" json:"bid"`         // 小说ID
	Second  int64  `gorm:"column:second" json:"second"`   // 当天阅读秒数
	Day     string `gorm:"column:day" json:"day"`         // 日期 例: 20131206
	Addtime int64  `gorm:"column:addtime" json:"addtime"` // 创建时间
	Uptime  int64  `gorm:"column:uptime" json:"uptime"`   // 修改时间
}

type GetBookTime struct {
	Bid int64 `form:"bid"  json:"bid"`
}

func (*McBookTime) TableName() string {
	return "mc_book_time"
}
