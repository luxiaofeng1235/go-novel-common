package utils

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"sync"
)

var (
	EmailCode sync.Map //验证码
)

type EmailObj struct {
	Email    string `json:"email" binding:"required"`
	Msg      string `json:"msg" binding:"required"`
	CreateAt int64  `json:"create_at"`
}

func (email *EmailObj) SendEmail() (err error) {
	return SendEmail(email.Email, email.Msg)
}

//google
//接收邮件 (IMAP) 服务器	imap.gmail.com 要求SSL：是 端口：993
//发送邮件 (SMTP) 服务器	smtp.gmail.com TLS 465 / STARTTLS 587

func SendEmail(to, code string) (err error) {
	//var host string = "smtp.gmail.com"
	//var port int = 587
	//var username string = "hi2210331918@gmail.com"
	//var passwd string = "cotbsyqyenlgfwdy"

	var host string = "smtp.sg.aliyun.com"
	var port int = 465
	//username := "admin@yaodudushu.site"
	//password := "rens123.123"

	username := "services@kuaiyanpingshu.com"
	password := "rens456.456"

	//username := "support@yaodudushu.site"
	//password := "rens678.678"
	// 1. 首先构建一个 Message 对象，也就是邮件对象
	msg := gomail.NewMessage()
	// 2. 填充 From，注意第一个字母要大写
	msg.SetHeader("From", msg.FormatAddress(username, "快眼云阁APP"))
	// 3. 填充 To
	msg.SetHeader("To", to)
	// 5. 设置邮件标题
	msg.SetHeader("Subject", fmt.Sprintf("验证码-%v", code))
	// 6. 设置要发送的邮件正文
	// 第一个参数是类型，第二个参数是内容
	// 如果是 html，第一个参数则是 `text/html`
	text := fmt.Sprintf("验证码为：%v，验证码将在10分钟后失效。请及时使用。如果非本人操作请忽略,有任何疑问与我们联系。", code)
	msg.SetBody("text/html", text)
	// 7. 添加附件，注意，这个附件是完整路径
	//msg.Attach("/Users/yufei/Downloads/1.jpg")
	// 8. 创建 Dialer
	dialer := gomail.NewDialer(host, port, username, password)
	// 9. 禁用 SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// 10. 发送邮件
	return dialer.DialAndSend(msg)
}

// 效验邮箱验证码是否存在
func IsEmailCode(email string, code string) (err error) {
	//效验验证码
	if v, ok := EmailCode.Load(email); ok {
		obj := v.(*EmailObj)
		now := GetUnix()
		if code != obj.Msg || now-obj.CreateAt > 600 {
			return fmt.Errorf("%v", "验证码不正确或已过期")
		}
	} else {
		return fmt.Errorf("%v", "验证码不正确或已过期!")
	}
	return nil
}
