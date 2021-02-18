package store

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Instance() (*BoltStorage, error) {
	return NewBolt(path.Join(os.TempDir(), "data.db"))
}

func TestBoltBucket_Set(t *testing.T) {
	ass := assert.New(t)
	sto, err := Instance()
	ass.Nil(err)

	ass.Nil(sto.Set([]byte("test_key"), []byte("test_value")))
	ass.Nil(sto.Set([]byte("test_key"), []byte("test_value_2")))

	value, err := sto.Get([]byte("test_key"))
	ass.Nil(err)

	ass.Equal([]byte("test_value_2"), value)
}
