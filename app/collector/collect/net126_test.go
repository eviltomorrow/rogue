package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMetadataFromNet126(t *testing.T) {
	_assert := assert.New(t)
	var codes = []string{
		"sh600030", "sz300002", "sz000001", "sh688009",
	}
	data, err := FetchMetadataFromNet126(codes)
	_assert.Nil(err)
	t.Logf("data: %s\r\n", data)
}
