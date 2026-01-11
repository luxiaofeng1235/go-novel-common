package models

type McNotify struct {
	Id         int64  `gorm:"column:id" json:"id"`                   // 通知ID
	ReceiveUid int64  `gorm:"column:receive_uid" json:"receive_uid"` // 接收者用户ID
	ParentText string `gorm:"column:parent_text" json:"parent_text"` // 原来内容
	NotifyType string `gorm:"column:notify_type" json:"notify_type"` // 通知类型 notice:公告 follow:关注 unfollow:取消关注 praise:点赞 unpraise:取消点赞 comment:评论
	IsRead     int    `gorm:"column:is_read" json:"is_read"`         // 是否已读 0-未读 1-已读
	NotifyName string `gorm:"column:notify_name" json:"notify_name"` // 通知标题
	NotifyText string `gorm:"column:notify_text" json:"notify_text"` // 通知内容
	SendUid    int64  `gorm:"column:send_uid" json:"send_uid"`       //发送者用户ID
	SendPic    string `gorm:"column:send_pic" json:"send_pic"`       //发送者用户头像
	ReadTime   int64  `gorm:"column:read_time" json:"read_time"`     // 阅读时间
	TargetId   int64  `gorm:"column:target_id" json:"target_id"`     // 对应类型关联的id
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`         // 创建时间
}

func (*McNotify) TableName() string {
	return "mc_notify"
}
