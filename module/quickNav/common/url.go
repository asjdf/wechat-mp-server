package common

import "wechat-mp-server/module/quickNav/menu"

type url struct {
	link string
	name string
}

func (u url) Info(_ string) menu.ItemInfo {
	return menu.ItemInfo{
		Name: u.name,
		Url:  u.link,
	}
}

func Url(name, link string) menu.Item {
	return url{
		name: name,
		link: link,
	}
}
