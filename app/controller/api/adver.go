package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/adver_service"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/common_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
)

type Adver struct{}

// 获取新手广告的绑定状态
func (adver Adver) GetNewbieStatus(c *gin.Context) {
	//log.Printf("Request Param: %#v", req) //打印日志信息
	//获取用户的登录信息
	userIdStr, ok := c.Get("user_id")
	if !ok {
		utils.FailEncrypt(c, nil, "获取登陆用户信息失败")
		return
	}
	loginUserId := userIdStr.(int64)
	uuid := utils.GetRequestHeaderByName(c, "Uuid") //获取客户端的UUID
	oaid := utils.GetRequestHeaderByName(c, "Oaid") //获取客户端的oaid
	imei := utils.GetRequestHeaderByName(c, "Imei") //获取IMEI的标记信息
	if oaid != "" || imei == "" {

	}
	var userId int64
	fmt.Println("获取相关的oaid 和uuid、imei等信息", oaid, uuid, imei)
	//计算新手的剩余时间
	userInfo, err := user_service.GetUserByDeviceAndOaid(uuid, oaid, imei)
	if err != nil {
		//异常情况给一个默认的，防止有问题
		userId = loginUserId
		log.Println("用户error = ", err.Error())
	} else {
		userId = userInfo.Id //这里面用用户的UUID反查的去进行关联查询，防止新账号获取不到
		if userId == 0 {
			log.Printf("用户的登录ID为空，默认取登录的配置信息,userId = %v", userId)
			userId = loginUserId
		} else {
			log.Printf("通过uuid = %v 查询能匹配到新用户身份 userId =%v，登录用户ID，%v", uuid, userId, loginUserId)
		}
	}
	everyTimes := book_service.GetReadToUserBeyond(userId)
	res := gin.H{
		"surplus_times": everyTimes,
	}
	utils.SuccessEncrypt(c, res, "获取成功")
}

// 获取项目信息
func (adver Adver) GetProjectInfo(c *gin.Context) {
	var req models.AdverPackageAqiReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c) //获取客户端的IP
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	packageInfo, err := adver_service.GetAdvertPackageInfo(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "获取数据失败")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["id"] = packageInfo.Id
	mapData["project_name"] = packageInfo.ProjectName
	mapData["app_id"] = packageInfo.AppId
	mapData["package_name"] = packageInfo.PackageName
	mapData["device_type"] = packageInfo.DeviceType

	//获取关联的渠道信息 *********************start
	//for key, value := range c.Request.Header {
	//	fmt.Printf("%s: %s\n", key, value)
	//}
	//处理小米事件上报
	common_service.AsyncXiaomiReportEvent(c, req.DeviceType, req.PackageName, mark, 0)

	details, err := adver_service.GetAllProjectList(packageInfo.Id)
	if err != nil {
		global.Errlog.Infof("数据获取失败 %v ", err)
	}
	mapData["details"] = details
	res := gin.H{
		"package_info": mapData,
	}
	utils.SuccessEncrypt(c, res, "ok")
	return
}

// 获取所有广告列表信息
func (adver *Adver) GetAdverMapList(c *gin.Context) {
	var req models.AdverListReq
	//获取用户的广告列表信息
	list, _, err := adver_service.AdverApiSearch(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	//创建一个interface的对象去进行返回
	var mapData = make(map[int]interface{})
	for _, value := range list {
		mapData[value.AdverType] = value
	}
	//result := gin.H{
	//	"list": mapData,
	//}
	utils.SuccessEncrypt(c, mapData, "获取成功")
}

// 获取广告的基础设置配置信息
func (adver *Adver) GetAdvertInfoById(c *gin.Context) {
	var req models.GetAdverReq

	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}

	adverId := req.AdverId //获取对应的ID信息
	log.Printf("Request Param: %#v", req)
	adverInfo, err := adver_service.GetAdverInfoById(int64(adverId))
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	//判断对应的信息
	if adverInfo.Pic != "" {
		adverInfo.Pic = utils.GetFileUrl(adverInfo.Pic)
	}
	adverPosition := adverInfo.AdverPosition
	killProcessTimes := 0
	if adverPosition > 0 && adverPosition == 3 {
		killProcessTimes = 4
	}
	//判断如果满足对应的杀进程次数的配置条件后，直接返回总的杀进程总次数
	typeInfo := map[string]interface{}{
		"id":             adverInfo.Id,
		"adver_type":     adverInfo.AdverType,
		"adver_name":     adverInfo.AdverName,
		"adver_value":    adverInfo.AdverValue,
		"adver_link":     adverInfo.AdverLink,
		"status":         adverInfo.Status,
		"pic":            adverInfo.Pic,
		"is_local":       adverInfo.IsLocal,
		"weight":         adverInfo.Weight,
		"adver_position": adverInfo.AdverPosition,
		"error_num":      adverInfo.ErrorNum,
		//"every_show_num": adverInfo.EveryShowNum, //激励广告设置每天的次数
		//"adver_time":     adverInfo.AdverTime,    //设置广告的时间（只对激励广告有效）
		//"adver_num":      adverInfo.AdverNum,     //设置广告的翻页或者启动切换次数
		"addtime":    adverInfo.Addtime,
		"uptime":     adverInfo.Uptime,
		"kill_times": killProcessTimes, //杀进程的总数
	}
	utils.SuccessEncrypt(c, typeInfo, "获取成功")
}

// 获取广告列表信息
func (adver *Adver) GetAdverMap(c *gin.Context) {

	mdata, err := adver_service.GetAdverMap()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"data": mdata,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (adver *Adver) UpdateClickCount(c *gin.Context) {
	var req models.UpdateClickCountReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	err := adver_service.UpdateClickCount(req.AdverValue)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}
