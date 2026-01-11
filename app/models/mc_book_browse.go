package models

type McBookBrowse struct {
	Id      int64 `gorm:"column:id" json:"id"`
	Uid     int64 `gorm:"column:uid" json:"uid"`         // 用户ID
	Bid     int64 `gorm:"column:bid" json:"bid"`         // 小说ID
	Addtime int64 `gorm:"column:addtime" json:"addtime"` // 添加时间
	Uptime  int64 `gorm:"column:uptime" json:"uptime"`   // 更新时间
}

func (*McBookBrowse) TableName() string {
	return "mc_book_browse"
}
