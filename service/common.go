package service

import "errors"

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
	switch sms_type {
	case "login":
	default:
		return errors.New("不存在的短信类型")
	}
	return err
}

