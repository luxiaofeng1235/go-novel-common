package models

// 使用对应的数据结构体进行解析
type McAdverPackage struct {
	Id          int64  `gorm:"column:id" json:"id"`                     //自增主键ID
	ProjectName string `gorm:"column:project_name" json:"project_name"` //项目名称
	AppId       string `gorm:"column:app_id" json:"app_id"`             //appid信息信息
	PackageName string `gorm:"column:package_name" json:"package_name"` //广告包名
	DeviceType  string `gorm:"column:device_type" json:"device_type"`   //渠道类型 ios|android
	CreateUser  string `gorm:"column:create_user" json:"create_user"`   //创建用户名
	Mark        string `gorm:"column:mark" json:"mark"`                 //渠道号
	AddTime     int64  `gorm:"column:addtime" json:"addtime"`           //添加时间
	Uptime      int64  `gorm:"column:uptime" json:"uptime"`             //更新时间
}

func (*McAdverPackage) TableName() string {
	return "mc_adver_package"
}

type AdverSettingRes struct {
	Id               int64  `form:"id" json:"id"`
	ProjectName      string `form:"project_name" json:"project_name"`
	AppId            string `form:"app_id" json:"app_id"`
	PackageName      string `form:"package_name" json:"package_name"`
	DeviceType       string `form:"device_type" json:"device_type"`
	CreateUser       string `form:"create_user" json:"create_user"`
	Mark             string `form:"mark" json:"mark"` //渠道号统计
	Addtime          int64  `form:"addtime" json:"addtime"`
	Uptime           int64  `form:"uptime" json:"uptime"`
	UserAgreementUrl string `gorm:"-"` //关联项目名称
	PrivacyUrl       string `gorm:"-"` //关联包名称
}

type AdvertProjectInfoRes struct {
	Id          int64                  `form:"id" json:"id"`
	ProjectName string                 `form:"project_name" json:"project_name"`
	AppId       string                 `form:"app_id" json:"app_id"`
	PackageName string                 `form:"package_name" json:"package_name"`
	DeviceType  string                 `form:"device_type" json:"device_type"`
	CreateUser  string                 `form:"create_user" json:"create_user"`
	Mark        string                 `form:"mark" json:"mark"` //渠道号统计
	Extra       map[string]interface{} `gorm:"-"`
}

// 搜索包管理需要用到的参数
type AdverPackageAqiReq struct {
	PackageName string `form:"package_name" json:"package_name"`
	DeviceType  string `form:"device_type" json:"device_type"`
	Ip          string `form:"ip" json:"ip"`      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"` //渠道号
}

// 搜索包管理需要用到的参数
type AdverPackageReq struct {
	ProjectName string `form:"project_name"  json:"project_name"`
	AppId       string `form:"app_id"  json:"app_id"`
	PackageName string `form:"package_name" json:"package_name"`
	DeviceType  string `form:"device_type" json:"device_type"`
	PageNum     int    `form:"pageNum" json:"pageNum"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
}

// 创建广告包的管理名称
type CreateAdverPackageReq struct {
	ProjectName string             `form:"project_name"  json:"project_name"`
	AppId       string             `form:"app_id"  json:"app_id"`
	PackageName string             `form:"package_name" json:"package_name"`
	DeviceType  string             `form:"device_type" json:"device_type"`
	CreateUser  string             `form:"create_user" json:"create_user"`
	Mark        string             `form:"mark" json:"mark"`
	PushData    []MyAjaxModelsList `json:"push_data"` //解析创建的接送对象信息
}

// 处理解析的JSON对象
type MyAjaxModelsList struct {
	AdverType        int    `json:"adver_type"`
	AdverTypeName    string `json:"adver_type_name"`
	AdverValueString string `json:"adver_value_string"`
}

// 更新包需要提交的数据
type UpdateAdverPackageReq struct {
	PackageId int64 `form:"id"  json:"id"`
	CreateAdverPackageReq
}

// 删除广告的包ID信息
type DeleteAdverPackageReq struct {
	PackageIds []int64 `json:"ids" form:"ids"`
}

// 获取广告包的传参
type GetAdverPackageReq struct {
	PackageId int64 `form:"id"  json:"id"`
}
