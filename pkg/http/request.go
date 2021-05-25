package http

import (
	"free-im/pkg/logger"
	"github.com/gin-gonic/gin"
)

func ReqBin(c *gin.Context, data interface{}) error {
	if err := c.Bind(&data); err != nil {
		logger.Sugar.Info(err.Error())
		RespFail(c, "数据格式不合法")
		return err
	}
	return nil
}

func ShouldBindUri(c *gin.Context, data interface{}) error {
	if err := c.ShouldBindUri(data); err != nil {
		RespFail(c, err.Error())
		return err
	}
	return nil
}

func GetUid(c *gin.Context) int64 {
	if value, exists := c.Get("authorized_member_id"); exists {
		return value.(int64)
	}
	return 0
}

func GetDeviceId(c *gin.Context) string {
	value := c.Request.Header.Get("device_id")
	return value
}

func GetClientType(c *gin.Context) string {
	value := c.Request.Header.Get("client_type")
	return value
}
