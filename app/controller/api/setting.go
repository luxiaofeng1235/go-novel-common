package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/setting_service"
	"go-novel/utils"
)

type Setting struct{}

func (setting *Setting) GetValue(c *gin.Context) {
	var req models.GetValueReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	value, err := setting_service.GetValue(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, value, "ok")
}

// 获取APP的配置信息
func (setting *Setting) GetAppConfigInfo(c *gin.Context) {
	//获取配置信息
	itemData, err := setting_service.GetValueByNameInfo("advertNewbieSet")
	if err != nil {
		utils.FailEncrypt(c, err, "获取失败")
		return
	}

	//转换切片数据返回
	mMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(itemData), &mMap)
	if err != nil {
		utils.FailEncrypt(c, err, "转换失败")
		return
	}
	ret := gin.H{
		"app_config": mMap,
	}
	utils.SuccessEncrypt(c, ret, "获取appID数据")
}

func (setting *Setting) GetRanksName(c *gin.Context) {
	ranks, err := setting_service.GetRanksName()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, ranks, "ok")
}
