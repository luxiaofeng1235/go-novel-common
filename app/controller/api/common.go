package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/common/common_service"
	"go-novel/app/service/common/upload_service"
	"go-novel/global"
	"go-novel/pkg/config"
	"go-novel/utils"
	"go-novel/utils/e"
	"net/http"
	"strings"
)

type Common struct{}

func (common *Common) UploadImg(c *gin.Context) {
	uploadType := c.PostForm("upload_type")
	if uploadType == "" {
		uploadType = "pic"
	}
	url, err := upload_service.UploadFile(c, uploadType, "")
	if err != nil {
		utils.Fail(c, err, "上传图片失败")
		return
	}
	//替换对应的的长路经
	newPath := strings.ReplaceAll(url, utils.REPLACEFOLDER, "")
	res := gin.H{
		"domain":   utils.GetApiUrl(),
		"url":      url,
		"show_url": newPath, //回显的字段
	}
	utils.Success(c, res, "ok")
	return
}

// vivo点击检测回传
func (common *Common) VivoCallBackClk(c *gin.Context) {
	var req []models.ClickCallback
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "error")
		return
	}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		fmt.Println("转换json失败", err)
	}
	global.VivoClicklog.Info(string(jsonStr))

	configMap, _ := common_service.GetConfigByName(config.VIVO)
	//config是否有值
	if len(configMap) == 0 {
		fmt.Println("vivo配置为空")
		return
	}
	common_service.SaveVivoClick(req)
	res := gin.H{
		"message": "successful",
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": res,
	})
}

// 神马平台统计
func (common *Common) ChaojihuiStatistics(c *gin.Context) {
	params := c.Request.URL.Query()
	var callbackUrl string
	// 打印参数
	for key, value := range params {
		//解析获取到的callback_url内容
		if key == "callback" {
			callbackUrl = value[0]
		}
		fmt.Printf("chaojihui Parameter %s: %v\n", key, value)
	}
	//log.Printf("callback: %s", callbackUrl)
	global.SmClicklog.Infof("记录神马平台的请求的参数为： params=%+v", params)
	res := gin.H{
		"message": "successful",
	}
	//处理神马搜索的回调统计
	common_service.SyncShenmaClickData(c, callbackUrl)
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS200,
		"data": res,
	})

}

// 百度统计代码监听地址/上报广告主回调地址给百度
func (common *Common) BaiduStatistics(c *gin.Context) {

	params := c.Request.URL.Query()
	var callbackUrl string
	// 打印参数
	for key, value := range params {
		//解析获取到的callback_url内容
		if key == "callback_url" {
			callbackUrl = value[0]
		}
		fmt.Printf("Parameter %s: %v\n", key, value)
	}
	global.BaiduClicklog.Info("记录当前百度请求的参数为： params=", params)
	var activateUrl string //这里面直接用百度的激活地址保存为当前的，为后面的上报铺垫
	//处理百度回调地址上报
	if callbackUrl != "" && strings.Contains(callbackUrl, "ocpc.baidu") != false {
		global.BaiduClicklog.Info("获取到的url中的callbak_url参数为 : ", callbackUrl)
		//请求百度的激活上报地址
		activateUrl = utils.GetReplaceBaiduCallbak(callbackUrl, "activate", 0) //获取激活的url
		//对其进行加密
		sign := utils.Md5(activateUrl)
		activateUrl = activateUrl + "&sign=" + sign //获取拼装的sign的url进行数据上报
		global.BaiduClicklog.Info("转换后的激活上报地址 : ", activateUrl)
		//请求上报接口进行上报
		utils.GetBaiduResponse(activateUrl) //请求激活

		//请求百度的注册上报地址
		registerUrl := utils.GetReplaceBaiduCallbak(callbackUrl, "register", 0) //获取注册的url
		sign1 := utils.Md5(registerUrl)
		registerUrl = registerUrl + "&sign=" + sign1 //获取拼装的sign的url进行数据上报
		global.BaiduClicklog.Info("转换后的注册上报地址:", registerUrl)
		////请求上报接口进行上报
		utils.GetBaiduResponse(registerUrl) //请求注册

	} else {
		global.BaiduClicklog.Info("解析callback参数为空或者回调异常 callback=", activateUrl)
	}
	//同步上报的百度的点击记录
	common_service.SyncBaiduClickData(c, callbackUrl)
	res := gin.H{
		"message": "successful",
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS200,
		"data": res,
	})
}

// 发送短信验证码
func (common *Common) SendCode(c *gin.Context) {
	var req models.SendCode
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//发送短信验证
	tel := strings.TrimSpace(req.Tel)
	if tel != "" {
		//发送短信验证码主流程
		//log.Printf("request phone is : %s", tel)
		//发送验证码
		smsResult, err := common_service.TencentSmsSend(tel)
		if err != nil {
			utils.FailEncrypt(c, err, "")
			return
		}
		//log.Printf("smsResult is : %s", smsResult)
		//检查下发状态，成功返回是OK
		sendRes := utils.GetSmsSendCode(smsResult)
		//返回发送的标记状态默认使用里面的状态值
		if sendRes != "Ok" {
			utils.FailEncrypt(c, nil, "发送失败")
			return
		}
		//_, err := common_service.TelSend(tel)
		//if err != nil {
		//	utils.FailEncrypt(c, err, "")
		//	return
		//}
	}

	//发送邮箱验证码
	email := strings.TrimSpace(req.Email)
	if email != "" {
		_, err := common_service.EmailSend(email)
		if err != nil {
			utils.FailEncrypt(c, err, "")
			return
		}
	}

	res := gin.H{
		"tel":   tel,
		"email": email,
		//"code":  smsCode,
	}

	utils.SuccessEncrypt(c, res, "ok")
}
