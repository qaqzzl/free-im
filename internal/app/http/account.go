package http

import (
	"encoding/json"
	"free-im/config"
	"free-im/internal/app/service"
	"free-im/pkg/logger"
	"free-im/pkg/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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
	if verify, err := CommonService.IsPhoneVerifyCode(formData["phone"].(string), formData["verify_code"].(string), "login"); err != nil {
		util.RespFail(writer, err.Error())
		return
	} else if verify == false {
		util.RespFail(writer, "短信验证码错误")
		return
	}
	ret, err := AccountService.Login(formData["phone"].(string), "phone", "", nil)
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
	data := make(map[string]string)
	data["nickname"] = authData["nickname"].(string)
	if authData["gender"].(string) == "女" {
		data["gender"] = "w"
	} else {
		data["gender"] = "m"
	}
	times, _ := time.Parse("2006", authData["year"].(string))
	timeUnix := times.Unix()
	var birthdate = "0"
	if timeUnix > 0 {
		birthdate = strconv.Itoa(int(timeUnix))
	}
	data["birthdate"] = birthdate
	data["avatar"] = authData["figureurl_qq_2"].(string)
	data["city"] = authData["city"].(string)
	data["province"] = authData["province"].(string)
	ret, err := AccountService.Login(formData["openid"].(string), "qq_openid", formData["access_token"].(string), data)
	if err != nil {
		logger.Sugar.Error(err)
		util.RespFail(writer, "系统繁忙")
		return
	}

	util.RespOk(writer, ret, "ok")
}
