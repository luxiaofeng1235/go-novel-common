package models

type McAdver struct {
	Id                   int64  `gorm:"column:id" json:"id"`
	AdverType            int    `gorm:"column:adver_type" json:"adver_type"`           // 广告类型  1-书架广告 2-开屏广告 3-激励视频 4-小说详情页 5-分类下小说列表广告
	AdverName            string `gorm:"column:adver_name" json:"adver_name"`           // 广告名称
	AdverValue           string `gorm:"column:adver_value" json:"adver_value"`         // 广告id -默认安卓广告
	AdverValueIos        string `gorm:"column:adver_value_ios" json:"adver_value_ios"` //广告id ios
	Pic                  string `gorm:"column:pic" json:"pic"`                         // 广告图片
	AdverLink            string `gorm:"column:adver_link" json:"adver_link"`           // 广告链接
	Weight               int    `gorm:"column:weight" json:"weight"`                   // 权重高显示本地广告 权重低显示三方广告
	ClickCount           int    `gorm:"column:click_count" json:"click_count"`         // 广告点击次数
	IsLocal              int    `gorm:"column:is_local" json:"is_local"`               // 是否使用本地广告
	Status               int    `gorm:"column:status" json:"status"`
	AdverPosition        int    `gorm:"column:adver_position" json:"adver_position"`                   //广告浏览设置 1.新用户触发广告时间（半个小时） 2.当天已经使用广告时间（默认半小时） 3.新用户杀进程次数
	AdverNum             int    `gorm:"column:adver_num" json:"adver_num"`                             //广告的总数，阅读中为记录页码，插半屏和书城为具体的时间（默认为秒）
	AdverTime            int    `gorm:"column:adver_time" json:"adver_time"`                           //激励广告免广告时间
	EveryShowNum         int    `gorm:"column:every_show_num" json:"every_show_num"`                   //每天显示的广告的总次数控制(激励视频的有效，一天多少次最大)
	ErrorNum             int    `gorm:"column:error_num" json:"error_num"`                             //失败次数
	ReadTurningSHowTimes int    `gorm:"column:read_turning_show_times" json:"read_turning_show_times"` //阅读中翻页的次数限制
	ReadRollTime         int    `gorm:"column:read_roll_time" json:"read_roll_time"`                   //阅读中滚动的时间
	Addtime              int64  `gorm:"column:addtime" json:"addtime"`                                 // 添加时间
	Uptime               int64  `gorm:"column:uptime" json:"uptime"`                                   // 更新时间
}

// 只作为类型返回数据字段信息
type AdverProjectTypeListRes struct {
	AdverType int    `form:"adver_type"  json:"adver_type"`
	AdverName string `form:"adver_name" json:"adver_name"`
}

func (*McAdver) TableName() string {
	return "mc_adver"
}

type AdverListReq struct {
	AdverName string `form:"adver_name" json:"adver_name"`
	AdverCode string `form:"adver_code" json:"adver_code"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type CreateAdverReq struct {
	AdverType            int    `form:"adver_code" json:"adver_type"`
	IsLocal              int    `form:"is_local" json:"is_local"`
	Weight               int    `form:"weight" json:"weight"`
	AdverName            string `form:"adver_name"  json:"adver_name"`
	Pic                  string `form:"pic"  json:"pic"`
	AdverLink            string `form:"adver_link"  json:"adver_link"`
	AdverValue           string `form:"adver_value" json:"adver_value"`
	AdverValueIos        string `form:"adver_value_ios" json:"adver_value_ios"`
	Status               int    `form:"status" json:"status"`
	AdverPosition        int    `form:"adver_position" json:"adver_position"`
	AdverNum             int    `form:"adver_num" json:"adver_num"`
	AdverTime            int    `form:"adver_time" json:"adver_time"`
	ErrorNum             int    `form:"error_num" json:"error_num"`
	EveryShowNum         int    `form:"every_show_num" json:"every_show_num"`
	ReadTurningSHowTimes int    `form:"read_turning_show_times" json:"read_turning_show_times"`
	ReadRollTime         int    `form:"read_roll_time" json:"read_roll_time"`
}

type UpdateAdverReq struct {
	AdverId int64 `form:"id"  json:"id"`
	CreateAdverReq
}

type DeleteAdverReq struct {
	AdverId int64 `json:"id" form:"id"`
}

// 获取单挑参数配置信息
type GetAdverReq struct {
	AdverId int64 `form:"id"  json:"id"`
}

type UpdateClickCountReq struct {
	AdverValue string `form:"adver_value" json:"adver_value"`
}
