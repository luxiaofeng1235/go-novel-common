/*
 * @Descripttion: 可选能力占位（默认关闭：OSS/Qiniu 上传、短信验证码等）
 * @Author: red
 * @Date: 2026-01-12 10:25:00
 * @LastEditors: red
 * @LastEditTime: 2026-01-12 10:25:00
 */
package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
)

// 说明：
// 这些能力（OSS/Qiniu 上传、短信超级验证码等）在当前仓库版本中未纳入实现。
// 为保证默认启动可编译运行，这里提供最小 stub：默认关闭，调用会返回明确错误。

// 上传开关：默认关闭（走本地存储逻辑）。
var (
	OssUpload   = false
	QiNiuUpload = false
)

func UploadOss(_ *multipart.FileHeader, _ string) (string, error) {
	return "", errors.New("OSS 上传未启用/未实现")
}

func UploadQiNiu(_ *multipart.FileHeader, _ string) (string, error) {
	return "", errors.New("七牛上传未启用/未实现")
}

// 超级验证码：默认关闭。
var (
	SuperCodeOpenStatus = false
	SuperCode           = ""
)

// 短信验证码的内存缓存（仅用于本地/默认实现）。
var (
	SmsCode sync.Map
)

type Sms struct {
	AreaCode string
	Mobiles  string
	Msg      string
	CreateAt int64
}

// SendSMS / SendSmsTencent：默认不发送，仅返回错误。
func (s *Sms) SendSMS() (string, error) {
	return "", fmt.Errorf("短信发送未配置/未实现")
}

func (s *Sms) SendSmsTencent() (string, error) {
	return "", fmt.Errorf("腾讯云短信发送未配置/未实现")
}

// IsYzm 校验短信验证码（默认实现：仅校验内存缓存中的验证码是否存在且未过期）。
func IsYzm(phone, code string) error {
	if v, ok := SmsCode.Load(phone); ok {
		obj := v.(*Sms)
		now := GetUnix()
		if code != obj.Msg || now-obj.CreateAt > 600 {
			return fmt.Errorf("%v", "验证码不正确或已过期")
		}
		return nil
	}
	return fmt.Errorf("%v", "验证码不正确或已过期!")
}

// GetSmsSendCode 解析短信服务商的下发结果。
// 当前默认 stub：只要返回值非空就视为成功（兼容 common.go 的 "Ok" 判断）。
func GetSmsSendCode(result string) string {
	if result == "" {
		return ""
	}
	return "Ok"
}
