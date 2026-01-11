package models

type SysAdmin struct {
	Id            int64  `gorm:"column:id" json:"id" structs:"id"`
	Username      string `gorm:"column:username" json:"username" structs:"username"`                      // 用户名 structs 结构体转map使用的字段
	Mobile        string `gorm:"column:mobile" json:"mobile" structs:"mobile"`                            // 中国手机不带国家代码，国际手机号格式为：国家代码-手机号
	Nickname      string `gorm:"column:nickname" json:"nickname" structs:"nickname"`                      // 用户昵称
	Password      string `gorm:"column:password" json:"password" structs:"password"`                      // 登录密码;
	Status        int    `gorm:"column:status" json:"status" structs:"status"`                            // 用户状态;0:禁用,1:正常,2:未验证
	Email         string `gorm:"column:email" json:"email" structs:"email"`                               // 用户登录邮箱
	Sex           int    `gorm:"column:sex" json:"sex" structs:"sex"`                                     // 性别;0:保密,1:男,2:女
	Avatar        string `gorm:"column:avatar" json:"avatar" structs:"avatar"`                            // 用户头像
	LastLoginTime int64  `gorm:"column:last_login_time" json:"last_login_time" structs:"last_login_time"` // 最后登录时间
	LastLoginIp   string `gorm:"column:last_login_ip" json:"last_login_ip" structs:"last_login_ip"`       // 最后登录ip
	CreateTime    int64  `gorm:"column:create_time" json:"create_time" structs:"create_time"`             // 注册时间
	UpdateTime    int64  `gorm:"column:update_time" json:"update_time" structs:"update_time"`             // 注册时间
	Remark        string `gorm:"column:remark" json:"remark" structs:"remark"`                            // 备注
}

func (*SysAdmin) TableName() string {
	return "sys_admin"
}

type Login struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	IdKeyC   string `form:"idKeyC"   json:"idKeyC"`
	IdValueC string `form:"idValueC" json:"idValueC"`
}

// 获取管理员列表结构体
type AdminListReq struct {
	Username  string `form:"username" json:"username"`
	Mobile    string `form:"mobile" json:"mobile"`
	Nickname  string `form:"nickname" json:"nickname"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	Status    string `form:"status" json:"status"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

// 创建用户结构体
type CreateAdminReq struct {
	Nickname string `form:"nickname" json:"nickname"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Mobile   string `form:"mobile" json:"mobile"`
	Email    string `form:"email" json:"email"`
	Avatar   string `form:"avatar" json:"avatar"`
	Sex      int    `form:"sex"    json:"sex"`
	Remark   string `form:"remark" json:"remark"`
	Status   int    `form:"status" json:"status"`
	RoleIds  []int  `form:"roleIds" json:"roleIds"`
}
type UpdateAdminReq struct {
	AdminId int64 `json:"adminId"`
	CreateAdminReq
}

// 批量删除用户结构体
type DeleteAdminReq struct {
	AdminIds []int64 ` form:"adminIds" json:"adminIds"`
}

// 重置用户密码状态参数
type ResetPwdReq struct {
	AdminId  int64  `form:"adminId" json:"adminId"`
	Password string `form:"password" json:"password"`
}

// 设置用户状态参数
type AdminStatusReq struct {
	AdminId int64 `form:"adminId" json:"adminId"`
	Status  int   `form:"status" json:"status"`
}

// 更新密码结构体
type ChangePwdReq struct {
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	NewPassword string `form:"newPassword" json:"newPassword"`
}

// 修改个人信息
type EditProfile struct {
	Avatar   string `form:"avatar"   json:"avatar"`
	Nickname string `form:"nickname" json:"nickname"`
	Mobile   string `form:"mobile"   json:"mobile"`
	Email    string `form:"email"    json:"email"`
	Sex      int    `form:"sex"      json:"sex"`
	Remark   string `form:"remark"   json:"remark"`
}
