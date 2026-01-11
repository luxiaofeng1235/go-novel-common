package common_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/pkg/config"
	"go-novel/utils"
	"io/ioutil"
	"net/http"
	"time"
)

// 获取config配置
func GetConfigByName(key string) (configMap map[string]string, err error) {
	config, err := models.GetConfigByKey(key)
	if err != nil {
		return nil, err
	}

	configMap = make(map[string]string)
	for _, value := range config {
		configMap[value.Name] = value.Value
	}

	return configMap, nil
}

// 获取vivo的access_token
func GetVivoAccessToken() (accessToken string, err error) {
	configMap, err := GetConfigByName(config.VIVO)
	if err != nil {
		return "", err
	}
	//判断access_token是否过期
	if utils.StrToInt64(configMap["token_date"])+7200 < time.Now().UnixMilli() {
		accessToken = configMap["access_token"]
		//access_token过期，重新获取
		//accessToken, err = GetVivoAccessTokenFromServer()
	} else {
		accessToken = configMap["access_token"]
	}
	return accessToken, nil
}

// access_token过期重新获取
func GetVivoAccessTokenFromServer() (accessToken string, err error) {
	configMap, _ := GetConfigByName(config.VIVO)
	//config是否有值
	if len(configMap) == 0 {
		fmt.Println("vivo配置为空")
		return
	}

	//拿refresh_token去获取access_token
	url := fmt.Sprintf("http://marketing-api.vivo.com.cn/openapi/v1/oauth2/refreshToken?client_id=%s&client_secret=%s&refresh_token=%s", configMap["client_id"], configMap["client_secret"], configMap["refresh_token"])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var GetAccessTokenResponse models.GetAccessTokenResponse
	err = json.Unmarshal(body, &GetAccessTokenResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Post: %+v\n", GetAccessTokenResponse)
	if GetAccessTokenResponse.Code != 0 {
		global.Errlog.Error("GetVivoAccessTokenFromServer error: ", GetAccessTokenResponse.Message)
		return "", err
	}
	tx := global.DB.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}
	// 更新 access_token
	err = tx.Model(&models.Config{}).Where("name = ? and key = ?", "access_token", config.VIVO).Update("value", GetAccessTokenResponse.Data.AccessToken).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}
	// 更新 token_date
	err = tx.Model(&models.Config{}).Where("name = ? and key = ?", "token_date", config.VIVO).Update("value", GetAccessTokenResponse.Data.TokenDate).Error
	if err != nil {
		tx.Rollback()
		return "", err
	}
	//更新refresh_token
	err = tx.Model(&models.Config{}).Where("name = ? and key = ?", "refresh_token", config.VIVO).Update("value", GetAccessTokenResponse.Data.RefreshToken).Error
	if err != nil {
		tx.Rollback()
	}
	//更新refresh_token_date
	err = tx.Model(&models.Config{}).Where("name = ? and key = ?", "refresh_token_date", config.VIVO).Update("value", GetAccessTokenResponse.Data.RefreshTokenDate).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit().Error
	if err != nil {
		return "", err
	}
	return GetAccessTokenResponse.Data.AccessToken, nil
}
