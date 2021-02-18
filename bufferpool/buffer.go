package bufferpool

import (
	"encoding/base64"
	"strconv"
	"time"
	"unsafe"
)

// Buffer  buffer
type Buffer struct {
	buf  []byte
	pool Pool
}

// Len 获取 len
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Cap 获取cap
func (b *Buffer) Cap() int {
	return cap(b.buf)
}

// Bytes 以 []byte 格式输出
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// String 以 string 格式输出
func (b *Buffer) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}

// Base64 base64编码后输出
func (b *Buffer) Base64() string {
	return base64.StdEncoding.EncodeToString(b.buf)
}

// Reset 清除内容
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
}

// Writer 写入数据
func (b *Buffer) Write(bs []byte) (int, error) {
	b.buf = append(b.buf, bs...)
	return len(bs), nil
}

// WriteByte 写入 一个 byte
func (b *Buffer) WriteByte(v byte) {
	b.buf = append(b.buf, v)
}

// Free 释放 buffer
func (b *Buffer) Free() {
	b.pool.pl.Put(b)
}

// WriteBytes 写入 []byte
func (b *Buffer) WriteBytes(bs []byte) {
	b.buf = append(b.buf, bs...)
}

// WriteString 写入 string
func (b *Buffer) WriteString(s string) {
	b.buf = append(b.buf, s...)
}

// WriteTime 写入时间
func (b *Buffer) WriteTime(t time.Time, layout ...string) {
	if len(layout) == 0 {
		b.buf = t.AppendFormat(b.buf, time.RFC3339)
	} else {
		b.buf = t.AppendFormat(b.buf, layout[0])
	}
}

// WriteInt64 写入 int64
func (b *Buffer) WriteInt64(i int64) {
	b.buf = strconv.AppendInt(b.buf, i, 10)
}

// WriteInt 写入 int
func (b *Buffer) WriteInt(i int) {
	b.WriteInt64(int64(i))
}

// WriteUint64 写入 uint64
func (b *Buffer) WriteUint64(i uint64) {
	b.buf = strconv.AppendUint(b.buf, i, 10)
}

// WriteUint 写入 uint
func (b *Buffer) WriteUint(i uint) {
	b.WriteUint64(uint64(i))
}

// WriteBool 写入 bool
func (b *Buffer) WriteBool(v bool) {
	b.buf = strconv.AppendBool(b.buf, v)
}

// WriteFloat 写入 float
func (b *Buffer) WriteFloat(f float64) {
	b.buf = strconv.AppendFloat(b.buf, f, 'f', -1, 64)
}
