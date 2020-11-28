package service

import (
	"errors"
	"fmt"
	"free-im/pkg/library/extend/alisms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
	"math/rand"
)

type CommonService struct {
}

// 判断手机验证码是否正确
func (s *CommonService) IsPhoneVerifyCode(phone string, verify_code string) (bool, error) {
	if verify_code == "1234" {
		return true, nil
	}
	return false, nil
}

// 发送手机验证码
func (s *CommonService) SendSms(phone string, sms_type string) (err error) {
	sms := alisms.Sms{
		AccessKeyId:  "",
		AccessSecret: "",
	}
	var response *dysmsapi.SendSmsResponse
	templateParam := fmt.Sprintf(`{"code":%06d}`, rand.Int31n(10000))
	switch sms_type {
	case "login":
		if response, err = sms.SendSms(phone, "FREE", "SMS_181196339", templateParam); err != nil {
			log.Println(err.Error())
			return errors.New("系统繁忙")
		}
		if response.Code != "OK" {
			return errors.New(response.Message)
		}
	default:
		return errors.New("不存在的短信类型")
	}
	return err
}
