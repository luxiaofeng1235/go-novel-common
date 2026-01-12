package user_service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goroom/rand"
	"go-novel/app/models"
	"go-novel/app/service/api/shelf_service"
	"go-novel/app/service/api/task_service"
	"go-novel/app/service/common/notify_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/app/service/common/user_service"
	"go-novel/config"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func hashLoginPasswd(plain string) string {
	salt := strings.TrimSpace(config.GetString("auth.passwordSalt"))
	if salt == "" {
		return utils.Md5(plain)
	}
	return utils.GetMd5(plain, salt)
}

// 账号密码注册
func Register(c *gin.Context, req *models.RegisterReq) (token string, expireTime int64, err error) {
	username := strings.TrimSpace(req.Username)
	passwd := strings.TrimSpace(req.Passwd)
	nickname := strings.TrimSpace(req.Nickname)
	referrer := strings.TrimSpace(req.Referrer)
	deviceid := strings.TrimSpace(req.Deviceid)

	if username == "" {
		return "", 0, fmt.Errorf("%v", "账号不能为空")
	}
	if passwd == "" {
		return "", 0, fmt.Errorf("%v", "密码不能为空")
	}

	// 用户名唯一性校验
	var count int64
	if err = global.DB.Model(models.McUser{}).Where("username = ?", username).Count(&count).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return "", 0, err
	}
	if count > 0 {
		return "", 0, fmt.Errorf("%v", "账号已存在")
	}

	parentId, parentLink, err := GetParentLinkByReffer(referrer)
	if err != nil {
		return "", 0, err
	}

	mark := utils.GetRequestHeaderByName(c, "Mark")
	devicePackage := utils.GetRequestHeaderByName(c, "Package")
	imei := utils.GetRequestHeaderByName(c, "Imei")
	oaid := utils.GetRequestHeaderByName(c, "Oaid")
	ip := utils.RemoteIp(c)

	if nickname == "" {
		nickname = rand.GetRand().ChineseName()
	}

	user := models.McUser{
		ParentId:   parentId,
		ParentLink: parentLink,
		Username:   username,
		Passwd:     hashLoginPasswd(passwd),
		Nickname:   nickname,
		Invitation: utils.RandomString("code", 6),
		Status:     1,
		IsGuest:    0,
		Deviceid:   deviceid,
		Mark:       mark,
		Oaid:       oaid,
		Package:    devicePackage,
		Ip:         ip,
		Imei:       imei,
		Addtime:    utils.GetUnix(),
	}

	if err = global.DB.Model(models.McUser{}).Create(&user).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return "", 0, err
	}

	// token 中不放明文密码，使用已存储的 hash
	token, expireTime, err = utils.GenerateToken(user.Username, user.Passwd, 1)
	if err != nil {
		global.Errlog.Errorf("注册生成token失败 username=%v err=%v", username, err.Error())
		return "", 0, err
	}
	return token, expireTime, nil
}

