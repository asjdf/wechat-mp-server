package menu

import "unicode/utf8"

const (
	maxRowLength = 12
)

type ItemInfo struct {
	Hide bool   // 是否显示
	Name string // 名称 之前想把Name和Url 直接拿Raw 但发现直接用Raw会导致无法合理处理换行等问题 这里涉及到一个问题：换行应该是自动实现还是手动实现？
	Url  string // 链接
}

type Item interface {
	Info(staffId string) ItemInfo
}

func Temp(staffId string, item ...Item) string {
	var output string
	rowLength := 0
	for _, v := range item {
		info := v.Info(staffId)
		if !info.Hide {
			if rowLength+utf8.RuneCountInString(info.Name) > maxRowLength && len(output) > 0 && output[len(output)-1] != '\n' { // 超出最大字数上限且之前没有换行 则换行
				output += "\n"
				rowLength = 0
			}
			if rowLength > 0 { // 之前有内容
				output += "  |  "
				rowLength++
			}
			if info.Url != "" {
				output += "<a href=\"" + info.Url + "\">" + info.Name + "</a>"
			} else {
				output += info.Name
			}
			rowLength += utf8.RuneCountInString(info.Name)
			if info.Name[len(info.Name)-1] == '\n' { // 如果主动进行了换行 则清零行长度
				rowLength = 0
			}
		}
	}
	return output
}
