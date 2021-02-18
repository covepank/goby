package idx

import (
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sanbsy/goby/bufferpool"
)

var (
	seed       = int64(1589858037028322000)
	merchantID uint64
	idCounter  uint64
)

type ID uint64

func init() {
	rand.Seed(time.Now().Unix())
	merchantID = (uint64(os.Getpid()) & 0x3FF) << 10
	idCounter = rand.Uint64()
}

// NewID snowflake生成唯一分布式ID
// 1位符号位0+41位时间戳+10位机器ID+10位计数器
func NewID() ID {
	timestamp := uint64(time.Now().UnixNano()-seed) &^ (uint64(0x3FFFFF)) >> 1
	counter := (atomic.AddUint64(&idCounter, 1)) & 0xFFF
	return ID(timestamp | merchantID | counter)
}

// String 以36进制格式化为字符串
func (id ID) String() string {
	buf := bufferpool.Get()
	defer buf.Free()
	_, _ = buf.Write(strconv.AppendUint(buf.Bytes(), uint64(id), 36))
	return buf.String()
}