func GuestLogin(c *gin.Context, req *models.GuestLoginReq) (userInfo *models.McUser, token string, expireTime int64, err error) {
	deviceid := strings.TrimSpace(req.Deviceid)
	referrer := strings.TrimSpace(req.Referrer)
	sex := req.Sex

	mark := utils.GetRequestHeaderByName(c, "Mark")             //获取渠道号
	devicePackage := utils.GetRequestHeaderByName(c, "Package") //客户端包名
	imei := utils.GetRequestHeaderByName(c, "Imei")             //imei设备号
	oaid := utils.GetRequestHeaderByName(c, "Oaid")             //获取oaid信息
	ip := utils.RemoteIp(c)                                     //获取客户端IP

	global.Requestlog.Info("用户请求参数 oaid = ", oaid, " deviceId = ", deviceid)
	//同时判断oaid和devide_id不能为空
	if deviceid == "" && oaid == "" {
		global.Errlog.Infof("device_empty and oaid_empty,key=device_nil|oaid_nil")
		err = fmt.Errorf("%v", "设备ID或者oaid不能同时为空")
		return
	}

	//userInfo, err = GetUserByDeviceid(deviceid)
	//通过关联设备ID+oaid来进行反查
	userInfo, err = GetUserByDeviceAndOaid(deviceid, oaid, imei)
	if err != nil {
		global.Errlog.Infof("not found user record-info empty")
		return
	}
	userId := userInfo.Id
	username := userInfo.Username
	passwd := userInfo.Passwd
	if userId <= 0 {
		var user models.McUser
		username = utils.GetGuestName()
		user.Username = username
		user.Nickname = rand.GetRand().ChineseName()
		user.ParentId, user.ParentLink, err = GetParentLinkByReffer(referrer)
		if err != nil {
			return
		}
		user.Invitation = utils.RandomString("code", 6)
		user.Status = 1
		user.Cion = 0
		user.Vip = 0
		user.Sex = int(sex)
		user.Viptime = 0
		user.IsGuest = 1
		//判断deviceId不为空
		if deviceid != "" {
			user.Deviceid = deviceid
		}
		user.Addtime = utils.GetUnix()
		user.Mark = mark
		//oaid可能为空需要判断下
		if oaid != "" {
			user.Oaid = oaid
		}
		user.Package = devicePackage
		user.Imei = imei
		user.Ip = ip
		// 创建用户
		if err = global.DB.Debug().Create(&user).Error; err != nil {
			return
		}
		user.Pic = utils.GetFileUrl(user.Pic)
		userInfo = &user
		global.Errlog.Infof("create_userId =%d,key=create_success", userInfo.Id)
	} else {
		global.Errlog.Infof("old_userId,此用户ID = %d 是一个老用户", userInfo.Id)
		userInfo.Passwd = ""
		userInfo.Pic = utils.GetFileUrl(userInfo.Pic)
		//处理下用户的关联更新
		mu := make(map[string]interface{})
		//如果deviceID不为空
		if deviceid != "" {
			mu["deviceid"] = deviceid
		}
		mu["mark"] = mark
		mu["imei"] = imei
		mu["package"] = devicePackage
		if oaid != "" {
			mu["oaid"] = oaid
		}
		mu["ip"] = ip
		//mu["uptime"] = utils.GetUnix()
		mu["last_login_time"] = utils.GetUnix() //获取最后一次的登录时间
		_ = UpdateUserByUserId(userId, mu)      //根据用户ID进行更新
	}
	//直接返回游客信息 并发送token
	token, expireTime, err = utils.GenerateToken(username, passwd, 1)
	if err != nil {
		global.Errlog.Errorf("游客登录生成token失败 username=%v err=%v", username, err.Error())
		return
	}
	return
}

