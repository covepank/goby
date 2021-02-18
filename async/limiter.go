package async

import "context"

// Limiter 令牌池
type Limiter struct {
	ch chan struct{}
}

// NewLimiter 初始化令牌桶
func NewLimiter(cap int) *Limiter {
	return &Limiter{
		ch: make(chan struct{}, cap),
	}
}

// Acquire 申请令牌
func (l *Limiter) Acquire(ctx context.Context) error {
	if cap(l.ch) == 0 {
		return nil
	}

	select {
	case l.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release 释放令牌
func (l *Limiter) Release() {
	select {
	case <-l.ch:
	default:
	}
}
