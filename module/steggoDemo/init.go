package steggoDemo

import (
	"fmt"
	"github.com/asjdf/steggo"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"regexp"
	"sync"
	"wechat-mp-server/config"
	"wechat-mp-server/hub"
)

type Mod struct {
}

func (m *Mod) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		ID:       hub.NewModuleID("atom", "steggoDemo"),
		Instance: m,
	}
}

func (m *Mod) Init() {

}

func (m *Mod) PostInit() {

}

func (m *Mod) Serve(s *hub.Server) {
	s.MsgEngine.MsgText("^告诉我一个秘密$", 1, func(msg *hub.Message) {
		tracker, _ := steggo.Encode([]byte(msg.GetOpenID()))

		msg.Reply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(fmt.Sprintf("这%s个秘密你可不能和别人说哦！这个后端的版本是%s", tracker, config.Version)),
		}
	})
	s.MsgEngine.MsgText("^追踪.*$", 1, func(msg *hub.Message) {
		embedTracker := regexp.MustCompile(`[^\x{200c}\x{200d}\x{2060}\x{2062}\x{2063}\x{2064}]+`).
			ReplaceAllString(msg.Content, "")
		if len(embedTracker) == 0 {
			msg.Reply = &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText("无法追踪该消息"),
			}
			return
		}
		from, err := steggo.Decode(embedTracker)
		if err != nil {
			msg.Reply = &message.Reply{
				MsgType: message.MsgTypeText,
				MsgData: message.NewText("追踪器受损，无法追踪该消息"),
			}
			return
		}
		msg.Reply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(fmt.Sprintf("这个消息的原接收者为%s", from)),
		}
	})
}

func (m *Mod) Start(_ *hub.Server) {

}

func (m *Mod) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	wg.Done()
}
