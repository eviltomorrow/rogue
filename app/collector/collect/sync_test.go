package collect

import (
	"testing"

	"github.com/eviltomorrow/rogue/lib/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestSyncDataQuick(t *testing.T) {
	_assert := assert.New(t)

	mongodb.DSN = "mongodb://127.0.0.1:27017"
	err := mongodb.Build()
	_assert.Nil(err)

	affected, err := SyncDataQuick("net126")
	_assert.Nil(err)
	t.Logf("affected: %v", affected)
}
