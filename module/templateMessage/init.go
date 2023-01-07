// 发送模板消息api
// 消息发送相关的配置请在method中修改
// todo: 确认message缓冲池大小

package templateMessage

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"sync"
	"wechat-mp-server/hub"
	"wechat-mp-server/utils"
)

func init() {
	instance = &Mod{}
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var logger *logrus.Entry

var Template *message.Template

var instance *Mod

type Mod struct {
	MessageQueue           chan *TemplateMessage
	MessageSenderWaitGroup sync.WaitGroup
	senderLimit            chan struct{} // 最大协程限制
}

func (m *Mod) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		ID:       hub.NewModuleID("atom", "templateMessage"),
		Instance: instance,
	}
}

func (m *Mod) Init() {
	logger.Info("Init template message sender...")

	m.MessageQueue = make(chan *TemplateMessage, maxSenderNumber)
	m.senderLimit = make(chan struct{}, maxSenderNumber)
}

func (m *Mod) PostInit() {

}

func (m *Mod) Serve(s *hub.Server) {
	Template = message.NewTemplate(s.WechatEngine.GetContext())
	go m.registerMessageSender(m.MessageQueue)
}

func (m *Mod) Start(_ *hub.Server) {
	// example:
	//m.PushMessage(&TemplateMessage{
	//	Message: &message.TemplateMessage{
	//		ToUser:     "unionId",
	//		TemplateID: "_HtuD7TrFKxquwizJwICXv4sWg5AeZBvaHBIRvYKeKk",
	//		URL:        "",
	//		Color:      "",
	//		Data: map[string]*message.TemplateDataItem{
	//			"keyword1": &message.TemplateDataItem{
	//				Value: "测试消息",
	//				Color: "",
	//			},
	//			"keyword2": &message.TemplateDataItem{
	//				Value: "中午问候",
	//				Color: "",
	//			},
	//		},
	//		MiniProgram: struct {
	//			AppID    string `json:"appid"`
	//			PagePath string `json:"pagepath"`
	//		}{},
	//	},
	//	Resend:       true,
	//	RetriedTime:  0,
	//	MaxRetryTime: 0,
	//})
}

func (m *Mod) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	close(m.senderLimit)
	close(m.MessageQueue) // 关闭此channel后sender携程会自己关闭
	m.MessageSenderWaitGroup.Wait()
	wg.Done()
}
