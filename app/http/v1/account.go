package v1

import (
	"encoding/json"
	"fmt"
	"free-im/service"
	"free-im/util"
	"net/http"
	"strconv"
)

var AccountService *service.AccountService
var CommonService *service.CommonService

// 手机号登录 / 注册
func PhoneLogin(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)

	//判断手机号是否存在
	var is_register bool
	var err error
	if is_register,err = AccountService.IsRegister(formData["phone"].(string), "phone"); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	// 验证短信验证码是否正确
	if verify, err := CommonService.IsPhoneVerifyCode(formData["phone"].(string), formData["verify_code"].(string)); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	} else {
		if verify == false {
			util.RespFail(writer, "短信验证码错误")
			return
		}
	}

	var member_id int64
	if is_register {
		member_id, err = AccountService.PhoneLogin(formData["phone"].(string))
	} else {
		member_id, err = AccountService.Register(formData["phone"].(string), "phone", "")
	}
	if err != nil {
		fmt.Println(err)
		util.RespFail(writer, "系统繁忙")
		return
	}

	// 获取token
	var token string
	if token, err = AccountService.GetToken(member_id, "app"); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}

	data := make(map[string]string)
	data["access_token"] = token
	data["uid"] = strconv.Itoa(int(member_id))
	util.RespOk(writer, data, "")
}

// 发送登录短信验证码
func SendLoginSms(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	if err := CommonService.SendSms(formData["phone"].(string), "login"); err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	util.RespOk(writer, nil, "短信验证码发送成功")
}