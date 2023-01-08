package common

import (
	"time"
	"wechat-mp-server/module/quickNav/menu"
)

type when struct {
	start time.Time
	end   time.Time
	item  menu.Item
}

func (w when) Info(staffId string) (info menu.ItemInfo) {
	if !(time.Now().After(w.start) && time.Now().Before(w.end)) {
		info.Hide = true
		return
	}
	info = w.item.Info(staffId)
	return
}

// When 当在时间区间内才显示
// warning: 为方便使用的关系 不方便抛出错误 请自行保证参数正确 隐藏的优先级关系：When > Item 自带的隐藏
func When(start string, end string, item menu.Item) menu.Item {
	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)
	return when{start: startTime, end: endTime, item: item}
}
