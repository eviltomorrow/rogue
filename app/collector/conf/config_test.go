package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAndLoad(t *testing.T) {
	_assert := assert.New(t)
	err := Global.FindAndLoad("../global.conf", nil)
	_assert.NoError(err)
}
