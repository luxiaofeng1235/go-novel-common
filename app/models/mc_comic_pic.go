package models

type McComicPic struct {
	Id     int64  `gorm:"id" json:"id"`
	Cid    int64  `gorm:"cid" json:"cid"`       // 章节ID
	Mid    int64  `gorm:"mid" json:"mid"`       // 漫画ID
	Img    string `gorm:"img" json:"img"`       // 图片url地址
	Width  int    `gorm:"width" json:"width"`   // 图片宽度
	Height int    `gorm:"height" json:"height"` // 图片高度
	Xid    int    `gorm:"xid" json:"xid"`       // 排序ID
	Md5    string `gorm:"md5" json:"md5"`       // 源地址MD5
}

func (*McComicPic) TableName() string {
	return "mc_comic_pic"
}
