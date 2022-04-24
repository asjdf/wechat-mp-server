package common

import "wechat-mp-server/module/quickNav/menu"

type plain struct {
	str string
}

func (p plain) Info(_ string) menu.ItemInfo {
	return menu.ItemInfo{
		Name: p.str,
	}
}

func Plain(str string) menu.Item {
	return plain{str}
}
