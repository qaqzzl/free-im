package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type H struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Rows interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}
//
func RespFail(w http.ResponseWriter,msg string){
	Resp(w,500,nil,msg)
}
func RespOk(w http.ResponseWriter,data interface{},msg string){
	Resp(w,0,data,msg)
}
func Resp(w http.ResponseWriter,code int,data interface{},msg string)  {

	w.Header().Set("Content-Type","application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	//输出
	//定义一个结构体
	h := H{
		Code:code,
		Msg:msg,
		Data:data,
	}
	//将结构体转化成JSOn字符串
	ret,err := json.Marshal(h)
	if err!=nil{
		log.Println(err.Error())
	}
	//输出
	w.Write(ret)
}