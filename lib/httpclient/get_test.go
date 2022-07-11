package httpclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	_assert := assert.New(t)
	data, err := Get("http://www.baidu.com", 10*time.Second, DefaultHeader)
	_assert.Nil(err)
	t.Logf("data: %v\r\n", data)
}
