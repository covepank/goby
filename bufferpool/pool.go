package bufferpool

import (
	"sync"
)

const defaultCap = 256

// Pool buffer pool
type Pool struct {
	pl *sync.Pool
}

// NewPool create a pool
func NewPool() Pool {
	return Pool{pl: &sync.Pool{
		New: func() interface{} {
			return &Buffer{
				buf: make([]byte, 0, defaultCap),
			}
		},
	}}
}

// Get get a buffer from pool
func (p Pool) Get() *Buffer {
	buf := p.pl.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	p.pl.Put(buf)
}
