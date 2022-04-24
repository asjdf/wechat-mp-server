package common

import "wechat-mp-server/module/quickNav/menu"

type br struct {
}

func (br) Info(_ string) menu.ItemInfo {
	return menu.ItemInfo{
		Name: "\n",
	}
}

func Br() menu.Item {
	return br{}
}
