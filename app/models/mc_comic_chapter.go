package models

type McComicChapter struct {
	Id      int64  `gorm:"id" json:"id"`
	Mid     int64  `gorm:"mid" json:"mid"`         // 漫画ID
	Xid     int    `gorm:"xid" json:"xid"`         // 排序ID
	Image   string `gorm:"image" json:"image"`     // 图片
	Name    string `gorm:"name" json:"name"`       // 标题
	Jxurl   string `gorm:"jxurl" json:"jxurl"`     // 解析地址
	Vip     int    `gorm:"vip" json:"vip"`         // VIP阅读，0否1是
	Cion    int    `gorm:"cion" json:"cion"`       // 章节需要金币
	Pnum    int    `gorm:"pnum" json:"pnum"`       // 图片数量
	Yid     int    `gorm:"yid" json:"yid"`         // 0已审核，1待审核，2未通过
	Msg     string `gorm:"msg" json:"msg"`         // 未通过原因
	Addtime int64  `gorm:"addtime" json:"addtime"` // 入库时间
}

func (*McComicChapter) TableName() string {
	return "mc_comic_chapter"
}
