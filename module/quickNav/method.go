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
		common.Url("按钮1", "http://www.baidu.com"),
		common.Url("按钮2", "http://www.baidu.com"),
		common.Url("按钮3", "http://www.baidu.com"),
	}
	msg.Reply = &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(menu.Temp(msg.OpenID, template...))}
}
