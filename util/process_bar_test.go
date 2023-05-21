package util

import (
	"testing"
	"time"
)

// TestDrawProgressBar 需要在console执行才能看到效果
func TestDrawProgressBar(t *testing.T) {
	for i := 0; i <= 100; i++ {
		DrawProgressBar("进度条:", i, 100)
		time.Sleep(time.Millisecond * 100)
	}
}
