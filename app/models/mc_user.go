package models

type McUser struct {
	Id              int64   `gorm:"column:id" json:"id"`
	ParentId        int64   `gorm:"parent_id" json:"parent_id,omitempty"`              // 父级ID
	Username        string  `gorm:"column:username" json:"username"`                   // 用户名
	Passwd          string  `gorm:"column:passwd" json:"passwd"`                       // 密码
	Nickname        string  `gorm:"column:nickname" json:"nickname"`                   // 昵称
	Tel             string  `gorm:"column:tel" json:"tel"`                             // 手机
	Pic             string  `gorm:"column:pic" json:"pic"`                             // 头像地址
	Email           string  `gorm:"column:email" json:"email"`                         // 邮箱
	Sex             int     `gorm:"column:sex" json:"sex"`                             // 性别 0-未知 1-男 2-女
	Text            string  `gorm:"column:text" json:"text"`                           // 介绍
	Referrer        string  `gorm:"referrer" json:"referrer,omitempty"`                // 上级邀请码
	Invitation      string  `gorm:"invitation" json:"invitation,omitempty"`            // 邀请码
	ParentLink      string  `gorm:"parent_link" json:"parent_link,omitempty"`          // 上级推荐链条
	Vip             int64   `gorm:"column:vip" json:"vip"`                             // 是否VIP 0-否 1-是
	Rmb             float64 `gorm:"column:rmb" json:"rmb"`                             // 账户人民币金额
	Cion            int64   `gorm:"column:cion" json:"cion"`                           // 金币
	Viptime         int64   `gorm:"column:viptime" json:"viptime"`                     // vip到期时间
	Status          int     `gorm:"column:status" json:"status"`                       // 状态 0-锁定 1-正常
	IsGuest         int     `gorm:"column:is_guest" json:"is_guest"`                   // 游客模式 1-是 0-否
	Deviceid        string  `gorm:"column:deviceid" json:"deviceid"`                   // 游客匿名ID
	Mark            string  `gorm:"column:mark" json:"mark"`                           //渠道号
	Package         string  `gorm:"column:package" json:"package"`                     //包名
	Imei            string  `gorm:"column:imei" json:"imei"`                           //imei设备号
	Oaid            string  `gorm:"column:oaid" json:"oaid"`                           //Oaid & imei  同属字段 都为空使用deviceid
	LastLoginTime   int64   `gorm:"column:last_login_time" json:"last_login_time"`     //上一次的登录时间
	Ip              string  `gorm:"column:ip" json:"ip"`                               //Ip地址
	RegistId        string  `gorm:"column:regist_id" json:"regist_id"`                 // 推送设备ID
	IsCheckinRemind int     `gorm:"column:is_checkin_remind" json:"is_checkin_remind"` // 游客模式 1-是 0-否
	LastRemindTime  int64   `gorm:"column:last_remind_time" json:"last_remind_time"`   // 最后一次签到提醒时间
	ReportStatus    int     `gorm:"column:report_status" json:"report_status"`         //上报状态
	BookType        int     `gorm:"column:book_type" json:"book_type"`                 // 是否打开签到提醒 0-关闭 1-正常
	Addtime         int64   `gorm:"column:addtime" json:"addtime"`                     // 注册时间
	Uptime          int64   `gorm:"column:uptime" json:"uptime"`                       // 更新时间
}

func (*McUser) TableName() string {
	return "mc_user"
}

// 获取用户列表结构体
type UserListReq struct {
	Id        int64  `form:"id" json:"id"`
	Nickname  string `form:"nickname" json:"nickname"`
	Username  string `form:"username" json:"username"`
	Referrer  string `form:"referrer" json:"referrer"`
	Tel       string `form:"tel" json:"tel"`
	Email     string `form:"email" json:"email"`
	Status    string `form:"status" json:"status"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type UpdateUserReq struct {
	UserId   int64   `json:"id"`
	Nickname string  `form:"nickname" json:"nickname"`
	Tel      string  `form:"tel" json:"tel"`
	Email    string  `form:"email" json:"email"`
	Rmb      float64 `form:"column:rmb" json:"rmb"`
	Cion     int64   `form:"column:cion" json:"cion"`
	Status   int     `form:"column:status" json:"status"`
	Sex      int     `form:"column:sex" json:"sex"`
}

type DelUserReq struct {
	UserIds []int64 `json:"ids" form:"ids"`
}

type SendCode struct {
	Tel   string `form:"tel" json:"tel"`
	Email string `form:"email" json:"email"`
}

type LoginReq struct {
	LoginType int    `form:"login_type" json:"login_type"`
	Tel       string `form:"tel" json:"tel"`
	Email     string `form:"email" json:"email"`
	Passwd    string `form:"passwd" json:"passwd"`
	Code      string `form:"code" json:"code"`
	Deviceid  string `form:"deviceid" json:"deviceid"`
	Referrer  string `form:"referrer" json:"referrer"`
}

type LogoffReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type GuestLoginReq struct {
	Deviceid string `form:"deviceid" json:"deviceid"`
	Referrer string `form:"referrer" json:"referrer"`
	Sex      int64  `form:"sex" json:"sex"`
}

type ForgotLoginPasswdReq struct {
	Tel  string `form:"tel" json:"tel"`
	Pass string `form:"pass" json:"pass"`
	Code string `form:"code" json:"code"`
}

type UserInfoReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type EditUserReq struct {
	Type      string `form:"type" json:"type"`
	Tel       string `form:"tel" json:"tel"`
	Code      string `form:"code" json:"code"`
	Email     string `form:"email" json:"email"`
	OldPasswd string `form:"old_passwd" json:"old_passwd"`
	Passwd    string `form:"passwd" json:"passwd"`
	Nickname  string `form:"nickname" json:"nickname"`
	Sex       int    `form:"sex" json:"sex"`
	Pic       string `form:"pic" json:"pic"`
	UserId    int64  `form:"user_id" json:"user_id"`
	BookType  int    `form:"book_type" json:"book_type"`
}

type BindRegistIdReq struct {
	RegistId string `form:"regist_id" json:"regist_id"`
	UserId   int64  `form:"user_id" json:"user_id"`
}

type GetInfoRes struct {
	*McUser
	ToDayCion   float64 `form:"today_cion" json:"today_cion"`
	FollowCount int64   `form:"follow_count" json:"follow_count"`
	FansCount   int64   `form:"fans_count" json:"fans_count"`
	CionRmb     float64 `form:"cion_rmb" json:"cion_rmb"`
	Rate        int64   `form:"rate" json:"rate"`
}

type MyInvitRewardsReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}
