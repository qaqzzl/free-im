package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type H struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

//
func RespFail(c *gin.Context, msg string) {
	Resp(c, 500, nil, msg)
}
func RespOk(c *gin.Context, data interface{}, msg string) {
	Resp(c, 0, data, msg)
}
func Resp(c *gin.Context, code int, data interface{}, msg string) {
	//定义一个结构体
	h := H{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, h)
}
