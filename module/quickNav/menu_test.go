package quickNav

import (
	"fmt"
	"testing"
	"wechat-mp-server/module/quickNav/common"
	"wechat-mp-server/module/quickNav/menu"
)

func TestGenMenu(t *testing.T) {
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
	fmt.Println(menu.Temp("", template...))
}
