package wechatPong

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"sync"
	"wechat-mp-server/hub"
	"wechat-mp-server/module/templateMessage"
)

func init() {
	instance = &wechatPong{}
	hub.RegisterModule(instance)
}

var instance *wechatPong

type wechatPong struct {
}

func (m *wechatPong) GetModuleInfo() hub.ModuleInfo {
	return hub.ModuleInfo{
		// module 的 id 全部用 NewModule 函数生成
		// namespace 为开发者，moduleName 为模块名（一般用包名） 方便锤开发者以及定位出错的包
		ID:       hub.NewModuleID("atom", "wechatPong"),
		Instance: instance,
	}
}

func (m *wechatPong) Init() {
	// 初始化过程
	// 在此处可以进行 Module 的初始化配置
	// 如配置读取
}

func (m *wechatPong) PostInit() {
	// 第二次初始化
	// 再次过程中可以进行跨Module的动作
	// 如通用数据库等等
}

func (m *wechatPong) Serve(s *hub.Server) {
	// 注册服务函数部分
	// index 是匹配的优先级，index越大，优先级越高，优先测试该条匹配规则
	// key 就是正则匹配的规则 当匹配上之后就该条路由
	// 可以在 https://regex101.com/ 测试你的正则规则
	s.MsgEngine.Group("^ping$", func(msg *hub.Message) {
		go func() {
			_, _ = message.NewTemplate(s.WechatEngine.GetContext()).Send(&message.TemplateMessage{
				ToUser:     msg.GetOpenID(),
				TemplateID: "_HtuD7TrFKxquwizJwICXv4sWg5AeZBvaHBIRvYKeKk", // 这里模板ID需要根据你的公众号修改
				Data: map[string]*message.TemplateDataItem{
					"keyword1": {
						Value: "ping-pong模块",
					},
					"keyword2": {
						Value: "SDK模板消息测试",
					},
				},
			})
		}()

		templateMessageModule, _ := hub.GetModule(hub.NewModuleID("atom", "templateMessage"))
		templateMessageSender := templateMessageModule.Instance.(*templateMessage.Module)
		templateMessageSender.PushMessage(&templateMessage.TemplateMessage{
			Message: &message.TemplateMessage{
				ToUser:     msg.GetOpenID(),
				TemplateID: "_HtuD7TrFKxquwizJwICXv4sWg5AeZBvaHBIRvYKeKk", // 这里模板ID需要根据你的公众号修改
				Data: map[string]*message.TemplateDataItem{
					"keyword1": {
						Value: "ping-pong模块",
					},
					"keyword2": {
						Value: "模板消息模组测试",
					},
				},
			},
			Resend: true,
		})
		u := s.WechatEngine.GetUser()
		tidList, err := u.UserTidList(msg.GetOpenID())
		tidListStr := "tidList获取失败"
		if err == nil {
			marshal, _ := json.Marshal(tidList)
			tidListStr = string(marshal)
		}
		msg.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("pong " + msg.OpenID + " " + tidListStr + "\n" +
			"Current version: " + hub.Version)}
	}).MsgText("", 1).EventClick("")
	// 由于group已经定义了一个中间件实现查询功能 所以后面注册具体方法的时候并不需要带上查询函数
	// 同时由于路由是由baseKey+key拼接而成，所以也不需要单独设置key了
	// 这样还可以快速实现一个查询函数多个查询关键字，如：
	s.MsgEngine.Group("", func(msg *hub.Message) {
		msg.Reply = &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("ping"),
		}
	}).MsgText("^pong$", 1).MsgText("^poooong$", 1)
}

func (m *wechatPong) Start(_ *hub.Server) {
	// 请在可能出错的地方使用 sentry 接住错误 越早defer越好
	defer sentry.Recover()

	// 此函数会新开携程进行调用
	// ```go
	// 		go exampleModule.Start()
	// ```

	// 可以利用此部分进行后台操作
	// 如http服务器等等
}

func (m *wechatPong) Stop(_ *hub.Server, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
	// 结束部分
	// 一般调用此函数时，程序接收到 os.Interrupt 信号
	// 即将退出
	// 在此处应该释放相应的资源或者对状态进行保存
}
