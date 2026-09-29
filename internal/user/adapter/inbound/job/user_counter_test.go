package job

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
)

type fakeCounter struct {
	calls atomic.Int32
	err   error
}

func (f *fakeCounter) Count(context.Context) (int64, error) {
	f.calls.Add(1)
	return 42, f.err
}

type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func run(t *testing.T, counter *fakeCounter) string {
	t.Helper()
	var out safeBuffer
	logger.Init(&out, "info")
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		RunUserCounter(ctx, counter, 10*time.Millisecond)
		close(done)
	}()

	time.Sleep(55 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("counter did not stop after cancel")
	}
	return out.String()
}

func TestRunUserCounterLogsTotal(t *testing.T) {
	counter := &fakeCounter{}
	logs := run(t, counter)
	if counter.calls.Load() < 3 {
		t.Errorf("expected several ticks, got %d", counter.calls.Load())
	}
	if !strings.Contains(logs, "total users") || !strings.Contains(logs, `"count":42`) {
		t.Errorf("missing count log:\n%s", logs)
	}
	if !strings.Contains(logs, "user counter stopped") {
		t.Errorf("missing stop log:\n%s", logs)
	}
}

func TestRunUserCounterKeepsGoingOnError(t *testing.T) {
	counter := &fakeCounter{err: errors.New("mongo down")}
	logs := run(t, counter)
	if counter.calls.Load() < 3 {
		t.Errorf("counter stopped after an error, calls=%d", counter.calls.Load())
	}
	if !strings.Contains(logs, "count users failed") {
		t.Errorf("missing error log:\n%s", logs)
	}
}
