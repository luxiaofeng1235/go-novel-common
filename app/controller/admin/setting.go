package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/adver_service"
	"go-novel/app/service/admin/setting_service"
	"go-novel/app/service/common/upload_service"
	"go-novel/utils"
	"io/ioutil"
	"strconv"
	"strings"
)

type Setting struct{}

// 更新站点信息
func (setting *Setting) UpdateInfo(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.SettingUpdateReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		err := setting_service.UpdateSave(req)
		if err != nil {
			utils.Fail(c, err, "更新失败!")
			return
		}
		utils.Success(c, "", "更新成功!")
		return
	}
	utils.Success(c, "", "ok")
}

func (setting *Setting) UpdateInfoOne(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.SettingUpdateOneReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		fmt.Println("c.Request.ContentLength：%v", c.Request.ContentLength)
		data, _ := ioutil.ReadAll(c.Request.Body)
		fmt.Println("c.Request.GetBody: %v", string(data))
		err := setting_service.UpdateSaveOne(req)
		if err != nil {
			utils.Fail(c, err, "更新失败!")
			return
		}
		utils.Success(c, "", "更新成功!")
		return
	}
	utils.Success(c, "", "ok")
}

// 获取相关的项目包里列表
func (setting *Setting) SettingPackageList(c *gin.Context) {
	var req models.AdverPackageReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	list, total, err := adver_service.AdvertSettingPackageList(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}
	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
	return
}

// 更新隐私协议信息
func (setting *Setting) UpdateAgreementOne(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.AgreementUpdateOneReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		//更新隐私协议内容
		err := setting_service.UpdateAgreementData(req)
		if err != nil {
			utils.Fail(c, err, "更新失败!")
			return
		}
		utils.Success(c, "", "更新成功!")
		return
	}
	utils.Success(c, "", "ok")
}

// 查询站点信息
func (setting *Setting) GetInfo(c *gin.Context) {
	//获取配置列表
	list, err := setting_service.SelectList()
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	res := gin.H{
		"list": list,
	}
	utils.Success(c, res, "ok")
}

// 查看单条隐私协议配置信息
func (Setting *Setting) GetAgreementInfo(c *gin.Context) {
	var req models.AgreementDetailReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定错误")
		return
	}
	packageId, _ := strconv.Atoi(c.Query("package_id"))
	info, err := setting_service.GetAgreementInfo(int64(packageId))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	res := gin.H{
		"info": info,
	}
	utils.Success(c, res, "ok")
	return
}

// 根据对应的value获取相关的配置信息
func (setting *Setting) GetInfoByName(c *gin.Context) {
	getNameValue := c.Query("name")                //从远程获取参数
	getNameValue = strings.TrimSpace(getNameValue) //获取接收的参数值
	//参数检查配置
	if getNameValue == "" {
		utils.Fail(c, nil, "参数无效")
		return
	}
	settingInfo, err := setting_service.GetValueByName(getNameValue)
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}
	//获取基础设置的配置信息
	res := gin.H{
		"setting": settingInfo,
	}
	utils.Success(c, res, "获取成功")

}

// 站点设置上传图片和文件
func (setting *Setting) UploadSettingImg(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "setting", "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	utils.Success(c, gin.H{"url": url}, "ok")
	return
}
