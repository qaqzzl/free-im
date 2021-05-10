package api

import (
	"encoding/json"
	"fmt"
	"free-im/config"
	"free-im/internal/http_app/model"
	"free-im/internal/http_app/service"
	http2 "free-im/pkg/http"
	"free-im/pkg/logger"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

var AccountService *service.AccountService
var CommonService *service.CommonService

// 手机号登录 / 注册
func PhoneLogin(c *gin.Context) {
	// 初始化请求变量结构
	// Parse JSON
	var req struct {
		Phone      string `json:"phone" binding:"required"`
		VerifyCode string `json:"verify_code" binding:"required"`
	}
	if http2.ReqBin(c, &req) != nil {
		return
	}
	// 验证短信验证码是否正确
	if verify, err := CommonService.IsPhoneVerifyCode(req.Phone, req.VerifyCode, "login"); err != nil {
		fmt.Println(err)
		http2.RespFail(c, err.Error())
		return
	} else if verify == false {
		http2.RespFail(c, "短信验证码错误")
		return
	}
	ret, err := AccountService.
		Login(model.UserAuths{Identifier: req.Phone,
			IdentityType: "phone"}, model.UserMember{}, http2.GetDeviceId(c), http2.GetClientType(c))
	if err != nil {
		logger.Logger.Info(err.Error())
		http2.RespFail(c, "系统繁忙")
		return
	}

	http2.RespOk(c, ret, "")
}

// QQ登陆
func QQLogin(c *gin.Context) {
	// 初始化请求变量结构
	var req struct {
		Openid      string `json:"openid"`
		AccessToken string `json:"access_token"`
	}
	if http2.ReqBin(c, &req) != nil {
		return
	}
	resp, err := http.Get("https://graph.qq.com/user/get_user_info?access_token=" + req.AccessToken +
		"&oauth_consumer_key=" + config.CommonConf.QQAuthAppID +
		"&openid=" + req.Openid)
	if err != nil {
		logger.Sugar.Error(err)
		http2.RespFail(c, "系统繁忙")
		return
	}
	defer resp.Body.Close()
	authBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		logger.Sugar.Error(err)
		http2.RespFail(c, "系统繁忙")
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
	ret, err := AccountService.
		Login(model.UserAuths{Identifier: req.Openid, IdentityType: "qq_openid",
			Credential: req.AccessToken}, UserMember, http2.GetDeviceId(c), http2.GetClientType(c))
	if err != nil {
		logger.Sugar.Error(err)
		http2.RespFail(c, "系统繁忙")
		return
	}

	http2.RespOk(c, ret, "ok")
}
