package service

import (
	"errors"
	"fmt"
	"free-im/config"
	"free-im/internal/http_app/dao"
	"free-im/pkg/library/extend/alisms"
	"free-im/pkg/logger"
	"free-im/pkg/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"math/rand"
	"time"
)

type CommonService struct {
}

// 判断手机验证码是否正确
func (s *CommonService) IsPhoneVerifyCode(phone string, verify_code string, sms_type string) (bool, error) {
	if !util.PhoneVerify(phone) {
		return false, errors.New("手机号格式错误")
	}

	if verify_code == "2020" {
		return true, nil
	}

	verify_code_data, _ := dao.Dao.Ris().Do("GET", "sms_verify_code:"+sms_type+":"+phone)
	if verify_code_data == nil {
		return false, nil
	}
	if string(verify_code_data.([]uint8)) == verify_code {
		dao.Dao.Ris().Do("DEL", "sms_verify_code:"+sms_type+":"+phone)
		return true, nil
	}
	return false, nil
}

// 发送手机验证码
func (s *CommonService) SendSms(phone string, sms_type string) (err error) {
	if !util.PhoneVerify(phone) {
		return errors.New("手机号格式错误")
	}
	return errors.New("默认验证码: 2020")
	timeStr := time.Now().Format("2006-01-02")
	if sms_total_send_sum, err := dao.Dao.Ris().Do("INCR", "sms_total_send_sum:"+timeStr); err != nil {
		return errors.New("系统繁忙")
	} else {
		dao.Dao.Ris().Do("EXPIRE", "sms_total_send_sum:"+timeStr, 3600*24)
		if sms_total_send_sum.(int64) > 100 {
			return errors.New("短信通道今日关闭,请使用QQ登陆")
		}
	}
	sms := alisms.Sms{
		AccessKeyId:  config.CommonConf.AliYunSmsAccessKeyID,
		AccessSecret: config.CommonConf.AliYunSmsAccessKeySecret,
	}
	var response *dysmsapi.SendSmsResponse
	code := rand.Int31n(9000) + 1000
	templateParam := fmt.Sprintf(`{"code":%06d}`, code)
	switch sms_type {
	case "login":
		if response, err = sms.SendSms(phone, "Free", "SMS_181196339", templateParam); err != nil {
			logger.Sugar.Error(err)
			err = errors.New("系统繁忙")
		}
		if response.Code != "OK" {
			err = errors.New(response.Message)

		}
	default:
		err = errors.New("不存在的短信类型")
	}
	if err != nil {
		logger.Sugar.Error(err)
		dao.Dao.Ris().Do("DECR", "sms_total_send_sum:"+timeStr)
	} else {
		dao.Dao.Ris().Do("SET", "sms_verify_code:"+sms_type+":"+phone, code)
		dao.Dao.Ris().Do("EXPIRE", "sms_verify_code:"+sms_type+":"+phone, 60*30)
	}
	return err
}
