package models

type McBookSearch struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Uid        int64  `gorm:"column:uid" json:"uid"`                 // 用户ID
	SearchName string `gorm:"column:search_name" json:"search_name"` // 搜索名称
	Num        int64  `gorm:"column:num" json:"num"`                 // 搜索次数
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`         // 时间
	Uptime     int64  `gorm:"column:uptime" json:"uptime"`           // 更新时间
}

func (*McBookSearch) TableName() string {
	return "mc_book_search"
}

type SearchHistoryReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type SearchHotReq struct {
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}
