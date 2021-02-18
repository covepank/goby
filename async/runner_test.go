package async

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRunner_Run(t *testing.T) {
	runner := NewRunner()

	var v1 int
	runner.Run(func() {
		time.Sleep(1 * time.Second)
		v1 = 1
	})

	var v2 int
	runner.Run(func() {
		time.Sleep(5 * time.Second)
		v2 = 2
	})

	var v3 int

	runner.Run(func() {
		time.Sleep(10 * time.Second)
		v3 = 3
	})

	runner.Await()

	assert.Equal(t, 1, v1)
	assert.Equal(t, 2, v2)
	assert.Equal(t, 3, v3)
}

func TestRunner_AwaitWithContext(t *testing.T) {
	runner := NewRunner()

	var v1 int
	runner.Run(func() {
		time.Sleep(1 * time.Second)
		v1 = 1
	})

	var v2 int
	runner.Run(func() {
		time.Sleep(5 * time.Second)
		v2 = 2
	})

	var v3 int

	runner.Run(func() {
		time.Sleep(10 * time.Second)
		v3 = 3
	})

	ctx, _ := context.WithTimeout(context.Background(), 6*time.Second)

	err := runner.AwaitWithContext(ctx)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 1, v1)
	assert.Equal(t, 2, v2)
	assert.Equal(t, 0, v3)
}
