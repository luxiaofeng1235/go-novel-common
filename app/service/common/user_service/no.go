package user_service

//func Register(req *models.RegisterReq) (userId int64, err error) {
//	tel := strings.TrimSpace(req.Tel)
//	code := strings.TrimSpace(req.Code)
//	passwd := strings.TrimSpace(req.Passwd)
//	referrer := strings.TrimSpace(req.Referrer)
//	deviceid := strings.TrimSpace(req.Deviceid)
//
//	if tel == "" {
//		err = fmt.Errorf("%v", "手机号不能为空")
//		return
//	}
//	if code == "" {
//		err = fmt.Errorf("%v", "验证码不能为空")
//		return
//	}
//	if deviceid == "" {
//		err = fmt.Errorf("%v", "设备id不能为空")
//		return
//	}
//	if len(passwd) < 6 {
//		err = fmt.Errorf("%v", "密码长度不得少于6位")
//		return
//	}
//	if !utils.CheckMobile(tel) {
//		err = fmt.Errorf("%v", "手机号格式错误")
//		return
//	}
//	err = checkTelCode(tel, code)
//	if err != nil {
//		return
//	}
//
//	count := GetUserCountByTel(tel)
//	if count > 0 {
//		err = fmt.Errorf("%v", "该手机已注册")
//		return
//	}
//	var user models.McUser
//
//	//注册赠送vip天数
//	if utils.User_Reg_Vip_Day > 0 {
//		user.Viptime = utils.GetUnix() + 86400*utils.User_Reg_Vip_Day
//	}
//
//	var hashpwd string
//	hashpwd, err = utils.HashPassword(passwd)
//	if err != nil {
//		return
//	}
//
//	user.Addtime = utils.GetUnix()
//	user.Username = generateName(user.Addtime)
//	user.Nickname = rand.GetRand().ChineseName()
//	user.ParentId, user.ParentLink, err = GetParentLinkByReffer(referrer)
//	if err != nil {
//		return
//	}
//	user.Invitation = utils.RandomString("code", 6)
//
//	user.Status = 1
//	user.Tel = tel
//	user.Deviceid = deviceid
//
//	user.Passwd = hashpwd
//	user.Cion = 0
//	user.Vip = 0
//	user.Viptime = 0
//
//	// 创建用户
//	if err = global.DB.Create(&user).Error; err != nil {
//		return
//	}
//
//	userId = user.Id
//
//	if referrer != "" {
//		var parentUser *models.McUser
//		parentUser, err = GetUserByReferrer(referrer)
//		if err != nil {
//			return
//		}
//		if parentUser.Id <= 0 {
//			err = fmt.Errorf("%v", "邀请码不存在")
//			return
//		}
//		var isReward bool
//		isReward, err = task_service.AppTaskReward(5, parentUser)
//		if !isReward {
//			if err != nil {
//				global.Errlog.Errorf("注册失败 发放邀请奖励失败 referrer=%v userId=%v err=%v", referrer, userId, err.Error())
//			}
//			return
//		}
//	}
//	return
//}
//
//func BindPhone(req *models.BindPhoneReq) (err error) {
//	tel := strings.TrimSpace(req.Tel)
//	code := strings.TrimSpace(req.Code)
//	passwd := strings.TrimSpace(req.Passwd)
//	referrer := strings.TrimSpace(req.Referrer)
//	deviceid := strings.TrimSpace(req.Deviceid)
//	userId := req.UserId
//
//	if tel == "" {
//		err = fmt.Errorf("%v", "手机号不能为空")
//		return
//	}
//	if code == "" {
//		err = fmt.Errorf("%v", "验证码不能为空")
//		return
//	}
//	if len(passwd) < 6 {
//		err = fmt.Errorf("%v", "密码长度不得少于6位")
//		return
//	}
//	if !utils.CheckMobile(tel) {
//		err = fmt.Errorf("%v", "手机号格式错误")
//		return
//	}
//	err = checkTelCode(tel, code)
//	if err != nil {
//		return
//	}
//
//	count := GetUserCountByTel(tel)
//	if count > 0 {
//		err = fmt.Errorf("%v", "该手机号已被其他用户绑定")
//		return
//	}
//	var user *models.McUser
//	user, err = GetGuestUserById(userId)
//	if user.Id <= 0 {
//		err = fmt.Errorf("%v", "该游客不存在")
//		return
//	}
//	if user.IsGuest != 1 {
//		err = fmt.Errorf("%v", "该用户已经注册")
//		return
//	}
//	data := make(map[string]interface{})
//	//注册赠送vip天数
//	if utils.User_Reg_Vip_Day > 0 {
//		data["viptime"] = utils.GetUnix() + 86400*utils.User_Reg_Vip_Day
//		data["vip"] = 1
//	}
//
//	parentId, parentLink, err := GetParentLinkByReffer(referrer)
//	if err != nil {
//		return
//	}
//	data["parent_id"] = parentId
//	data["parent_link"] = parentLink
//
//	var hashpwd string
//	hashpwd, err = utils.HashPassword(passwd)
//	if err != nil {
//		return
//	}
//
//	data["addtime"] = utils.GetUnix()
//	if user.Username == "" {
//		data["username"] = generateName(user.Addtime)
//	}
//	if user.Nickname == "" {
//		data["nickname"] = rand.GetRand().ChineseName()
//	}
//	data["is_guest"] = 0
//	data["tel"] = tel
//	data["passwd"] = hashpwd
//	data["deviceid"] = deviceid
//
//	err = global.DB.Model(models.McUser{}).Where("id = ?", userId).Updates(data).Error
//	if err != nil {
//		global.Sqllog.Errorf("%v", err.Error())
//		return
//	}
//
//	if referrer != "" {
//		var parentUser *models.McUser
//		parentUser, err = GetUserByReferrer(referrer)
//		if err != nil {
//			return
//		}
//		var isReward bool
//		isReward, err = task_service.AppTaskReward(2, parentUser)
//		if !isReward {
//			if err != nil {
//				global.Errlog.Errorf("注册失败 发放邀请奖励失败 reffer=%v userId=%v err=%v", referrer, userId, err.Error())
//			}
//			return
//		}
//	}
//	return
//}

//func ForgotPasswd(req *models.ForgotLoginPasswdReq) (err error) {
//	tel := strings.TrimSpace(req.Tel)
//	pass := strings.TrimSpace(req.Pass)
//	code := strings.TrimSpace(req.Code)
//	if tel == "" {
//		err = fmt.Errorf("%v", "手机号不能为空")
//		return
//	}
//	if pass == "" {
//		err = fmt.Errorf("%v", "密码不能为空")
//		return
//	}
//	if code == "" {
//		err = fmt.Errorf("%v", "验证码不能为空")
//		return
//	}
//	if len(pass) < 6 {
//		err = fmt.Errorf("%v", "密码长度不得少于6位")
//		return
//	}
//	if !utils.CheckMobile(tel) {
//		err = fmt.Errorf("%v", "手机号格式错误")
//		return
//	}
//
//	count := GetUserCountByTel(tel)
//	if count <= 0 {
//		err = fmt.Errorf("%v", "手机不存在，请先注册~")
//		return
//	}
//
//	err = checkTelCode(tel, code)
//	if err != nil {
//		return
//	}
//
//	//新密码
//	var passwd string
//	passwd, err = utils.HashPassword(req.Pass)
//	if err != nil {
//		err = fmt.Errorf("%v", "注册失败,请稍后再试")
//		return err
//	}
//	err = UpdatePassByTel(tel, passwd)
//	if err != nil {
//		return
//	}
//	return err
//}
