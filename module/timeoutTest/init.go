package timeoutTest

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"strings"
	"sync"
	"time"
	"wechat-mp-server/hub"
)

func init() {
	Instance = &Mod{}
	hub.RegisterModule(Instance)
}

var Instance *Mod

type Mod struct {
}

func (m *Mod) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		ID:       hub.NewModuleID("atom", "timeoutTest"),
		Instance: Instance,
	}
}

func (m *Mod) Init() {

}

func (m *Mod) PostInit() {

}

func (m *Mod) Serve(s *hub.Server) {
	s.MsgEngine.MsgText("^超时回复测试$", 1, func(m *hub.Message) {
		time.Sleep(time.Second * 10)
		m.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("超时回复测试")}
	})
	s.MsgEngine.MsgText("^超时回复测试2$", 1, func(m *hub.Message) {
		time.Sleep(time.Second * 10)
		m.Reply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText(strings.Repeat("超时回复测试 超时回复测试\n", 60)),
		}
	})
}

func (m *Mod) Start(_ *hub.Server) {

}

func (m *Mod) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	wg.Done()
}
