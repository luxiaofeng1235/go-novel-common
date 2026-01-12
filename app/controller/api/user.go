/*
 * @Descripttion: API 用户控制器（登录/注册/用户信息等）
 * @Author: congz
 * @Date: 2020-07-15 14:48:46
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 11:10:00
 */
package api

import (
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/utils"

	"github.com/gin-gonic/gin"
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
	//游客登录判断
	userInfo, token, expireTime, err := user_service.GuestLogin(c, &req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
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

// 登录流程
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
