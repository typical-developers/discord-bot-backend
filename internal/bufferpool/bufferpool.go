package bufferpool

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool *sync.Pool
}

var (
	Buffers BufferPool
)

func init() {
	Buffers = BufferPool{
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (b *BufferPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

func (b *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}
