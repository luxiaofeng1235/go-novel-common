package models

type SysAdminOnline struct {
	Id         int64  `gorm:"column:id" json:"id"`
	Username   string `gorm:"column:username" json:"username"`       // 用户名
	Token      string `gorm:"column:token" json:"token"`             // 用户token
	CreateTime int64  `gorm:"column:create_time" json:"create_time"` // 登录时间
	ExpireTime int64  `gorm:"column:expire_time" json:"expire_time"` // 过期时间
	Ip         string `gorm:"column:ip" json:"ip"`                   // 登录ip
	Browser    string `gorm:"column:browser" json:"browser"`         // 浏览器
	Os         string `gorm:"column:os" json:"os"`                   // 操作系统
}

func (*SysAdminOnline) TableName() string {
	return "sys_admin_online"
}

// 获取在线用户列表结构体
type AdminOnlineListReq struct {
	Username string `json:"username" form:"username"`
	Ip       string `json:"ip" form:"ip"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 管理员强退结构体
type ForceLogoutReq struct {
	OnlineIds []int64 `json:"onlineIds" form:"onlineIds"`
}
