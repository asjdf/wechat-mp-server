package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"time"
	"wechat-mp-server/utils"
)

var ginLogger = logrus.WithField("hub", "gin")
var wechatMsgLogger = logrus.WithField("hub", "wechatMsg")

// ginRequestLog gin 请求纪录( 请求状态码 处理时间 请求方法 IP 路由 )
func ginRequestLog(c *gin.Context) {
	// 开始时间
	startTime := time.Now()

	// 处理请求
	c.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := c.Request.Method

	// 请求路由
	reqUri := c.Request.RequestURI

	// 状态码
	statusCode := c.Writer.Status()

	// 请求IP
	clientIP := c.ClientIP()

	// 日志格式
	ginLogger.Infof(
		"| %3d | %13v | %15s | %s | %s |",
		statusCode, latencyTime, clientIP, reqMethod, reqUri,
	)
}

func wechatMsgLog(m *Message) {
	// 开始时间
	startTime := time.Now()

	// 建立sentry
	sentryCtx := context.Background()
	span := sentry.StartSpan(sentryCtx, "passiveReply",
		sentry.TransactionName(fmt.Sprintf("type: %s, key: %s", m.MsgType, m.Pattern)),
	)
	m.Span = span.Context()

	// 处理请求
	m.Next()

	// sentry结束
	span.Finish()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(m.Reply)
	reply := buffer.String()

	var msgType, key, id string
	if m.MsgType == message.MsgTypeEvent {
		msgType = string(m.Event)
		key = m.EventKey
	} else {
		msgType = string(m.MsgType)
		key = m.Content
	}
	if m.UnionID != "" {
		id = m.UnionID
	} else {
		id = string(m.FromUserName)
	}
	wechatMsgLogger.Infof(
		"| %10v | %v | %v | %s | %s |",
		latencyTime, msgType, key, id, reply)
}

func wechatLongMsgHandle(m *Message) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		m.Next()
		cancel()
	}()

	select {
	case <-time.After(time.Millisecond * 4500): // 微信那边超时是5s 文档没写具体怎么算的5s 所以4.5s保险
		m.Reply = nil
		go func(m *Message) {
			select {
			case <-ctx.Done():
				// 超时后调了客服消息，所以得分类处理各类客服消息
				switch m.Reply.MsgType {
				case message.MsgTypeText:
					Content := m.Reply.MsgData.(*message.Text).Content
					reply := utils.CutMsg(string(Content), 2000, 1750)
					if len(reply) > 0 {
						for i := 0; i < len(reply); i++ {
							go func(msg string, order int) {
								time.Sleep(time.Duration(order) * 200 * time.Millisecond)
								textMessage := message.NewCustomerTextMessage(m.GetOpenID(), msg)
								sendCustomerMsg(textMessage)
							}(reply[i], i)
						}
					}
				case message.MsgTypeImage:
					imgMessage := message.NewCustomerImgMessage(m.GetOpenID(), m.Reply.MsgData.(*message.Image).Image.MediaID)
					sendCustomerMsg(imgMessage)
				case message.MsgTypeVoice:
					voiceMessage := message.NewCustomerVoiceMessage(m.GetOpenID(), m.Reply.MsgData.(*message.Voice).Voice.MediaID)
					sendCustomerMsg(voiceMessage)
				case message.MsgTypeMiniprogrampage:
					miniProgramPageMessage := message.NewCustomerMiniprogrampageMessage(m.GetOpenID(),
						m.Reply.MsgData.(*message.MediaMiniprogrampage).Title,
						m.Reply.MsgData.(*message.MediaMiniprogrampage).AppID,
						m.Reply.MsgData.(*message.MediaMiniprogrampage).Pagepath,
						m.Reply.MsgData.(*message.MediaMiniprogrampage).ThumbMediaID)
					sendCustomerMsg(miniProgramPageMessage)
				}
			}
		}(m)
	case <-ctx.Done():
		// 没有超时的情况 只需要处理长文本消息即可
		if m.Reply.MsgType == message.MsgTypeText {
			Content := m.Reply.MsgData.(*message.Text).Content
			reply := utils.CutMsg(string(Content), 2000, 1750)
			if len(reply) > 1 {
				m.Reply.MsgData = message.NewText(reply[0])
				for i := 1; i < len(reply); i++ {
					go func(msg string, order int) {
						time.Sleep(time.Duration(order) * 200 * time.Millisecond)
						m := message.NewCustomerTextMessage(m.GetOpenID(), msg)
						sendCustomerMsg(m)
					}(reply[i], i)
				}
			}
		}
	}
}

func sendCustomerMsg(customerMessage *message.CustomerMessage) {
	err := Instance.WechatEngine.GetCustomerMessageManager().Send(customerMessage)
	if err != nil {
		logger.Errorf("send custom text msg failed: %v", err)
	}
}
