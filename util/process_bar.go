package util

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const PROGRESS_WIDTH = 100

// DrawProgressBar 显示进度条
func DrawProgressBar(prefix string, val, max int) {
	proportion := float32(val) / float32(max)
	pos := int(proportion * PROGRESS_WIDTH)
	s := fmt.Sprintf("%s [%s%*s] %6.2f%% \t[%d/%d]",
		prefix, strings.Repeat("■", pos), PROGRESS_WIDTH-pos, "", proportion*100, val, max)
	fmt.Print(color.CyanString("\r" + s))
	if proportion >= 1 {
		fmt.Print("\n")
	}
}
