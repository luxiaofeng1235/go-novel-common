package models

type McComment struct {
	Id          int64  `gorm:"column:id" json:"id"`                     // 小说评论ID
	Bid         int64  `gorm:"column:bid" json:"bid"`                   // 被评价的小说ID
	Pid         int64  `gorm:"column:pid" json:"pid"`                   // 上级评论ID
	ParentLink  string `gorm:"parent_link" json:"parent_link"`          // 上级推荐链条
	Uid         int64  `gorm:"column:uid" json:"uid"`                   // 评论用户ID
	ReplyUid    int64  `gorm:"column:reply_uid" json:"reply_uid"`       // 回复评论@某人 被@用户ID
	Text        string `gorm:"column:text" json:"text"`                 // 评论内容
	PraiseCount int64  `gorm:"column:praise_count" json:"praise_count"` // 点赞数
	Score       int    `gorm:"column:score" json:"score"`               // 评分1-10
	Status      int    `gorm:"column:status" json:"status"`             // 状态，0-待审核 1-正常
	ReplyNum    int64  `gorm:"column:reply_num" json:"reply_num"`       // 回复总数
	Machine     string `gorm:"column:machine" json:"machine"`           // 来自pc、wap、app
	Ip          string `gorm:"column:ip" json:"ip"`                     // IP
	City        string `gorm:"column:city" json:"city"`                 // IP城市
	Addtime     int64  `gorm:"column:addtime" json:"addtime"`           // 评论时间
}

func (*McComment) TableName() string {
	return "mc_comment"
}

type CommentRes struct {
	Comments []*CommentListRes `form:"comments" json:"comments"`
	Total    int64             `form:"total" json:"total"`
	Score    float64           `form:"score" json:"score"`
}

type CommentListRes struct {
	Id          int64  `form:"id" json:"id"`
	UserId      int64  `form:"uid" json:"uid"`
	Nickname    string `form:"nickname" json:"nickname"`
	Pic         string `form:"pic" json:"pic"`
	Text        string `form:"text" json:"text"`
	ReplyNum    int64  `form:"reply_num" json:"reply_num"`
	Ip          string `form:"ip" json:"ip"`
	PraiseCount int64  `form:"praise_count" json:"praise_count"`
	IsPraise    int    `form:"is_praise" json:"is_praise"`
	Score       int    `form:"score" json:"score"`
	Addtime     int64  `form:"addtime" json:"addtime"`
}

type CommentListReq struct {
	BookId int64  `form:"bid" json:"bid"`
	Sort   string `form:"sort" json:"sort"`
	Page   int    `form:"page" json:"page"`
	Size   int    `form:"size" json:"size"`
	UserId int64  `form:"user_id" json:"user_id"`
}

type CommentAddReq struct {
	BookId   int64  `form:"bid" json:"bid"`
	Parentid int64  `form:"pid" json:"pid"`             //上级评论ID
	ReplyUid int64  `form:"reply_uid" json:"reply_uid"` //被@用户ID
	Text     string `form:"text" json:"text"`
	Score    int    `form:"score" json:"score"`
	Machine  string `form:"machine" json:"machine"` //来自pc、wap、app
	Ip       string `form:"ip" json:"ip"`
	UserId   int64  `form:"user_id" json:"user_id"`
}

type CommentAddRes struct {
	Id          int64  `form:"id" json:"id"`
	Nickname    string `form:"nickname" json:"nickname"`
	Pic         string `form:"pic" json:"pic"`
	Pid         int64  `form:"pid" json:"pid"`
	Addtime     int64  `form:"addtime" json:"addtime"`
	Text        string `form:"text" json:"text"`
	Uid         int64  `form:"uid" json:"uid"`
	ReplyUid    int64  `form:"reply_uid" json:"reply_uid"`
	ReplyNum    int64  `form:"reply_num" json:"reply_num"`
	PraiseCount int    `form:"praise_count" json:"praise_count"`
	IsPraise    int    `form:"is_praise" json:"is_praise"`
}

type GetCommentReq struct {
	BookId int64 `form:"bid" json:"bid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type CommentDelReq struct {
	CommentId int64 `form:"cid" json:"cid"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type StarGroupReq struct {
	BookId int64 `form:"bid" json:"bid"`
}

type StarGroupRes struct {
	Star  int64 `form:"star" json:"star"`
	Count int64 `form:"count" json:"count"`
}

type CommentReportReq struct {
	CommentId int64 `form:"cid" json:"cid"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type CommentListSearchReq struct {
	BookId   string `form:"book_id" json:"book_id"`
	UserId   string `form:"user_id" json:"user_id"`
	Text     string `form:"text" json:"text"`
	Status   string `form:"status" json:"status"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type UpdateCommentReq struct {
	CommentId   int64  `form:"id" json:"id"`
	Status      int    `form:"status" json:"status"`
	Text        string `form:"text" json:"text"`
	PraiseCount int64  `form:"praise_count" json:"praise_count"`
}

type DeleteCommentReq struct {
	CommentIds []int64 `json:"ids" form:"ids"`
}
