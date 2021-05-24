package errors_test

import (
	"testing"

	"github.com/sanbsy/gopkg/errors"
)

func BenchmarkNewError(b *testing.B) {
	err := errors.New("hhhh")
	for i := 0; i < b.N; i++ {
		_ = errors.Wrap(err, "test")
	}
}

func BenchmarkNewError2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = errors.New("hhhh")
	}
}

func TestNewError(t *testing.T) {
	t.Log(errors.Wrap(errors.New("test error"), "test"))
}
