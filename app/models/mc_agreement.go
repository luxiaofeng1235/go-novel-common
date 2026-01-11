package models

type McAgreement struct {
	Id                 int64  `gorm:"column:id" json:"id"`                                     // ID
	PackageId          int    `gorm:"column:package_id" json:"package_id"`                     //项目ID，等同于包ID
	Qdh                string `gorm:"column:qdh" json:"qdh"`                                   // 所属渠道 ，一般为baidu、huawei等
	UserAgreementUrl   string `gorm:"column:user_agreement_url" json:"user_agreement_url"`     //用户协议的url
	PrivacyUrl         string `gorm:"column:privacy_url" json:"privacy_url"`                   //隐私协议的url
	UserAgreementValue string `gorm:"column:user_agreement_value" json:"user_agreement_value"` // 用户协议配置值
	PrivacyValue       string `gorm:"column:privacy_value" json:"privacy_value"`               //隐私协议的配置值
}

func (*McAgreement) TableName() string {
	return "mc_agreement"
}

type AgreementUpdateOneReq struct {
	PackageId          int    `form:"package_id" json:"package_id"`
	Qdh                string `form:"qdh" json:"qdh"`
	UserAgreementUrl   string `form:"user_agreement_url" json:"user_agreement_url"`
	PrivacyUrl         string `form:"privacy_url" json:"privacy_url"`
	UserAgreementValue string `form:"user_agreement_value" json:"user_agreement_value"`
	PrivacyValue       string `form:"privacy_value" json:"privacy_value"`
	ProjectName        string `form:"project_name" json:"project_name"`
}

type AgreementDetailReq struct {
	PackageId int `form:"package_id" json:"package_id"` //项目包ID
}
