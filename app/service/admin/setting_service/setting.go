package setting_service

import (
	"errors"
	"go-novel/app/models"
	"go-novel/global"
	"html"
)

var keys string = "userInstructions,fireCoinApiExportGuide,parameterSetDesc,userInstructionsEn,fireCoinApiExportGuideEn,parameterSetDescEn,bianCoinApiExportGuide,bianCoinApiExportGuideEn"

func GetSetByName(name string) (*models.McSetting, error) {
	var setting *models.McSetting
	err := global.DB.Model(models.McSetting{}).Where("name", name).First(&setting).Error
	return setting, err
}

func GetSetByNameValue(name string) (value string) {
	var setting *models.McSetting
	err := global.DB.Model(models.McSetting{}).Where("name", name).First(&setting).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	value = setting.Value
	return
}

func GetList() (list []*models.McSetting, err error) {
	//从数据库获取
	err = global.DB.Model(models.McSetting{}).Order("id desc").Find(&list).Error
	if err != nil {
		return
	}
	return
}

func GetListByKey(keys []string) (list []*models.McSetting, err error) {
	db := global.DB.Model(models.McSetting{}).Order("id desc")
	if len(keys) > 0 {
		for _, key := range keys {
			db = db.Or("name = ?", key)
		}
	}
	//从数据库获取
	err = db.Find(&list).Error
	if err != nil {
		return
	}
	return
}

// 配置列表
func SelectList() (list []*models.McSetting, err error) {
	list, err = GetList()
	for _, val := range list {
		//index := strings.Index(keys,val.Name)
		//if index != -1{
		val.Value = html.UnescapeString(val.Value)
		//}
	}
	if err != nil {
		err = errors.New("获取数据失败")
		return
	}
	return
}

// 获取项目对应的协议信息
func GetAgreementInfo(package_id int64) (info *models.McAgreement, err error) {
	if package_id < 0 {
		return
	}
	err = global.DB.Model(models.McAgreement{}).Where("package_id", package_id).Debug().Find(&info).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//处理回显的数据变更
	info.UserAgreementValue = html.UnescapeString(info.UserAgreementValue)
	info.PrivacyValue = html.UnescapeString(info.PrivacyValue)
	return
}

// 更新站点信息
func UpdateAgreementData(req models.AgreementUpdateOneReq) error {
	var err error
	var count int64
	err = global.DB.Model(&models.McAgreement{}).Where("package_id", req.PackageId).Count(&count).Error
	var htmlUserAgreement = req.UserAgreementValue //原始html内容的协议
	var htmlPrivacy = req.PrivacyValue             //原始的html的隐私协议
	req.UserAgreementValue = html.EscapeString(req.UserAgreementValue)
	req.PrivacyValue = html.EscapeString(req.PrivacyValue)
	if count > 0 {
		//编辑
		var mapData = make(map[string]interface{})
		mapData["package_id"] = req.PackageId //项目包ID
		mapData["qdh"] = req.Qdh              //渠道
		//判断用户协议不为空
		if req.UserAgreementValue != "" {
			//同步FTP数据到服务器上
			mapData["user_agreement_value"] = req.UserAgreementValue                           //注册协议内容
			userAgreementUrl := AsyncFtpUpload(req.ProjectName, req.Qdh, htmlUserAgreement, 1) //同步用户协议到ftp
			if userAgreementUrl != "" {
				mapData["user_agreement_url"] = userAgreementUrl
			}
		}
		//判断隐私协议不为空
		if req.PrivacyValue != "" {
			mapData["privacy_value"] = req.PrivacyValue                            //隐私协议内容
			privacyUrl := AsyncFtpUpload(req.ProjectName, req.Qdh, htmlPrivacy, 2) //同步隐私协议到ftp
			if privacyUrl != "" {
				mapData["privacy_url"] = privacyUrl
			}
		}
		err = global.DB.Model(&models.McAgreement{}).Debug().Where("package_id", req.PackageId).Updates(mapData).Error
	} else {
		//添加
		setting := models.McAgreement{
			PackageId: req.PackageId, //项目包ID
			Qdh:       req.Qdh,       //渠道
		}
		//用户协议判断
		if req.UserAgreementValue != "" {
			setting.UserAgreementValue = req.UserAgreementValue
			userAgreementUrl := AsyncFtpUpload(req.ProjectName, req.Qdh, htmlUserAgreement, 1) //同步用户协议到ftp
			if userAgreementUrl != "" {
				setting.UserAgreementUrl = userAgreementUrl //用户协议关联
			}
		}
		//隐私协议判断
		if req.PrivacyValue != "" {
			setting.PrivacyValue = req.PrivacyValue
			privacyUrl := AsyncFtpUpload(req.ProjectName, req.Qdh, htmlPrivacy, 2) //同步隐私协议到ftp
			if privacyUrl != "" {
				setting.PrivacyUrl = privacyUrl //隐私协议关联
			}
		}
		err = global.DB.Model(&models.McAgreement{}).Debug().Create(&setting).Error
	}
	if err != nil {
		return errors.New("更新项目协议设置信息失败")
	}
	return nil
}

// 更新站点信息
func UpdateSave(req models.SettingUpdateReq) error {
	var err error
	for key, val := range req.WebContent {
		//index := strings.Index(keys,key)
		//if index != -1{
		val = html.EscapeString(val.(string))
		//}
		db := global.DB.Model(&models.McSetting{})
		var count int64
		err = db.Where("name", key).Count(&count).Error
		if count > 0 {
			//更新
			var mapData = make(map[string]interface{})
			mapData["value"] = val
			err = db.Where("name", key).Updates(mapData).Error
		} else {
			//添加
			setting := models.McSetting{
				Name:  key,
				Value: val.(string),
			}
			err = db.Create(&setting).Error
		}
	}

	if err != nil {
		return errors.New("更新配置信息失败")
	}
	return nil
}

func UpdateSaveOne(req models.SettingUpdateOneReq) error {
	var err error
	var count int64
	err = global.DB.Model(&models.McSetting{}).Where("name", req.Key).Count(&count).Error
	req.Value = html.EscapeString(req.Value)
	if count > 0 {
		var mapData = make(map[string]interface{})
		mapData["value"] = req.Value
		err = global.DB.Model(&models.McSetting{}).Where("name", req.Key).Updates(mapData).Error
	} else {
		//添加
		setting := models.McSetting{
			Name:  req.Key,
			Value: req.Value,
		}
		err = global.DB.Model(&models.McSetting{}).Create(&setting).Error
	}

	if err != nil {
		return errors.New("更新配置信息失败")
	}
	return nil
}

func GetValueByName(name string) (value string, err error) {
	err = global.DB.Model(models.McSetting{}).Select("value").Where("name = ?", name).Scan(&value).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	value = html.UnescapeString(value)
	return
}
