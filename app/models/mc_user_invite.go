package models

type McUserInvite struct {
	Id       int64  `gorm:"column:id" json:"id"`
	Uid      int64  `gorm:"column:uid" json:"uid"`           // 用户ID
	Inviteid int64  `gorm:"column:inviteid" json:"inviteid"` // 邀请人ID
	Deviceid string `gorm:"column:deviceid" json:"deviceid"` // 设备ID
	Addtime  int64  `gorm:"column:addtime" json:"addtime"`   // 邀请时间
}

func (*McUserInvite) TableName() string {
	return "mc_user_invite"
}
