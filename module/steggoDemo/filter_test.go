package steggoDemo

import (
	"fmt"
	"regexp"
	"testing"
)

func TestFilter(t *testing.T) {
	embedTracker := regexp.MustCompile(`[^\x{200c}\x{200d}\x{2060}\x{2062}\x{2063}\x{2064}]+`).
		ReplaceAllString("这⁤⁤⁤⁤⁤⁤⁤⁤⁤⁤个秘密你可不能和别人说哦！这个后端的版本是9e96b7c", "")
	fmt.Print(embedTracker)

}
