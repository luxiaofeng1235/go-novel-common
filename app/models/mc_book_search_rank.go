package models

type McBookSearchRank struct {
	Id          int64 `gorm:"column:id" json:"id"` //主键自增ID
	Bid         int64 `form:"bid" json:"bid"`      //小说ID
	SearchCount int64 `form:"search_count" json:"search_count"`
	CreatedAt   int64 `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt   int64 `gorm:"column:updated_at" json:"updated_at"` // 更新时间
}

func (McBookSearchRank) TableName() string {
	return "mc_book_search_rank"
}
