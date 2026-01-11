package models

type CommentReplyListReq struct {
	CommentId int64 `form:"cid" json:"cid"`
	BookId    int64 `form:"bid" json:"bid"`
	Page      int   `form:"page" json:"page"`
	Size      int   `form:"size" json:"size"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type CommentReplyRes struct {
	Id          int64                  `form:"id" json:"id"`
	Text        string                 `form:"text" json:"text"`
	City        string                 `form:"city" json:"city"`
	ReplyNum    int64                  `forms:"reply_num" json:"reply_num"`
	IsFollow    int                    `form:"is_follow" json:"is_follow"`
	IsPraise    int                    `form:"is_praise" json:"is_praise"`
	PraiseCount int64                  `form:"praise_count" json:"praise_count"`
	Addtime     int64                  `form:"addtime" json:"addtime"`
	User        *CommentReplyUserRes   `form:"user" json:"user"`
	Book        *CommentReplyBookRes   `form:"book" json:"book"`
	ReplyList   []*CommentReplyListRes `form:"reply_list" json:"reply_list"`
}

type CommentReplyUserRes struct {
	Uid      int64  `form:"uid" json:"uid"`
	Nickname string `form:"nickname" json:"nickname"`
	Pic      string `form:"pic" json:"pic"`
}

type CommentReplyBookRes struct {
	BookId   int64  `form:"bid" json:"bid"`
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author" json:"author"`
	Pic      string `form:"pic" json:"pic"`
}

type CommentReplyListRes struct {
	Id             int64  `forms:"id" json:"id"`
	IsPraise       int    `form:"is_praise" json:"is_praise"`
	PraiseCount    int64  `form:"praise_count" json:"praise_count"`
	Text           string `form:"text" json:"text"`
	City           string `form:"city" json:"city"`
	Addtime        int64  `form:"addtime" json:"addtime"`
	Uid            int64  `form:"uid" json:"uid"`
	Nickname       string `form:"nickname" json:"nickname"`
	Pic            string `form:"pic" json:"pic"`
	ParentId       int64  `forms:"parent_id" json:"parent_id"`
	ParentText     string `form:"parent_text" json:"parent_text"`
	ParentNickname string `form:"parent_nickname" json:"parent_nickname"`
}

type ReplyListReq struct {
	Page   int   `form:"page" json:"page"`
	Size   int   `form:"size" json:"size"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type ReplyDetailReq struct {
	CommentId int64 `form:"cid" json:"cid"`
	UserId    int64 `form:"user_id" json:"user_id"`
}

type CommentReplyList struct {
	Id             int64  `forms:"id" json:"id"`
	PraiseCount    int64  `form:"praise_count" json:"praise_count"`
	Text           string `form:"text" json:"text"`
	Uid            int64  `form:"uid" json:"uid"`
	Nickname       string `form:"nickname" json:"nickname"`
	Pic            string `form:"pic" json:"pic"`
	ParentText     string `form:"parent_text" json:"parent_text"`
	ParentId       int64  `forms:"parent_id" json:"parent_id"`
	ParentNickname string `form:"parent_nickname" json:"parent_nickname"`
	Ip             string `form:"ip" json:"ip"`
	City           string `form:"city" json:"city"`
	Status         int    `form:"status" json:"status"`
	Addtime        int64  `form:"addtime" json:"addtime"`
}
