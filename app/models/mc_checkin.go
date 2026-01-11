package models

type McCheckin struct {
	Id        int64 `gorm:"column:id" json:"id"`
	Uid       int64 `gorm:"column:uid" json:"uid"`               // 签到用户ID
	Day       int   `gorm:"column:day" json:"day"`               // 签到天数
	Cion      int64 `gorm:"column:cion" json:"cion"`             // 奖励金币
	Vip       int64 `gorm:"column:vip" json:"vip"`               // 奖励VIP天数
	IsReissue int   `gorm:"column:is_reissue" json:"is_reissue"` // 是否补签 0-否 1-是
	Retime    int64 `gorm:"column:retime" json:"retime"`         // 补签时间
	Addtime   int64 `gorm:"column:addtime" json:"addtime"`       // 签到时间
}

func (*McCheckin) TableName() string {
	return "mc_checkin"
}

type McCheckinReward struct {
	Id   int64 `gorm:"column:id" json:"id"`
	Day  int   `gorm:"column:day" json:"day"`   // 签到天数
	Cion int64 `gorm:"column:cion" json:"cion"` // 奖励金币
	Vip  int64 `gorm:"column:vip" json:"vip"`   // 奖励VIP天数
}

func (*McCheckinReward) TableName() string {
	return "mc_checkin_reward"
}

type CheckinListReq struct {
	UserId int64 `form:"user_id" json:"user_id"`
}

type CheckinListRes struct {
	List  []*CheckinDay `form:"list" json:"list"`
	Today int           `form:"today" json:"today"`
}

type CheckinDay struct {
	Day       int   `form:"day" json:"day"`
	IsCheck   int   `form:"is_check" json:"is_check"`
	Cion      int64 `form:"cion" json:"cion"`
	Vip       int64 `form:"vip" json:"vip"`
	IsReissue int   `form:"is_reissue" json:"is_reissue"`
}

type CheckinReq struct {
	IsReissue int   `form:"is_reissue" json:"is_reissue"`
	Day       int   `form:"day" json:"day"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type CheckinHistoryReq struct {
	Year   int   `form:"year" json:"year"`
	Month  int   `form:"month" json:"month"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type CheckinHistoryRes struct {
	Day       int   `form:"day" json:"day"`
	Cion      int64 `form:"cion" json:"cion"`
	Vip       int64 `form:"vip" json:"vip"`
	IsCheckin int64 `form:"is_checkin" json:"is_checkin"`
}

type OpenCheckinRemindReq struct {
	IsCheckinRemind int    `form:"is_checkin_remind" json:"is_checkin_remind"`
	RegistId        string `form:"regist_id" json:"regist_id"`
	UserId          int64  `form:"user_id" json:"user_id"`
}

type CheckinRewardListReq struct {
	Day      string `form:"day" json:"day"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateRewardReq struct {
	Day  int   `form:"day"  json:"day"`
	Cion int64 `form:"cion" json:"cion"`
	Vip  int64 `form:"vip" json:"vip"`
}

type UpdateRewardReq struct {
	RewardId int64 `form:"id" json:"id"`
	CreateRewardReq
}

type DeleteRewardReq struct {
	RewardIds []int64 `json:"ids" form:"ids"`
}

type CheckinListSearchReq struct {
	UserId    string `form:"user_id" json:"user_id"`
	Day       string `form:"day" json:"day"`
	IsReissue string `form:"is_reissue" json:"is_reissue"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}