// 用户登录
func Login(c *gin.Context, req *models.LoginReq) (token string, expireTime int64, err error) {
	loginType := req.LoginType
	username := strings.TrimSpace(req.Username)
	tel := strings.TrimSpace(req.Tel)
	email := strings.TrimSpace(req.Email)
	code := strings.TrimSpace(req.Code)
	passwd := strings.TrimSpace(req.Passwd)
	deviceid := strings.TrimSpace(req.Deviceid)
	referrer := strings.TrimSpace(req.Referrer)

	// 账号密码登录（推荐）
	if username != "" {
		if passwd == "" {
			return "", 0, fmt.Errorf("%v", "密码不能为空")
		}
		var user *models.McUser
		err = global.DB.Model(models.McUser{}).Where("username = ?", username).First(&user).Error
		if err != nil || user == nil || user.Id <= 0 {
			return "", 0, fmt.Errorf("%v", "账号不存在，请先注册")
		}
		if user.Status == 0 {
			return "", 0, fmt.Errorf("%v", "账户已被锁定~")
		}
		if user.Passwd == "" || hashLoginPasswd(passwd) != user.Passwd {
			return "", 0, fmt.Errorf("%v", "密码不正确~")
		}

		// 登录时同步更新设备信息（可选）
		mark := utils.GetRequestHeaderByName(c, "Mark")
		devicePackage := utils.GetRequestHeaderByName(c, "Package")
		imei := utils.GetRequestHeaderByName(c, "Imei")
		oaid := utils.GetRequestHeaderByName(c, "Oaid")
		ip := utils.RemoteIp(c)
		mu := map[string]interface{}{
			"mark":            mark,
			"imei":            imei,
			"package":         devicePackage,
			"ip":              ip,
			"last_login_time": utils.GetUnix(),
			"uptime":          utils.GetUnix(),
		}
		if deviceid != "" {
			mu["deviceid"] = deviceid
		}
		if oaid != "" {
			mu["oaid"] = oaid
		}
		_ = UpdateUserByUserId(user.Id, mu)

		token, expireTime, err = utils.GenerateToken(user.Username, user.Passwd, 1)
		if err != nil {
			global.Errlog.Errorf("登录生成token失败 username=%v err=%v", username, err.Error())
			return "", 0, err
		}
		return token, expireTime, nil
	}

	if tel == "" && email == "" {
		err = fmt.Errorf("%v", "手机号或邮箱不能为空")
		return
	}
	if deviceid == "" {
		err = fmt.Errorf("%v", "设备ID不能为空")
		return
	}

	if loginType > 0 {
		if code == "" {
			err = fmt.Errorf("%v", "验证码不能为空")
			return
		}
		if tel != "" {
			err = checkTelCode(tel, code)
		} else {
			err = checkEmailCode(email, code)
		}
		if err != nil {
			return
		}
	}

	var count int64
	if tel != "" {
		count = GetUserCountByTel(tel)
	} else {
		count = GetUserCountByEmail(email)
	}

	mark := utils.GetRequestHeaderByName(c, "Mark")             //获取渠道号
	devicePackage := utils.GetRequestHeaderByName(c, "Package") //客户端包名
	imei := utils.GetRequestHeaderByName(c, "Imei")             //imei设备号
	oaid := utils.GetRequestHeaderByName(c, "Oaid")             //获取oaid信息
	ip := utils.RemoteIp(c)                                     //获取客户端IP
	fmt.Println(mark, devicePackage, imei, oaid)
	if count <= 0 {
		var userInfo *models.McUser
		userInfo, err = GetUserByDeviceAndOaid(deviceid, oaid, imei)
		userId := userInfo.Id
		if userId <= 0 {
			err = fmt.Errorf("%v", "deviceid oaid 不存在")
			return
		}
		userName := utils.GetRandUserName()

		var parentId int64
		var parentLink string
		parentId, parentLink, err = GetParentLinkByReffer(referrer)
		if err != nil {
			return
		}
		if parentId <= 0 {
			if referrer != "" {
				err = fmt.Errorf("%v", "邀请码不存在")
				return
			}
		}

		var nickName, userPic string
		if userInfo.Nickname == "" {
			nickName = rand.GetRand().ChineseName()
		}

		if userInfo.Pic == "" {
			userPic = ""
		}
		userAdd := models.McUser{
			ParentId:   parentId,
			ParentLink: parentLink,
			Username:   userName,
			Nickname:   nickName,
			Deviceid:   deviceid,
			Mark:       mark,
			Oaid:       oaid,
			Package:    devicePackage,
			Ip:         ip,
			Imei:       imei,
			Status:     1,
			IsGuest:    0,
			Tel:        tel,
			Email:      email,
			Addtime:    utils.GetUnix(),
		}
		err = global.DB.Model(models.McUser{}).Create(&userAdd).Error
		if err != nil {
			global.Sqllog.Errorf("%v", err.Error())
			return
		}

		if referrer != userInfo.Invitation {
			err = task_service.InviteGiveReward(parentId, userId, nickName, userPic)
			if err != nil {
				global.Errlog.Errorf("%v", err.Error())
				return
			}
		}

		//发送token
		token, expireTime, err = utils.GenerateToken(userName, passwd, 1)
		if err != nil {
			global.Errlog.Errorf("登录生成token失败 tel=%v err=%v", tel, err.Error())
			return
		}
		return
	} else {
		var userInfo *models.McUser
		//查询是否关联用户的ID
		userInfo, err = GetUserInfoByMailOrTel(email, tel)
		if err != nil {
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
		//只有当设备ID为空的时候才进行更新
		if deviceid != "" || oaid != "" {
			//更新当前的设备信息
			userId := userInfo.Id //用户id
			mu := make(map[string]interface{})
			mu["deviceid"] = deviceid
			mu["uptime"] = utils.GetUnix()
			mu["deviceid"] = deviceid
			mu["mark"] = mark
			mu["imei"] = imei
			mu["package"] = devicePackage
			if oaid != "" {
				mu["oaid"] = oaid
			}
			mu["ip"] = ip
			mu["last_login_time"] = utils.GetUnix() //获取最后一次的登录时间
			_ = UpdateUserByUserId(userId, mu)      //根据用户ID进行更新
		}
	}
	var user *models.McUser
	if tel != "" {
		user, err = user_service.GetUserByTel(tel)
		if err != nil {
			return
		}
	}

	if email != "" {
		user, err = user_service.GetUserByEmail(email)
		if err != nil {
			return
		}
	}

	if loginType <= 0 {
		if passwd == "" {
			err = fmt.Errorf("%v", "密码不能为空")
			return
		}
		isPwd := utils.CheckPasswordHash(passwd, user.Passwd)
		if !isPwd {
			err = fmt.Errorf("%v", "密码不正确~")
			return
		}
	}

	userId := user.Id
	userName := user.Username

	if user.Status == 0 {
		err = fmt.Errorf("%v", "账户已被锁定~")
		return
	}

	mu := make(map[string]interface{})
	if user.Vip > 0 && user.Viptime < utils.GetUnix() {
		mu["vip"] = 0
		mu["viptime"] = 0
		err = global.DB.Model(models.McUser{}).Where("id", userId).Updates(mu).Error
		if err != nil {
			global.Sqllog.Errorf("%v", err.Error())
			return
		}
	}

	//发送token
	token, expireTime, err = utils.GenerateToken(userName, passwd, 1)
	if err != nil {
		global.Errlog.Errorf("登录生成token失败 tel=%v email=%v err=%v", tel, email, err.Error())
		return
	}
	return
}

func Logoff(req *models.LogoffReq) (err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	err = DeleteFollowByUserId(userId)
	if err != nil {
		err = fmt.Errorf("清楚账号信息失败 %v", err.Error())
		return
	}
	err = shelf_service.DeleteShelfByUserId(userId)
	if err != nil {
		err = fmt.Errorf("清楚账号信息失败 %v", err.Error())
		return
	}
	err = DeleteTaskListByUserId(userId)
	if err != nil {
		err = fmt.Errorf("清楚账号信息失败 %v", err.Error())
		return
	}
	err = DeleteCionByUserId(userId)
	if err != nil {
		err = fmt.Errorf("清楚账号信息失败 %v", err.Error())
		return
	}
	err = user_service.UserLogoff(userId)
	return
}

func UserInfo(userId int64) (res *models.GetInfoRes, err error) {
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	var user *models.McUser
	user, err = GetUserById(userId)
	user.Passwd = ""
	user.Pic = utils.GetFileUrl(user.Pic)

	res = new(models.GetInfoRes)
	res.McUser = user
	res.ToDayCion = getTodayCion(userId)
	res.FollowCount = getFollowCount(userId)
	res.FansCount = getFansCount(userId)
	rate, rmb, err := setting_service.GetMoneyByCion(user.Cion)
	if err != nil {
		return
	}
	res.Rate = rate
	res.CionRmb = rmb
	return
}

func EditUser(req *models.EditUserReq) (err error) {
	//类型 tel email nickname sex pic pass
	//type = tel tel code不能为空
	//type = email email code不能为空
	//type = nickname nickname不能为空
	//type = sex sex不能为空
	//type = pic pic不能为空
	rtype := strings.TrimSpace(req.Type)
	tel := strings.TrimSpace(req.Tel)
	code := strings.TrimSpace(req.Code)
	email := strings.TrimSpace(req.Email)
	nickname := strings.TrimSpace(req.Nickname)
	sex := req.Sex
	pic := strings.TrimSpace(req.Pic)
	userId := req.UserId
	bookType := req.BookType

	mu := make(map[string]interface{})
	if rtype == utils.Tel {
		if tel == "" {
			err = fmt.Errorf("%v", "手机号不能为空")
			return
		}
		if code == "" {
			err = fmt.Errorf("%v", "验证码不能为空")
			return
		}
		err = checkTelCode(tel, code)
		if err != nil {
			return
		}
		mu["tel"] = tel
	} else if rtype == utils.Email {
		if email == "" {
			err = fmt.Errorf("%v", "邮箱不能为空")
			return
		}
		if code == "" {
			err = fmt.Errorf("%v", "验证码不能为空")
			return
		}
		err = checkEmailCode(email, code)
		if err != nil {
			return
		}
		mu["email"] = email
	} else if rtype == utils.Nickname {
		if nickname == "" {
			err = fmt.Errorf("%v", "昵称不能为空")
			return
		}
		mu["nickname"] = nickname
	} else if rtype == utils.Sex {
		if sex <= 0 {
			err = fmt.Errorf("%v", "性别不能为空")
			return
		}
		mu["sex"] = sex
	} else if rtype == utils.Pic {
		if pic == "" {
			err = fmt.Errorf("%v", "头像不能为空")
			return
		}
		mu["pic"] = pic
	} else if rtype == utils.BookType {
		if bookType <= 0 {
			err = fmt.Errorf("%v", "阅读喜好不能为空")
			return
		}
		mu["book_type"] = bookType
	} else if rtype == utils.Passwd {
		//if oldPasswd == "" {
		//	err = fmt.Errorf("%v", "原密码不能为空")
		//	return
		//}
		//userpwd := GetUserPasswdIdById(userId)
		//isPwd := utils.CheckPasswordHash(oldPasswd, userpwd)
		//if !isPwd {
		//	err = fmt.Errorf("%v", "原密码不正确~")
		//	return
		//}
		//if passwd == "" {
		//	err = fmt.Errorf("%v", "新密码不能为空")
		//	return
		//}
		//
		//var newPass string
		//newPass, err = utils.HashPassword(passwd)
		//if err != nil {
		//	return
		//}
		//
		//mu["passwd"] = newPass
	}

	err = UpdateUserByUserId(userId, mu)
	user, _ := GetUserById(userId)

	if user.Pic != "" && user.Nickname != "" && user.Email != "" && user.Tel != "" {
		var taskId int64 = 1
		var task *models.McTask
		task, err = task_service.GetTaskById(taskId)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			//这里不返回给前端错误
			err = nil
			return
		}
		err = task_service.CompleteTask(task, userId)
		if err != nil {
			global.Errlog.Errorf("%v", err.Error())
			//这里不返回给前端错误
			err = nil
			return
		}
	}
	return
}

