package models

type McAppVersion struct {
	Id              int64  `gorm:"column:id" json:"id"`
	Device          string `gorm:"column:device" json:"device"`                     // 设备终端android ios
	Version         string `gorm:"column:version" json:"version"`                   // 最新版本号
	DownUrl         string `gorm:"column:down_url" json:"down_url"`                 // app下载地址
	IsForce         int    `gorm:"column:is_force" json:"is_force"`                 // 是否强制更新
	ForceType       int    `gorm:"column:force_type" json:"force_type"`             //更新类型
	ForceUrl        string `gorm:"column:force_url" json:"force_url"`               //本地上传的url信息
	ForceFileSize   int    `gorm:"column:force_file_size" json:"force_file_size"`   //本地文件大小
	UpdateText      string `gorm:"column:update_text" json:"update_text"`           // app更新文案
	CommentStatus   int    `gorm:"column:comment_status" json:"comment_status"`     //书评开关
	CopyrightStatus int    `gorm:"column:copyright_status" json:"copyright_status"` //版全面状态
	PackageId       int    `gorm:"column:package_id" json:"package_id"`             //项目关联的ID
	Addtime         int64  `gorm:"column:addtime" json:"addtime"`                   // 添加时间
	Uptime          int64  `gorm:"column:uptime" json:"uptime"`                     // 更新时间
}

func (*McAppVersion) TableName() string {
	return "mc_app_version"
}

// 后台返回的列表状态信息
type AppVersionAdminListRes struct {
	Id              int64  `form:"id" json:"id"`
	Device          string `form:"device" json:"device"`
	Version         string `form:"version" json:"version"`
	DownUrl         string `form:"down_url" json:"down_url"`
	IsForce         int    `form:"is_force" json:"is_force"`
	ForceType       int    `form:"force_type" json:"force_type"`
	ForceUrl        string `form:"force_url" json:"force_url"`
	UpdateText      string `form:"update_text" json:"update_text"`
	CommentStatus   int    `form:"comment_status" json:"comment_status"`
	CopyrightStatus int    `form:"copyright_status" json:"copyright_status"`
	PackageId       int    `form:"package_id" json:"package_id"`
	Addtime         int64  `form:"addtime" json:"addtime"`
	Uptime          int64  `form:"uptime" json:"uptime"`
	ProjectName     string `gorm:"-"` //关联项目名称
	PackageName     string `gorm:"-"` //关联包名称
	AppId           string `gorm:"-"`
	Mark            string `gorm:"-"` //渠道号
}

type AppVersionListReq struct {
	Device      string `form:"device" json:"device"`
	PageNum     int    `form:"pageNum" json:"pageNum"`
	ProjectName string `form:"projectName" json:"projectName"` //项目名称
	PackageName string `form:"packageName" json:"packageName"` //包名称
	PageSize    int    `form:"pageSize" json:"pageSize"`
}

type UpdateAppVersionReq struct {
	VersionId       int64  `form:"id"  json:"id"`
	Version         string `form:"version"  json:"version"`
	DownUrl         string `form:"down_url" json:"down_url"`
	PackageId       int    `form:"package_id" json:"package_id"`
	ForceFileSize   int    `form:"force_file_size" json:"force_file_size"`
	IsForce         int    `form:"is_force" json:"is_force"`
	ForceType       int    `form:"force_type" json:"force_type"`
	ForceUrl        string `form:"force_url" json:"force_url"`
	UpdateText      string `form:"update_text" json:"update_text"`
	CommentStatus   int    `form:"comment_status" json:"comment_status"`     //书评状态 1：开启 0：关闭
	CopyrightStatus int    `form:"copyright_status" json:"copyright_status"` //版权面状态1：开启 0：关闭
}

// 更新
type CreateAppVersionReq struct {
	Device          string `form:"device" json:"device"`
	Version         string `form:"version"  json:"version"`
	DownUrl         string `form:"down_url" json:"down_url"`
	PackageId       int    `form:"package_id" json:"package_id"`
	IsForce         int    `form:"is_force" json:"is_force"`
	ForceType       int    `form:"force_type" json:"force_type"`
	ForceUrl        string `form:"force_url" json:"force_url"`
	ForceFileSize   int    `form:"force_file_size" json:"force_file_size"`
	UpdateText      string `form:"update_text" json:"update_text"`
	CommentStatus   int    `form:"comment_status" json:"comment_status"`     //书评状态 1：开启 0：关闭
	CopyrightStatus int    `form:"copyright_status" json:"copyright_status"` //版权面状态1：开启 0：关闭
}

type GetVersionInfo struct {
	Device string `form:"device"  json:"device"`
}

type GetVersionInfoNew struct {
	Device      string `form:"device"  json:"device"`             //设备号
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
	UserId      int64  `form:"user_id" json:"user_id"`            //用户ID
}

// 删除版本
type DeleteVersionReq struct {
	VersionIds []int64 `json:"ids" form:"ids"`
}
