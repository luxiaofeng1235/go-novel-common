package models

type McVipCard struct {
	Id       int64   `gorm:"column:id" json:"id"`
	CardName string  `gorm:"column:card_name" json:"card_name"` // VIP套餐名称
	Day      int64   `gorm:"column:day" json:"day"`             // 会员天数
	Price    float64 `gorm:"column:price" json:"price"`         // 套餐原价格
	DisRate  float64 `gorm:"column:dis_rate" json:"dis_rate"`   // 打折折扣
	DisPrice float64 `gorm:"column:dis_price" json:"dis_price"` // 套餐打折后价格
	DisDesc  string  `gorm:"column:dis_desc" json:"dis_desc"`   // 打折折扣描述角标
	Status   int     `gorm:"column:status" json:"status"`       // 是否开启：0-关闭，1-关闭
	Daily    float64 `gorm:"column:daily" json:"daily"`         // 折合每天多少 例 ￥0.9/天
	IsRmb    int     `gorm:"column:is_rmb" json:"is_rmb"`       // 是否可以使用人民币购买
	IsCion   int     `gorm:"column:is_cion" json:"is_cion"`     // 是否可以使用金币购买
	Sort     int     `gorm:"column:sort" json:"sort"`           // 排序
	Addtime  int64   `gorm:"column:addtime" json:"addtime"`     // 入库时间
	Uptime   int64   `gorm:"column:uptime" json:"uptime"`       // 更新时间
}

func (*McVipCard) TableName() string {
	return "mc_vip_card"
}

type VipCardCionRes struct {
	Id       int64   `form:"id" json:"id"`
	CardName string  `form:"card_name" json:"card_name"`
	Price    float64 `form:"price" json:"price"`
	DisRate  float64 `form:"dis_rate" json:"dis_rate"`
	DisPrice float64 `form:"dis_price" json:"dis_price"`
	CionRate int64   `form:"cion_rate" json:"cion_rate"`
	Cion     int64   `form:"cion" json:"cion"`
}

type CardListReq struct {
	CardName string `form:"card_name" json:"card_name"`
	IsRmb    string `form:"is_rmb" json:"is_rmb"`
	IsCion   string `form:"is_cion" json:"is_cion"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type CreateCardReq struct {
	CardName string  `form:"card_name"  json:"card_name"`
	Price    float64 `form:"price" json:"price"`
	DisRate  float64 `form:"dis_rate" json:"dis_rate"`
	DisPrice float64 `form:"dis_price" json:"dis_price"`
	DisDesc  string  `form:"dis_desc" json:"dis_desc"`
	Day      int64   `form:"day" json:"day"`
	Daily    float64 `form:"daily" json:"daily"`
	IsRmb    int     `form:"is_rmb" json:"is_rmb"`
	IsCion   int     `form:"is_cion" json:"is_cion"`
	Status   int     `form:"status" json:"status"`
	Sort     int     `form:"sort" json:"sort"`
}

type UpdateCardReq struct {
	CardId int64 `form:"id"  json:"id"`
	CreateCardReq
}

type DeleteCardReq struct {
	CardId int64 `json:"id" form:"id"`
}
