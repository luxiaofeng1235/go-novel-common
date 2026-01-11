package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/app/models"
	"go-novel/app/service/admin/admin_service"
	"go-novel/app/service/admin/menu_service"
	"go-novel/app/service/admin/monitor_service"
	"go-novel/app/service/admin/role_service"
	"go-novel/app/service/common/upload_service"
	"go-novel/utils"
	"log"
	"strings"
)

type Index struct{}

func (index *Index) IndexGet(c *gin.Context) {
	c.Writer.Write([]byte("首页测试"))
}

// 登录验证码
func (index *Index) Captcha(c *gin.Context) {
	idKeyC, base64stringC := utils.GetVerifyImgString()
	utils.Success(c, gin.H{"idKeyC": idKeyC, "base64stringC": base64stringC}, "ok")
}

// 管理员登录
func (index *Index) Login(c *gin.Context) {
	var login models.Login
	if err := c.ShouldBindJSON(&login); err != nil {
		utils.Fail(c, err, "数据解析失败")
		return
	}

	ip := "" //登录暂时不需要记录IP了，boss交代
	userAgent := c.Request.Header.Get("User-Agent")

	if token, _, err := admin_service.Login(&login, ip); err != nil {
		if strings.Contains(err.Error(), "密码错误") {
			go monitor_service.LoginLog(0, login.Username, ip, userAgent, err.Error(), "系统后台")
		}
		utils.Fail(c, err, "")
		return
	} else {
		go func() {
			//创建登录成功日志
			monitor_service.LoginLog(1, login.Username, ip, userAgent, "登录成功", "系统后台")
		}()
		//返回结果
		utils.Success(c, gin.H{"token": token}, "登陆成功")
		return
	}
}

// 刷新token
func (index *Index) RefreshToken(c *gin.Context) {
	token := c.Request.Header.Get("Token")
	token, _, err := utils.RefreshToken(token)
	if err != nil {
		utils.Fail(c, err, "刷新token失败")
		return
	}
	utils.Success(c, gin.H{"token": token}, "ok")
}

// 退出登录
func (index *Index) Logout(c *gin.Context) {
	username, _ := c.Get("username")
	utils.Success(c, username, "ok")
	return
}

// 上传管理员头像

func (index *Index) Avatar(c *gin.Context) {
	url, err := upload_service.UploadFile(c, "avatar", "")
	if err != nil {
		utils.Fail(c, err, "上传头像失败")
		return
	}
	username, _ := c.Get("username")
	isAvatar, err := admin_service.EditAvatar(username.(string), url)
	if err != nil || isAvatar == false {
		utils.Fail(c, err, "上传头像失败")
		return
	}
	//上传回调完成后，获取新的url地址
	remote_url := utils.GetAdminUrl() + url
	utils.Success(c, gin.H{
		"url": remote_url,
	}, "ok")
	return
}

// 获取当前登录用户信息
func (index *Index) Profile(c *gin.Context) {
	//获取用户信息
	userInfo, err := admin_service.GetCurrentUserInfo(c)
	if err != nil {
		utils.Fail(c, err, "获取当前登录用户信息失败")
		return
	}

	utils.Success(c, userInfo, "ok")
	return
}

// 更新用户登录密码
func (index *Index) UpdatePwd(c *gin.Context) {
	var req models.ChangePwdReq

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "数据解析失败")
		return
	}

	username, _ := c.Get("username")
	if isUpdate, err := admin_service.UpdatePwd(&req, username.(string)); err != nil || isUpdate == false {
		utils.Fail(c, err, "修改密码失败")
		return
	} else {
		//返回结果
		utils.Success(c, "", "修改密码成功")
		return
	}
}

// 修改用户信息
func (index *Index) EditProfile(c *gin.Context) {
	var req models.EditProfile

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "数据解析失败")
		return
	}

	username, _ := c.Get("username")
	if isUpdate, err := admin_service.EditProfile(&req, username.(string)); err != nil || isUpdate == false {
		utils.Fail(c, err, "修改信息失败")
		return
	} else {
		//返回结果
		utils.Success(c, "", "修改信息成功")
		return
	}
}

// 获取管理员信息
func (index *Index) GetUserInfo(c *gin.Context) {
	username, _ := c.Get("username")
	err, user, roles, permissions, rolesInfo := admin_service.GetUserInfo(username.(string))
	if err != nil {
		utils.Fail(c, err, "获取数据失败")
		return
	}

	//返回结果
	mapData := gin.H{
		"user":        user,
		"roles":       roles,
		"permissions": permissions,
		"rolesInfo":   rolesInfo,
	}
	utils.Success(c, mapData, "ok")
	return
}

// 获取路由信息
func (index *Index) GetRouters(c *gin.Context) {
	username, _ := c.Get("username")
	admininfo, err := admin_service.GetAdminByUsername(username.(string))
	if err != nil {
		utils.Fail(c, err, "获取用户信息失败")
		return
	}
	isSuperAdmin := false
	var menuList []*models.SysAuthRule
	if admininfo != nil {
		userId := admininfo.Id
		//获取无需验证权限的用户id
		NotCheckAuthAdminIds := viper.GetIntSlice("adminInfo.notCheckAuthAdminIds")
		for _, v := range NotCheckAuthAdminIds {
			if int64(v) == userId {
				isSuperAdmin = true
				break
			}
		}

		//获取用户角色信息
		allRoles, err := role_service.GetRoleList()
		if err == nil {
			roles, err := admin_service.GetAdminRole(userId, allRoles)
			if err == nil {
				name := make([]string, len(roles))
				roleIds := make([]uint, len(roles))
				for k, v := range roles {
					name[k] = v.Name
					roleIds[k] = uint(v.Id)
				}
				//获取菜单信息
				if isSuperAdmin {
					//超管获取所有菜单
					menuList, _ = menu_service.GetAllMenusTree()
				} else {
					menuList, err = menu_service.GetAdminMenusByRoleIds(roleIds)
				}
				if err != nil {
					utils.Fail(c, err, "获取数据失败")
					return
				}
			} else {
				utils.Fail(c, err, "获取数据失败")
				return
			}
		}
	}
	utils.Success(c, menuList, "ok")

	//list,_:=menu_service.GetMapMenus()
	//c.JSON(http.StatusOK, utils.Success(list, "ok"))
}

// 获取token 测试
func (index *Index) GetToken(c *gin.Context) {
	token, _, err := utils.GenerateToken("admin", "123456", 1)
	if err != nil {
		log.Println(token, err)
		return
	}
	utils.SetCookie(c, "token", token, 86400)

	c.JSON(200, gin.H{"token": token, "error": err})
}

// 获取token 测试
func (index *Index) ReToken(c *gin.Context) {
	token := c.Request.Header.Get("Token")
	token1, err := utils.GetCookie(c, "token")
	log.Println(token1, err)

	claims, err := utils.ParseToken(token)
	log.Printf("claims:%#v", claims)

	token, _, err = utils.RefreshToken(token)
	if err != nil {
		log.Println(token, err)
		return
	}

	claims, err = utils.ParseToken(token)
	log.Printf("claims:%#v", claims)
	utils.SetCookie(c, "token", token, 10)

	c.JSON(200, gin.H{"token": token, "error": err})
}

// 测试是否携带token访问
func (index *Index) CheckToken(c *gin.Context) {
	username, ok := c.Get("username")
	fmt.Println(username, ok)
	token, err := utils.GetCookie(c, "token")
	log.Println(token, err)

	c.JSON(200, gin.H{"code": 200, "msg": "ok", "data": username})
}
