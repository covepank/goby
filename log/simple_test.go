package log

import (
	"context"
	"testing"
)

func TestDebugCtx(t *testing.T) {
	InfoCtx(context.Background(),"hello world")
}
