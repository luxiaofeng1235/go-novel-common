package models

type McCallBackData struct {
	Id           int64  `gorm:"column:id" json:"id"`
	RequestId    string `gorm:"column:request_id" json:"request_id"`
	CallBackData string `gorm:"column:callback_data" json:"callback_data"`
	IsRep        bool   `gorm:"column:is_rep" json:"is_rep"`
	Source       bool   `gorm:"column:source" json:"source"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}

// 定义用于请求的结构体
type ExtParam struct {
	PayAmount string `json:"payAmount"`
	UnionId   string `json:"unionId"`
}

type DataList struct {
	CreativeId string   `json:"creativeId"`
	CvParam    string   `json:"cvParam"`
	CvTime     int64    `json:"cvTime"`
	CvType     string   `json:"cvType"`
	UserId     string   `json:"userId"`
	DlrSrc     string   `json:"dlrSrc"`
	UserIdType string   `json:"userIdType"`
	RequestId  string   `json:"requestId"`
	CvCustom   string   `json:"CvCustom"`
	ExtParam   ExtParam `json:"extParam"`
}

type RequestData struct {
	DataList []DataList `json:"dataList"`
	PageUrl  string     `json:"pageUrl"`
	PkgName  string     `json:"pkgName"`
	SrcId    string     `json:"srcId"`
	SrcType  string     `json:"srcType"`
	DataFrom string     `json:"dataFrom"`
}

type GetAccessTokenResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken      string `json:"access_token"`
		RefreshToken     string `json:"refresh_token"`
		TokenDate        string `json:"token_date"`
		RefreshTokenDate string `json:"refresh_token_date"`
	}
}

func (*McCallBackData) TableName() string {
	return "mc_callback_data"
}
