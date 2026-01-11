package main

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func main() {
	fmt.Println(3334)
	return
	// 替换为你的SecretId和SecretKey
	credential := common.NewCredential(
		"YOUR_SECRET_ID",
		"YOUR_SECRET_KEY",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	request := sms.NewSendSmsRequest()

	// 替换为你的应用 ID
	request.SmsSdkAppId = common.StringPtr("YOUR_SDK_APP_ID")
	// 替换为你的签名内容
	request.SignName = common.StringPtr("YOUR_SIGN_NAME")
	// 替换为你的模板 ID
	request.TemplateId = common.StringPtr("YOUR_TEMPLATE_ID")

	// 假设验证码是1234，有效期是5分钟
	request.TemplateParamSet = common.StringPtrs([]string{"1234", "5"})

	// 替换为目标手机号码
	request.PhoneNumberSet = common.StringPtrs([]string{"+8613800138000"})

	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())
}
