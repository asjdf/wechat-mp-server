package quickNav

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"wechat-mp-server/hub"
	"wechat-mp-server/module/quickNav/common"
	"wechat-mp-server/module/quickNav/menu"
)

func menuHandler(msg *hub.Message) {
	template := []menu.Item{
		common.Divider(), common.Br(),
		common.Divider(true), common.Br(),
		common.Plain("一段话"),
		common.Divider(),
		common.Plain("两段话\n"),
		common.Url("按钮1", "https://homeboyc.cn"),
		common.Url("按钮2", "https://homeboyc.cn"),
		common.Url("按钮3", "https://homeboyc.cn"),
		common.When(
			"2021-01-01 00:00:00",
			"2021-01-02 00:00:00",
			common.Plain("这行看不见\n"),
		),
		common.When(
			"2021-01-01 00:00:00",
			"2050-01-02 00:00:00",
			common.Plain("这行看得见\n"),
		),
	}
	msg.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(menu.Temp(msg.OpenID, template...))}
}