func FollowUser(req *models.FollowUserReq) (err error) {
	followType := req.FollowType
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}
	byUserId := req.ByUserId
	if byUserId <= 0 {
		err = fmt.Errorf("%v", "被关注者用户ID不能为空")
		return
	}
	if byUserId == userId {
		err = fmt.Errorf("%v", "不能关注自己")
		return
	}
	var count, userCount int64
	count = GetFollowCountByUid(userId, byUserId)
	userCount = GetUserCountById(byUserId)
	if userCount <= 0 {
		err = fmt.Errorf("%v", "关注用户不存在")
		return
	}
	if followType <= 0 {
		if count <= 0 {
			err = fmt.Errorf("%v", "未关注对方")
			return
		}
		err = DeleteFollowUid(userId, byUserId)
		return
	} else {
		if count > 0 {
			err = fmt.Errorf("%v", "已经关注对方啦~")
			return
		}
		var follow models.McUserFollow
		follow.Uid = userId
		follow.ByUid = byUserId
		follow.Addtime = utils.GetUnix()
		if err = global.DB.Create(&follow).Error; err != nil {
			return
		}
		user, _ := GetUserById(userId)
		byUser, _ := GetUserById(byUserId)
		//关注通知
		_ = notify_service.SendNotify(utils.Follow, "", user.Pic, userId, byUserId, "关注通知", fmt.Sprintf("%v 关注了您", byUser.Nickname), userId)
	}

	return
}

