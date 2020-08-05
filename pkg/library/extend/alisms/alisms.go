package alisms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type Sms struct {
	AccessKeyId string
	AccessSecret string
}

// 发送短信验证码
// https://help.aliyun.com/document_detail/101414.html
func (s *Sms) SendSms(PhoneNumbers string, 		//手机号
	SignName string, 	//短信签名名称
	TemplateCode string,	//短信模板ID
	TemplateParam string) (*dysmsapi.SendSmsResponse, error) {	//短信模板变量对应的实际值，JSON格式
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", s.AccessKeyId, s.AccessSecret)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = PhoneNumbers
	request.SignName = SignName
	request.TemplateCode = TemplateCode
	request.TemplateParam = TemplateParam

	response, err := client.SendSms(request)
	return response, err
}
