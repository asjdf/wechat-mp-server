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
	fmt.Println(menu.Temp("", template...))
}