func FollowList(req *models.FollowListReq) (followList []*models.FollowListRes, err error) {
	followType := req.FollowType
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "请先登录")
		return
	}

	var userIds []int64
	if followType > 0 {
		userIds = GetFollowIdsByUid(userId)
	} else {
		userIds = GetFansIdsByUid(userId)
	}

	var userList []*models.McUser
	userList, err = GetUserListByIds(userIds)
	if err != nil {
		err = fmt.Errorf("获取关注列表错误 %v", err.Error())
		return
	}
	if len(userList) <= 0 {
		return
	}
	for _, user := range userList {
		var isBoth int
		count := GetFollowCountByUid(userId, user.Id)
		if count > 0 {
			isBoth = 1
		}
		follow := &models.FollowListRes{
			Id:       user.Id,
			Nickname: user.Nickname,
			Addtime:  utils.UnixToDatetime(user.Addtime),
			Pic:      utils.GetFileUrl(user.Pic),
			IsBoth:   isBoth,
		}
		followList = append(followList, follow)
	}
	return
}

func BindRegistId(req *models.BindRegistIdReq) (err error) {
	registId := strings.TrimSpace(req.RegistId)
	userId := req.UserId

	if registId == "" {
		err = fmt.Errorf("%v", "推送设备ID不能为空")
		return
	}
	var user *models.McUser
	user, err = GetUserById(userId)
	if user.Id <= 0 {
		err = fmt.Errorf("%v", "该用户不存在")
		return
	}
	err = global.DB.Model(models.McUser{}).Where("id = ?", userId).Update("regist_id", registId).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetMyInvitRewards(req *models.MyInvitRewardsReq) (inviteUserCount, inviteUserCion int64, err error) {
	parentId := req.UserId
	inviteUserCount = GetChildCountById(parentId)
	inviteUserCion = int64(getInviteCion(parentId))
	return
}
