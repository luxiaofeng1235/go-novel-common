package models

type AddNodeReq struct {
	UserId int64  `form:"user_id" json:"user_id"`
	Key    string `form:"key"    json:"key"`
}

type SendNoReadCount struct {
	Key string `form:"key"    json:"key"`
}

type ApiAesReq struct {
	Data string `form:"data"    json:"data"`
}

type SendChapters struct {
	BookId   int64  `form:"bid" json:"bid"`
	SourceId int64  `form:"source_id" json:"source_id"`
	Sort     string `form:"sort" json:"sort"` //asc-正序 desc-倒序
	Key      string `form:"key"    json:"key"`
}
