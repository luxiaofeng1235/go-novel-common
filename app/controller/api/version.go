package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/api/version_service"
	"go-novel/app/service/common/common_service"
	"go-novel/global"
	"go-novel/pkg/config"
	"go-novel/utils"
	"log"
	"strings"
	"time"
)

type Version struct{}

// 获取新的客户端信息
func (version *Version) GetVersionNewInfo(c *gin.Context) {
	var req models.GetVersionInfoNew
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}

	mark := utils.GetRequestHeaderByName(c, "Mark") //获取设备渠道号
	//调用vivo上报
	if strings.Contains(mark, "vivo") {
		// 查询用户信息
		userInfo, err := user_service.GetUserById(req.UserId)
		if err != nil {
			return
		}
		t1 := time.Unix(userInfo.Addtime, 0)
		t2 := time.Now()

		// 计算时间差
		daysDiff := t2.YearDay() - t1.YearDay()
		if daysDiff == 1 {
			oaid := utils.GetRequestHeaderByName(c, "oaid")
			imei := utils.GetRequestHeaderByName(c, "imei")
			packageName := utils.GetRequestHeaderByName(c, "Package") // 获取包名
			fmt.Println("开始上报vivo 次留")
			global.VivoClicklog.Info("开始上报vivo 次留 package_name:", packageName, "  mark:", mark, "  oaid:", oaid, "  imei:", imei)
			common_service.VivoReportedEvent(c, packageName, mark, oaid, imei, config.VIVO_RETENTION_1)
		}
	}
	//神马平台的搜索
	if strings.Contains(mark, "sm") {
		userInfo, err := user_service.GetUserById(req.UserId)
		if err != nil {
			return
		}
		t1 := time.Unix(userInfo.Addtime, 0)
		t2 := time.Now()
		t1 = time.Unix(1726167600, 0)
		t2 = time.Unix(1726246800, 0)
		daysDiff := t2.YearDay() - t1.YearDay()
		if daysDiff == 1 {
			oaid := utils.GetRequestHeaderByName(c, "oaid")
			imei := utils.GetRequestHeaderByName(c, "imei")
			packageName := utils.GetRequestHeaderByName(c, "Package") // 获取包名
			global.SmClicklog.Info("开始上报shenma 次留")
			global.SmClicklog.Info("开始上报sm 次留 package_name:", packageName, "  mark:", mark, "  oaid:", oaid, "  imei:", imei)
			cvType := config.SM_RETENTION_1 //激活类型
			channel := config.SM_CHALLEN    //渠道标识
			common_service.SmReportEvent(c, packageName, mark, oaid, imei, channel, cvType)
		} else {
			global.SmClicklog.Error("神马次留的上报时间未到")
		}
	}

	//获取渠道号
	req.Mark = mark
	//获取绑定的任务状态信息
	//获取版本更新信息 ,根据端号+包名+渠道来进行匹配
	versionInfo, err := version_service.GetVersionByQdh(req.Device, req.PackageName, req.Mark)
	if err != nil {
		utils.FailEncrypt(c, err, "获取数据失败")
		return
	}
	clientIp := utils.RemoteIp(c)
	global.Errlog.Infof("当前客户端IP ip=%v ", clientIp)
	log.Printf("version info:%+v\n", versionInfo)
	res := gin.H{
		"versionInfo": versionInfo,
	}
	log.Printf("当前的用户ID = %d", req.UserId)

	//接收同步保存游客的基础信息，防止没有请求guest
	var guestReq models.GuestLoginReq
	uuid := utils.GetRequestHeaderByName(c, "Uuid") //获取客户端的UUID
	guestReq.Deviceid = uuid
	guestReq.Sex = 1 //默认男
	_, _, _, _ = user_service.GuestLogin(c, &guestReq)
	utils.SuccessEncrypt(c, res, "ok")
	return
}

func (version *Version) GetVersionInfo(c *gin.Context) {
	var req models.GetVersionInfo
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	var err error
	versionInfo, err := version_service.GetVersionByDevice(req.Device)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"versionInfo": versionInfo,
	}
	utils.SuccessEncrypt(c, res, "获取成功")
}
