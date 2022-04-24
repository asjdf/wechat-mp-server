package templateMessage

import (
	"encoding/json"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"strings"
)

const maxSenderNumber = 50          // 最大协程数
const globalMaxRetryTime = int64(3) // 全局最大重试次数 优先级低于用户自己配置的

type TemplateMessage struct {
	Message      *message.TemplateMessage
	Resend       bool  // 发送失败后是否重新发送
	RetriedTime  int64 // 记录发送失败的重试次数
	MaxRetryTime int64 // 此消息的最大重试次数 -1为一直重试
}

// PushMessage 向消息队列中塞待发送消息
func (m *Module) PushMessage(superMessage *TemplateMessage) {
	go func() {
		m.MessageQueue <- superMessage
	}()
}

// registerMessageSender 创建新的协程来发送消息
func (m *Module) registerMessageSender(msgChannel <-chan *TemplateMessage) {
	for msg := range msgChannel {
		m.senderLimit <- struct{}{} // 等待有多余的协程量
		m.MessageSenderWaitGroup.Add(1)
		go func(t *TemplateMessage) {
			defer m.MessageSenderWaitGroup.Done()
			m.sendMessage(t)
			<-m.senderLimit // 释放
		}(msg)
	}
}

// sendMessage 发送模板消息
func (m *Module) sendMessage(templateMessage *TemplateMessage) {
	_, err := Template.Send(templateMessage.Message)

	msgMarshal, _ := json.Marshal(templateMessage)
	if err != nil {
		if strings.Contains(err.Error(), "43004") {
			logger.Warn("Send templateMsg failed, errcode 43004")
			return
		}
		logger.Warnf("Send templateMsg failed, msg: %v , errMsg: %v", string(msgMarshal), err.Error())
		if templateMessage.Resend {
			if templateMessage.MaxRetryTime == -1 {
				m.MessageQueue <- templateMessage
			} else if templateMessage.MaxRetryTime != 0 && templateMessage.RetriedTime < templateMessage.MaxRetryTime {
				templateMessage.RetriedTime += 1
				m.MessageQueue <- templateMessage
			} else if templateMessage.MaxRetryTime == 0 && templateMessage.RetriedTime < globalMaxRetryTime {
				templateMessage.RetriedTime += 1
				m.MessageQueue <- templateMessage
			} else {
				logger.Warnf("Send template message failed, msg: %v , stop retry.", string(msgMarshal)) // 这部分log需要优化
			}
		}
		return
	}
	logger.Info("Send templateMsg sucess:" + string(msgMarshal))
}
