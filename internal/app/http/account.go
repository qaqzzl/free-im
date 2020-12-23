package http

import (
	"encoding/json"
	"fmt"
	"free-im/config"
	"free-im/internal/app/service"
	"free-im/pkg/util"
	"io/ioutil"
	"net/http"
)

var AccountService *service.AccountService
var CommonService *service.CommonService

// 手机号登录 / 注册
func PhoneLogin(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var err error
	// 验证短信验证码是否正确
	if verify, err := CommonService.IsPhoneVerifyCode(formData["phone"].(string), formData["verify_code"].(string)); err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	} else if verify == false {
		util.RespFail(writer, "短信验证码错误")
		return
	}
	ret, err := AccountService.Login(formData["phone"].(string), "phone","",nil)
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}

	util.RespOk(writer, ret, "")
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

// QQ登陆
func QQLogin(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var err error
	resp, err := http.Get("https://graph.qq.com/user/get_user_info?access_token=" +formData["access_token"].(string) +
		"&oauth_consumer_key=" + config.CommonConf.QQAuthAppID +
		"&openid="+formData["openid"].(string))
	if err != nil {
		util.RespFail(writer, "系统繁忙")
		return
	}
	defer resp.Body.Close()
	authBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		util.RespFail(writer, "系统繁忙")
		return
	}
	// 初始化请求变量结构
	authData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.Unmarshal(authBody, &authData)
	fmt.Println(string(authBody))
	fmt.Println(authData)
	// 获取unionid https://graph.qq.com/oauth2.0/me?access_token=ACCESSTOKEN&unionid=1

	//ret, err := AccountService.Login(formData["phone"].(string), "qq_","",nil)
	//if err != nil {
	//	util.RespFail(writer, "系统繁忙")
	//	return
	//}

	util.RespOk(writer, nil, "ok")
}
