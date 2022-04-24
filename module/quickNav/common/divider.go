package common

import "wechat-mp-server/module/quickNav/menu"

type divider struct {
	Vertical bool
}

func (d divider) Info(_ string) menu.ItemInfo {
	info := menu.ItemInfo{}
	if d.Vertical {
		info.Name = "  |  "
	} else {
		info.Name = "─────────────"
	}
	return info
}

func Divider(vertical ...bool) menu.Item {
	if len(vertical) > 0 && vertical[0] {
		return divider{
			Vertical: true,
		}
	}
	return divider{}
}
