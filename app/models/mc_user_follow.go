package models

type McUserFollow struct {
	Id      int64 `gorm:"column:id" json:"id"`
	Uid     int64 `gorm:"column:uid" json:"uid"`         // 用户ID
	ByUid   int64 `gorm:"column:by_uid" json:"by_uid"`   // 被关注者用户ID
	Addtime int64 `gorm:"column:addtime" json:"addtime"` // 邀请时间
}

func (*McUserFollow) TableName() string {
	return "mc_user_follow"
}

type FollowUserReq struct {
	FollowType int   `form:"follow_type" json:"follow_type"`
	ByUserId   int64 `form:"by_uid" json:"by_uid"`
	UserId     int64 `form:"user_id" json:"user_id"`
}

type FollowListReq struct {
	FollowType int   `form:"follow_type" json:"follow_type"`
	UserId     int64 `form:"user_id" json:"user_id"`
}

type FollowListRes struct {
	Id       int64  `form:"id" json:"id"`
	Nickname string `form:"nickname" json:"nickname"`
	Addtime  string `form:"addtime" json:"addtime"`
	Pic      string `form:"pic" json:"pic"`
	IsBoth   int    `form:"is_both" json:"is_both"`
}
