package models

type McRank struct {
	Id       int64  `gorm:"column:id" json:"id"`
	RankName string `gorm:"column:rank_name" json:"rank_name"` // 排行榜名称
	RankCode string `gorm:"column:rank_code" json:"rank_code"` // 排行榜标识
	Sort     int    `gorm:"column:sort" json:"sort"`           // 排序ID
	Status   int    `gorm:"column:status" json:"status"`       // 状态,0-停用,1-正常
	Addtime  int64  `gorm:"column:addtime" json:"addtime"`     // 添加时间
	Uptime   int64  `gorm:"column:uptime" json:"uptime"`       // 更新时间
}

func (*McRank) TableName() string {
	return "mc_rank"
}

type RankListReq struct {
	RankName string `form:"rank_name" json:"rank_name"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateRankReq struct {
	RankName string `form:"rank_name"  json:"rank_name"`
	RankCode string `form:"rank_name"  json:"rank_code"`
	Sort     int    `form:"sort" json:"sort"`
	Status   int    `form:"status" json:"status"`
}

type UpdateRankReq struct {
	RankId int64 `form:"id"  json:"id"`
	CreateRankReq
}

type DeleteRankReq struct {
	RankId int64 `json:"id" form:"id"`
}
