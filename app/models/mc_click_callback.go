package models

import (
	"go-novel/global"
)

type ClickCallback struct {
	Id              int    `json:"id" gorm:"column:id"`
	AdId            int    `json:"adId" gorm:"column:ad_id"`
	RequestId       string `json:"requestId" gorm:"column:request_id"`
	Imei            string `json:"imei" gorm:"column:imei"`
	ClickTime       int64  `json:"clickTime" gorm:"column:click_time"`
	Callback        string `json:"callback" gorm:"column:callback"` //回调地址
	Ua              string `json:"ua" gorm:"column:ua"`
	Oaid            string `json:"oaid" gorm:"column:oaid"`
	Ip              string `json:"ip" gorm:"column:ip"`
	CreativeId      int    `json:"creativeId" gorm:"column:creative_id"`
	MediaType       int    `json:"mediaType" gorm:"column:media_type"`
	AdvertiserId    string `json:"advertiserId" gorm:"column:advertiser_id"`
	AdvertiserName  string `json:"advertiserName" gorm:"column:advertiser_name"`
	AdvertisementId int    `json:"advertisementId" gorm:"column:advertisement_id"`
	PlaceType       int    `json:"placeType" gorm:"column:place_type"`
	AdName          string `json:"adName" gorm:"column:ad_name"`
	GroupId         int    `json:"groupId" gorm:"column:group_id"`
	GroupName       string `json:"groupName" gorm:"column:group_name"`
	CampaignId      int    `json:"campaignId" gorm:"column:campaign_id"`
	Channel         string `json:"channel" gorm:"column:channel"`
	CampaignName    string `json:"campaignName" gorm:"column:campaign_name"`
	CreatedAt       int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       int64  `json:"updated_at" gorm:"column:updated_at"`
}
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 根据requestId查询数据 返回count
func (m *ClickCallback) GetCountByRequestId(requestId string) (int, error) {
	var count int64
	err := global.DB.Model(&ClickCallback{}).Where("request_id = ?", requestId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// 根据imei和oaid 查询数据
func (m *ClickCallback) GetCountByImeiAndOaid(str string) (ClickCallback, error) {
	var clickCallback ClickCallback
	err := global.DB.Model(&ClickCallback{}).Where("imei = ? OR oaid = ?", str, str).First(&clickCallback).Debug().Error
	if err != nil {
		return ClickCallback{}, err
	}
	return clickCallback, nil
}

// 根据IMEI和oaid查询对应的数据信息
func (m *ClickCallback) GetCountByImeiAndOaidType(str string, channel string) (ClickCallback, error) {
	var clickCallback ClickCallback
	err := global.DB.Model(&ClickCallback{}).Where("(imei = ? OR oaid = ?) and channel = ? ", str, str, channel).First(&clickCallback).Debug().Error
	if err != nil {
		return ClickCallback{}, err
	}
	return clickCallback, nil
}

// 批量插入
func (m *ClickCallback) AddCallback() error {
	if m.RequestId == "" {
		return nil
	}
	err := global.DB.Create(&m).Error
	if err != nil {
		return err
	}
	return nil

}

// 添加神马平台的关联数据
func (m *ClickCallback) AddCallbackShenma() error {
	err := global.DB.Create(&m).Error
	if err != nil {
		return err
	}
	return nil

}

// 根据id修改status
func (m *ClickCallback) UpdateStatusById(id int, status int) error {
	err := global.DB.Model(&ClickCallback{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
