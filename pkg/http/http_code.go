package http

const (
	HTTP_CODE_OK                  = 0   // 成功
	HTTP_CODE_ERROR               = 500 // 服务器错误
	HTTP_CODE_ACCOUNT_TOKEN_ERROR = 401 // 重新登陆
	HTTP_CODE_SMS_ERROR           = 405 // 短信验证码错误
)
