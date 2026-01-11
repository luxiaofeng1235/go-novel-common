package models

type McCity struct {
	Id          int64  `gorm:"column:id" json:"id"`
	CityName    string `gorm:"column:city_name" json:"city_name"`
	EnglishName string `gorm:"column:english_name" json:"english_name"`
	Addtime     int64  `gorm:"column:addtime" json:"addtime"` // 创建
	Uptime      int64  `gorm:"column:uptime" json:"uptime"`   // 修改时间
}

type GetCityReq struct {
	CityId int64 `form:"id"  json:"id"`
}

func (*McCity) TableName() string {
	return "mc_city"
}
