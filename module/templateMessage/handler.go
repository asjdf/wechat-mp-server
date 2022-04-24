package templateMessage

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"net/http"
	"time"
)

// handleOldWechatTemplateMessage 兼容旧接口的模板消息处理
func handleOldWechatTemplateMessage(c *gin.Context) {
	data := &message.TemplateMessage{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusOK,
			buildGinResp("20", err.Error(), nil),
		)
		return
	}
	if data.TemplateID == "" {
		c.JSON(http.StatusOK,
			buildGinResp("21", "template_id can't be empty", nil),
		)
		return
	}
	if data.ToUser == "" {
		c.JSON(http.StatusOK,
			buildGinResp("22", "touser can't be empty", nil),
		)
		return
	}
	superMsg := &TemplateMessage{
		Message:      data,
		Resend:       true,
		RetriedTime:  0,
		MaxRetryTime: 0,
	}
	instance.PushMessage(superMsg)
	c.JSON(http.StatusOK, buildGinResp("0", "ok", ""))
}

func buildGinResp(errcode, errmsg string, data interface{}) gin.H {
	return gin.H{
		"error":     errcode,
		"msg":       errmsg,
		"isCache":   false,
		"timestamp": time.Now().Unix(),
		"data":      data,
	}
}
