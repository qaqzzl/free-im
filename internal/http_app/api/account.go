package api

import (
	"encoding/json"
	"fmt"
	"free-im/config"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	"free-im/pkg/logger"
	"free-im/pkg/util"
	"io/ioutil"
	"net/http"
	"time"
)

var AccountService *service.AccountService
var CommonService *service.CommonService

// 手机号登录 / 注册
func PhoneLogin(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("PhoneLogin")
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var err error
	// 验证短信验证码是否正确
	if verify, err := CommonService.IsPhoneVerifyCode(formData["phone"].(string), formData["verify_code"].(string), "login"); err != nil {
		fmt.Println(err)
		util.RespFail(writer, err.Error())
		return
	} else if verify == false {
		util.RespFail(writer, "短信验证码错误")
		return
	}
	ret, err := AccountService.Login(formData["phone"].(string), "phone", "", model.UserMember{})
	if err != nil {
		logger.Logger.Info(err.Error())
		util.RespFail(writer, "系统繁忙")
		return
	}

	util.RespOk(writer, ret, "")
}

// QQ登陆
func QQLogin(writer http.ResponseWriter, request *http.Request) {
	// 初始化请求变量结构
	formData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.NewDecoder(request.Body).Decode(&formData)
	var err error
	resp, err := http.Get("https://graph.qq.com/user/get_user_info?access_token=" + formData["access_token"].(string) +
		"&oauth_consumer_key=" + config.CommonConf.QQAuthAppID +
		"&openid=" + formData["openid"].(string))
	if err != nil {
		logger.Sugar.Error(err)
		util.RespFail(writer, "系统繁忙")
		return
	}
	defer resp.Body.Close()
	authBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		logger.Sugar.Error(err)
		util.RespFail(writer, "系统繁忙")
		return
	}
	// 初始化请求变量结构
	authData := make(map[string]interface{})
	// 调用json包的解析，解析请求body
	json.Unmarshal(authBody, &authData)
	// 获取unionid https://graph.qq.com/oauth2.0/me?access_token=ACCESSTOKEN&unionid=1
	UserMember := model.UserMember{}
	UserMember.Nickname = authData["nickname"].(string)
	if authData["gender"].(string) == "女" {
		UserMember.Gender = "w"
	} else {
		UserMember.Gender = "m"
	}
	times, _ := time.Parse("2006", authData["year"].(string))
	timeUnix := times.Unix()
	if timeUnix > 0 {
		UserMember.Birthdate = int(timeUnix)
	}
	UserMember.Avatar = authData["figureurl_qq_2"].(string)
	UserMember.City = authData["city"].(string)
	UserMember.Province = authData["province"].(string)
	ret, err := AccountService.Login(formData["openid"].(string), "qq_openid", formData["access_token"].(string), UserMember)
	if err != nil {
		logger.Sugar.Error(err)
		util.RespFail(writer, "系统繁忙")
		return
	}

	util.RespOk(writer, ret, "ok")
}
