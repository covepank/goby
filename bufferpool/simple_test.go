package bufferpool_test

import (
	"testing"

	"github.com/sanbsy/goby/bufferpool"
	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {

	_cases := []struct {
		name   string
		action func(buf *bufferpool.Buffer)
		want   string
	}{
		{
			"write byte",
			func(buf *bufferpool.Buffer) {
				buf.WriteString("hello world")
			},
			"hello world",
		},
		{
			"write uint",
			func(buf *bufferpool.Buffer) {
				buf.WriteUint(1023)
			},
			"1023",
		},
		{
			"base64",
			func(buf *bufferpool.Buffer) {
				buf.WriteString("hello world")
				bb := buf.Base64()
				buf.Reset()
				buf.WriteString(bb)
			},
			"aGVsbG8gd29ybGQ=",
		},
	}

	for _, cc := range _cases {
		t.Run(cc.name, func(t *testing.T) {
			buf := bufferpool.Get()
			defer buf.Free()
			cc.action(buf)
			assert.Equal(t, cc.want, buf.String())
		})
	}
}
