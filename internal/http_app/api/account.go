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

func BindPushID(c *gin.Context) {

	var req struct {
		PushID string `uri:"push_id" binding:"required"`
	}
	if err := http2.ShouldBindUri(c, &req); err != nil {
		return
	}
	if err := AccountService.BindPushID(http2.GetUid(c), req.PushID); err != nil {
		http2.RespFail(c, err.Error())
	}
	http2.RespOk(c, nil, "ok")
}

// 手机号验证码 登录 / 注册
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
		http2.Resp(c, http2.HTTP_CODE_SMS_ERROR, nil, "验证码错误")
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

// 手机号密码登陆
func PhonePasswordLogin(c *gin.Context) {
	// 初始化请求变量结构
	// Parse JSON
	var req struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if http2.ReqBin(c, &req) != nil {
		return
	}
	// 验证密码是否正确

	ret, err := AccountService.
		Login(model.UserAuths{Identifier: req.Account,
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
	// 初始化变量结构
	var authData struct {
		Ret             int    `json:"ret"`
		Msg             string `json:"msg"`
		IsLost          int    `json:"is_lost"`
		Nickname        string `json:"nickname"`
		Gender          string `json:"gender"`
		GenderType      int    `json:"gender_type"`
		Province        string `json:"province"`
		City            string `json:"city"`
		Year            string `json:"year"`
		Constellation   string `json:"constellation"`
		Figureurl       string `json:"figureurl"`
		Figureurl1      string `json:"figureurl_1"`
		Figureurl2      string `json:"figureurl_2"`
		FigureurlQq1    string `json:"figureurl_qq_1"`
		FigureurlQq2    string `json:"figureurl_qq_2"`
		FigureurlQq     string `json:"figureurl_qq"`
		FigureurlType   string `json:"figureurl_type"`
		IsYellowVip     string `json:"is_yellow_vip"`
		Vip             string `json:"vip"`
		YellowVipLevel  string `json:"yellow_vip_level"`
		Level           string `json:"level"`
		IsYellowYearVip string `json:"is_yellow_year_vip"`
	}
	// 调用json包的解析，解析请求body
	json.Unmarshal(authBody, &authData)
	if authData.Ret != 0 {
		http2.RespFail(c, authData.Msg)
		return
	}

	// 获取unionid https://graph.qq.com/oauth2.0/me?access_token=ACCESSTOKEN&unionid=1
	UserMember := model.UserMember{}
	UserMember.Nickname = authData.Nickname
	if authData.Gender == "女" {
		UserMember.Gender = "w"
	} else {
		UserMember.Gender = "m"
	}
	times, _ := time.Parse("2006", authData.Year)
	timeUnix := times.Unix()
	if timeUnix > 0 {
		UserMember.Birthdate = int(timeUnix)
	}
	UserMember.Avatar = authData.Figureurl2
	UserMember.City = authData.City
	UserMember.Province = authData.Province
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
