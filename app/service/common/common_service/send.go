package common_service

import (
	"fmt"
	"go-novel/global"
	"go-novel/utils"
	"log"
)

// 发送腾讯云短信验证码
func TencentSmsSend(phone string) (smsCode string, err error) {
	if phone == "" {
		err = fmt.Errorf("%v", "手机号不能为空")
		return
	}
	//检查手机号
	isPhone := utils.CheckMobile(phone)
	if isPhone == false {
		err = fmt.Errorf("%v", "手机号格式不正确")
		return
	}

	//拼接手机号，直接按照+86来进行发送拼接
	countryCode := "86"
	mobile := "+" + countryCode + phone
	log.Printf("send-mobile:%v", mobile)
	//判断是否30S内发送过短信了
	if v, ok := utils.SmsCode.Load(phone); ok {
		sms := v.(*utils.Sms)
		now := utils.GetUnix()
		if now-sms.CreateAt < 30 {
			err = fmt.Errorf("%v", "发送过于频繁，稍后再试")
			return
		}
	}

	smsCode = utils.GenValidateCode(6) //获取验证码
	//定义参数
	sms := utils.Sms{
		AreaCode: countryCode,
		Mobiles:  phone,
		Msg:      smsCode,
		CreateAt: utils.GetUnix(),
	}
	log.Printf("send-mobile-param:%#v", sms)
	//使用腾讯云发送验证验证
	var res = ""
	res, err = sms.SendSmsTencent()
	if err != nil {
		err = fmt.Errorf("%v", "短信发送失败")
		return
	}
	//存储短信信息
	utils.SmsCode.Store(phone, &sms)
	return res, nil
}

func TelSend(phone string) (smsCode string, err error) {
	if phone == "" {
		err = fmt.Errorf("%v", "手机号不能为空")
		return
	}

	isPhone := utils.CheckMobile(phone)
	if isPhone == false {
		err = fmt.Errorf("%v", "手机号格式不正确")
		return
	}

	if v, ok := utils.SmsCode.Load(phone); ok {
		sms := v.(*utils.Sms)
		now := utils.GetUnix()
		if now-sms.CreateAt < 30 {
			err = fmt.Errorf("%v", "发送过于频繁，稍后再试")
			return
		}
	}

	//短信验证码
	smsCode = utils.GenValidateCode(6)

	sms := utils.Sms{
		AreaCode: "86",
		Mobiles:  phone,
		Msg:      smsCode,
		CreateAt: utils.GetUnix(),
	}

	//发送短信验证码
	_, err = sms.SendSMS()
	if err != nil {
		err = fmt.Errorf("%v", "短信发送失败")
		return
	}

	//存储短信信息
	utils.SmsCode.Store(phone, &sms)
	return
}

func EmailSend(email string) (emailCode string, err error) {
	if email == "" {
		err = fmt.Errorf("%v", "邮箱不能为空")
		return
	}

	isEmail := utils.CheckEmail(email)
	if isEmail == false {
		err = fmt.Errorf("%v", "邮箱格式不正确")
		return
	}

	if v, ok := utils.EmailCode.Load(email); ok {
		obj := v.(*utils.EmailObj)
		now := utils.GetUnix()
		if now-obj.CreateAt < 30 {
			err = fmt.Errorf("%v", "发送过于频繁，稍后再试")
			return
		}
	}

	//短信验证码
	emailCode = utils.GenValidateCode(6)

	oEmail := utils.EmailObj{
		Email:    email,
		Msg:      emailCode,
		CreateAt: utils.GetUnix(),
	}

	//发送邮箱验证码
	err = oEmail.SendEmail()
	if err != nil {
		global.Errlog.Errorf("%v", err.Error())
		err = fmt.Errorf("%v", "邮箱发送失败")
		return
	}

	utils.EmailCode.Store(email, &oEmail)
	return
}
