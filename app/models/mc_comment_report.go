package models

type McCommentReport struct {
	Id      int64 `gorm:"column:id" json:"id"`           // 小说评论举报ID
	Cid     int64 `gorm:"column:cid" json:"cid"`         // 评论ID
	Uid     int64 `gorm:"column:uid" json:"uid"`         // 评论用户id
	Addtime int64 `gorm:"column:addtime" json:"addtime"` // 创建时间
}

func (*McCommentReport) TableName() string {
	return "mc_comment_report"
}
