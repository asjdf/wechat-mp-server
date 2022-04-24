package utils

import (
	"regexp"
)

//const maxMsgLen = 2048

func CutMsg(msg string, maxLen, minLen int) (reply []string) {
	rest := msg
	for len(rest) > 2048 {
		var msg string
		msg, rest = CutByLine(rest, maxLen, minLen)
		if msg != "" {
			reply = append(reply, msg)
			continue
		}
		msg, rest = CutByParagraph(rest, maxLen, minLen)
		if msg != "" {
			reply = append(reply, msg)
			continue
		}
		msg, rest = CutBySentence(rest, maxLen, minLen)
		if msg != "" {
			reply = append(reply, msg)
			continue
		}
	}
	reply = append(reply, rest)
	return
}

// StrSplit 分割字符串
// text为需要处理的文本
// reg为分割点的正则对象
// maxLen为切割点的最大位置 如果超过最大位置还未找到切割点 则返回的c为空
// minLen为切割点的最小位置 切割点位置小于最小的切割点位置 则返回的c为空
// firstPos是切割时起始点的位置 false表示从切割点开始的位置开始切除 true则为切割点结束的地方开始切除
// secondPos是切割时起始点的位置 false表示从切割点开始的位置开始切除 true则为切割点结束的地方开始切除
func StrSplit(text string, reg *regexp.Regexp, maxLen, minLen int, firstPos, secondPos bool) (c, rest string) {
	idx := reg.FindAllStringIndex(text, -1)
	if len(idx) != 0 {
		last := getClosestIndex(idx, maxLen)
		// 切的时候确保距离最大字数限制较小
		if last[1] > minLen {
			c = text[0:last[b2i(firstPos)]]
			rest = text[last[b2i(secondPos)]:]
		} else {
			rest = text
		}
	} else {
		rest = text
	}
	return
}

// CutByLine 根据分割线裁切 如果没法切 则c返回为空
func CutByLine(msg string, maxLen int, minLen int) (c, rest string) {
	cutLine := regexp.MustCompile("[─]+\n")
	return StrSplit(msg, cutLine, maxLen, minLen, true, true)
}

// CutByParagraph 根据段落裁切
func CutByParagraph(msg string, maxLen int, minLen int) (c, rest string) {
	cutParagraph := regexp.MustCompile(`\n`)
	return StrSplit(msg, cutParagraph, maxLen, minLen, false, true)
}

// CutBySentence 根据语句裁切
func CutBySentence(msg string, maxLen int, minLen int) (c, rest string) {
	cutSentence := regexp.MustCompile(`[.。]`) // 查找句末
	return StrSplit(msg, cutSentence, maxLen, minLen, true, true)
}
func getClosestIndex(indexes [][]int, maxPos int) (last []int) {
	last = indexes[0]
	// 找距离字数限制最近的一条可用的分割位置
	for i := 1; i < len(indexes); i++ {
		if indexes[i][1] > maxPos { // 该条已经超出字数限制
			break
		}
		last = indexes[i] // 未超出 继续找
	}
	return
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
