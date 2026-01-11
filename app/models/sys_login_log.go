package models

type SysLoginLog struct {
	Id            int64  `gorm:"column:id" json:"id"`                         // 访问ID
	LoginName     string `gorm:"column:login_name" json:"login_name"`         // 登录账号
	Ipaddr        string `gorm:"column:ipaddr" json:"ipaddr"`                 // 登录IP地址
	LoginLocation string `gorm:"column:login_location" json:"login_location"` // 登录地点
	Browser       string `gorm:"column:browser" json:"browser"`               // 浏览器类型
	Os            string `gorm:"column:os" json:"os"`                         // 操作系统
	Status        int    `gorm:"column:status" json:"status"`                 // 登录状态（1成功 0失败）
	Msg           string `gorm:"column:msg" json:"msg"`                       // 提示消息
	LoginTime     int64  `gorm:"column:login_time" json:"login_time"`         // 访问时间
	Module        string `gorm:"column:module" json:"module"`                 // 登录模块
}

func (*SysLoginLog) TableName() string {
	return "sys_login_log"
}

// 获取登录日志列表结构体
type LoginLogListReq struct {
	LoginName string `json:"login_name" form:"login_name"`
	Status    string `json:"status" form:"status" `
	Ipaddr    string `json:"ipaddr" form:"ipaddr"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除登录日志结构体
type DeleteLoginLogReq struct {
	LoginLogIds []int64 `json:"loginLogIds" form:"loginLogIds"`
}
