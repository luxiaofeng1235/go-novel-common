package models

type McBookRecommand struct {
	Id            int64  `gorm:"column:id" json:"id"`
	RecommandType string `gorm:"column:recommand_type" json:"recommand_type"`
	BookId        int64  `gorm:"column:book_id" json:"book_id"`
	Addtime       int64  `gorm:"column:addtime" json:"addtime"` // 添加时间
}

func (*McBookRecommand) TableName() string {
	return "mc_book_recommand"
}

// 创建推荐书籍的提交提交数据细腻
type CreateBookRecommandReq struct {
	RecommandType string          `form:"recommand_type" json:"recommand_type"` //推荐类型
	RecommandIds  []MySyncRecList `json:"rec_ids"`                              //解析创建的接送对象信息
}

// 搜索指定的推荐类型信息
type SearchBookRecReq struct {
	RecommandType string `form:"recommand_type" json:"recommand_type"` //推荐类型
}

// 处理解析的JSON对象
type MySyncRecList struct {
	BookId        int64  `json:"book_id"`
	RecommandType string `json:"recommand_type"`
}
