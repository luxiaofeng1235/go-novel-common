package models

type McCommentPraise struct {
	Id      int64 `gorm:"column:id" json:"id"`           // 小说评论ID
	Cid     int64 `gorm:"column:cid" json:"cid"`         // 评论ID
	ByUid   int64 `gorm:"column:by_uid" json:"by_uid"`   // 被点赞的用户ID
	Uid     int64 `gorm:"column:uid" json:"uid"`         // 用户ID
	Type    int   `gorm:"column:type" json:"type"`       // 评论类型 1-书评 2-回复
	Addtime int64 `gorm:"column:addtime" json:"addtime"` // 点赞时间
}

func (*McCommentPraise) TableName() string {
	return "mc_comment_praise"
}

type PraiseListReq struct {
	Page   int   `form:"page" json:"page"`
	Size   int   `form:"size" json:"size"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type PraiseListRes struct {
	Id          int64  `form:"id" json:"id"`
	Cid         int64  `form:"cid" json:"cid"`
	Uid         int64  `form:"uid" json:"uid"`
	Pic         string `form:"pic" json:"pic"`
	Nickname    string `form:"nickname" json:"nickname"`
	CommentText string `form:"comment_text" json:"comment_text"`
	PraiseText  string `form:"praise_text" json:"praise_text"`
	IsRead      int64  `forms:"is_read" json:"is_read"`
	Addtime     int64  `form:"addtime" json:"addtime"`
}

type PraiseUserReq struct {
	PraiseType int   `form:"praise_type" json:"praise_type"`
	Commentid  int64 `form:"cid" json:"cid"`
	UserId     int64 `form:"user_id" json:"user_id"`
}
