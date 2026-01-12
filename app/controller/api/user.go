/*
 * @Descripttion: API 用户控制器（登录/注册/用户信息等）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:30:00
 */
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/common/common_service"
	"go-novel/global"
	"go-novel/pkg/config"
	"go-novel/utils"
	"strings"
)

type User struct{}

func (user *User) Guest(c *gin.Context) {
	var req models.GuestLoginReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}

	//获取客户端的IP
	clientIp := utils.RemoteIp(c)
	global.Requestlog.Infof("游客登录的客户端IP ip=%v ", clientIp)
	//游客登录判断
	userInfo, token, expireTime, err := user_service.GuestLogin(c, &req)
	if err != nil {
		global.Errlog.Infof("获取游客登录信息失败 key=guest_error")
		utils.FailEncrypt(c, err, "")
		return
	}

	global.Requestlog.Infof("记录当前请求游客返回数据： results=%+v", userInfo)

	//************************调用小米上报数据状态信息 start************************
	deviceType := utils.GetRequestHeaderByName(c, "Os")       //端号
	packageName := utils.GetRequestHeaderByName(c, "Package") //获取包名
	mark := utils.GetRequestHeaderByName(c, "Mark")           //获取设备渠道号
	//fmt.Println("deviceType|packageName|mark ...", deviceType, packageName, mark)
	//处理小米上报
	common_service.AsyncXiaomiReportEvent(c, deviceType, packageName, mark, 1)
	//************************调用小米上报数据状态信息 end**************************

	//调用vivo上报
	if strings.Contains(mark, "vivo") != false {
		fmt.Println("开始上报vivo")
		oaid := utils.GetRequestHeaderByName(c, "oaid")
		imei := utils.GetRequestHeaderByName(c, "imei")
		common_service.VivoReportedEvent(c, packageName, mark, oaid, imei, config.VIVO_ACTIVATION)
	}

	//调用神马平台上报
	if strings.Contains(mark, "sm") != false {
		global.SmClicklog.Info("开始上报sm【神马平台】")
		oaid := utils.GetRequestHeaderByName(c, "oaid")
		imei := utils.GetRequestHeaderByName(c, "imei")
		//fmt.Println(oaid, imei)
		cvType := config.SM_ACTIVATION //激活类型
		channel := config.SM_CHALLEN   //渠道标识
		common_service.SmReportEvent(c, packageName, mark, oaid, imei, channel, cvType)
	}

	res := gin.H{
		"user":       userInfo,
		"token":      token,
		"expireTime": expireTime,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (user *User) Register(c *gin.Context) {
	var req models.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	token, expireTime, err := user_service.Register(c, &req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"token":      token,
		"expireTime": expireTime,
	}
	utils.SuccessEncrypt(c, res, "注册成功~")
}

func (user *User) Login(c *gin.Context) {
	var req models.LoginReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}
	token, expireTime, err := user_service.Login(c, &req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"token":      token,
		"expireTime": expireTime,
	}
	utils.SuccessEncrypt(c, res, "登陆成功~")
}

func (user *User) Logoff(c *gin.Context) {
	var req models.LogoffReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		utils.Fail(c, nil, "获取登陆用户信息失败")
		return
	}
	req.UserId = userId.(int64)

	err := user_service.Logoff(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "注销账号成功")
}

func (user *User) Info(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		utils.FailEncrypt(c, nil, "获取登陆用户信息失败")
		return
	}

	userId := userIdStr.(int64)
	userInfo, err := user_service.UserInfo(userId)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"user": userInfo,
	}
	utils.SuccessEncrypt(c, res, "获取用户信息成功")
}

func (user *User) Edit(c *gin.Context) {
	var req models.EditUserReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.Fail(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		utils.Fail(c, nil, "获取登陆用户信息失败")
		return
	}

	req.UserId = userId.(int64)
	fmt.Println("用户登录的相关信息", req.UserId)
	err := user_service.EditUser(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	utils.Success(c, "", "修改成功")
}

func (user *User) Follow(c *gin.Context) {
	var req models.FollowUserReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.Fail(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		utils.Fail(c, nil, "获取登陆用户信息失败")
		return
	}
	req.UserId = userId.(int64)
	err := user_service.FollowUser(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	var msg string
	if req.FollowType > 0 {
		msg = "关注成功"
	} else {
		msg = "取消关注成功"
	}
	utils.Success(c, "", msg)
}

func (user *User) FollowList(c *gin.Context) {
	var req models.FollowListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, err := user_service.FollowList(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.Success(c, res, "获取关注列表成功")
}

func (user *User) BindRegistId(c *gin.Context) {
	var req models.BindRegistIdReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.Fail(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := user_service.BindRegistId(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	utils.Success(c, "", "ok")
}

func (user *User) MyInvitRewards(c *gin.Context) {
	var req models.MyInvitRewardsReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.Fail(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	inviteUserCount, inviteUserCion, err := user_service.GetMyInvitRewards(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	res := gin.H{
		"inviteUserCount": inviteUserCount,
		"inviteUserCion":  inviteUserCion,
	}
	utils.Success(c, res, "ok")
}
