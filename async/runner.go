package async

import (
	"context"
	"sync"

	"github.com/sanbsy/gopkg/log"
)

// Runner 异步运行器
type Runner struct {
	wg *sync.WaitGroup
}

// NewRunner 初始化 runner
func NewRunner() *Runner {
	return &Runner{
		wg: &sync.WaitGroup{},
	}
}

// Run 异步运行一个任务
func (runner *Runner) Run(fn func()) {
	runner.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WarnF("recover a panic:%v", err)
			}
			runner.wg.Done()
		}()
		fn()
	}()
}

// RunCtx 异步运行一个任务，并接收一个Context
func (runner *Runner) RunCtx(ctx context.Context, fn func()) {
	runner.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WarnCtxF(ctx, "recover a panic:%v", err)
			}
			runner.wg.Done()
		}()
		fn()
	}()
}

// Await 等待所有任务执行完成
func (runner *Runner) Await() {
	runner.wg.Wait()
}

// AwaitWithContext 通过 context 控制等待
func (runner *Runner) AwaitWithContext(ctx context.Context) error {
	done := make(chan struct{})
	// 异步等待
	go func() {
		runner.Await()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	return nil
}
