package ticker

import (
	"testing"
	"time"
)

func TestNewAlignedTicker(t *testing.T) {
	var ticker = NewAlignedTicker(time.Now(), 10*time.Second, 0, 0)
	for range ticker.Elapsed() {
		t.Log("Hello world")
	}

}
