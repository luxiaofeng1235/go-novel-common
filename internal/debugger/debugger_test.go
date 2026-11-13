/*
 * @Descripttion: 调试输出的快照、队列、并发及收尾行为测试
 * @Author: red
 * @Date: 2026-11-13
 * @LastEditors: red
 * @LastEditTime: 2026-11-13
 */
package debugger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// useTestQueue 为串行测试隔离队列，后台默认消费者仍持有原队列。
func useTestQueue(t *testing.T, size int) chan entry {
	t.Helper()
	originalQueue, originalDropped := queue, dropped.Load()
	queue = make(chan entry, size)
	dropped.Store(0)
	t.Cleanup(func() {
		queue = originalQueue
		dropped.Store(originalDropped)
	})
	return queue
}

func TestHangCapturesValuesAndCaller(t *testing.T) {
	items := useTestQueue(t, 1)
	value := map[string]any{"name": "测试用户", "user_id": 9}
	_, file, line, _ := runtime.Caller(0)
	Hang(111, value)
	value["name"] = "修改后的用户"

	item := <-items
	for _, want := range []string{"[DEBUG]", fmt.Sprintf("%s:%d", file, line+1), "[1] (int): 111", "[2] (map[string]interface {})", "测试用户"} {
		if !strings.Contains(item.text, want) {
			t.Errorf("output %q missing %q", item.text, want)
		}
	}
	if strings.Contains(item.text, "修改后的用户") {
		t.Fatal("output contains a value changed after Hang returned")
	}
}

func TestHangDropsWithoutWaitingForFullQueue(t *testing.T) {
	items := useTestQueue(t, 1)
	items <- entry{text: "existing"}
	done := make(chan struct{})
	go func() {
		Hang("queue is full")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		<-items
		<-done
		t.Fatal("Hang waited for space in the output queue")
	}
	if got := Dropped(); got != 1 {
		t.Fatalf("Dropped() = %d, want 1", got)
	}
	if got := (<-items).text; got != "existing" {
		t.Fatalf("queued output = %q, want existing", got)
	}
}

func TestFlushWaitsForQueuedOutput(t *testing.T) {
	items := useTestQueue(t, 2)
	items <- entry{text: "pending"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- Flush(ctx) }()

	if got := (<-items).text; got != "pending" {
		t.Fatalf("first queued output = %q, want pending", got)
	}
	var barrier entry
	select {
	case barrier = <-items:
	case <-ctx.Done():
		t.Fatal("Flush did not enqueue a completion marker")
	}
	if barrier.done == nil {
		t.Fatal("Flush completion marker is missing")
	}
	select {
	case err := <-result:
		t.Fatalf("Flush returned before the output completed: %v", err)
	default:
	}

	close(barrier.done)
	if err := <-result; err != nil {
		t.Fatalf("Flush() = %v", err)
	}
}

func TestFlushHonorsCancellation(t *testing.T) {
	for _, size := range []int{0, 1} {
		t.Run(fmt.Sprintf("queue_capacity_%d", size), func(t *testing.T) {
			useTestQueue(t, size)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if err := Flush(ctx); !errors.Is(err, context.Canceled) {
				t.Fatalf("Flush() = %v, want context.Canceled", err)
			}
		})
	}

	t.Run("cancel_after_enqueue", func(t *testing.T) {
		items := useTestQueue(t, 1)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		result := make(chan error, 1)
		go func() { result <- Flush(ctx) }()
		select {
		case <-items:
		case <-time.After(5 * time.Second):
			t.Fatal("Flush did not enqueue a completion marker")
		}
		cancel()
		if err := <-result; !errors.Is(err, context.Canceled) {
			t.Fatalf("Flush() = %v, want context.Canceled", err)
		}
	})
}

func TestConcurrentHangAndFlush(t *testing.T) {
	items := useTestQueue(t, queueSize)
	var output bytes.Buffer
	finished := make(chan struct{})
	go func() {
		consume(items, &output)
		close(finished)
	}()
	defer func() {
		close(items)
		<-finished
	}()

	const workers, messages = 16, 20
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < messages; j++ {
				Hang(j)
			}
		}()
	}
	wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Flush(ctx); err != nil {
		t.Fatal(err)
	}

	printed := strings.Count(output.String(), "[DEBUG]")
	if got := uint64(printed) + Dropped(); got != workers*messages {
		t.Fatalf("printed + dropped = %d, want %d", got, workers*messages)
	}
	if got := strings.Count(output.String(), "[1] (int):"); got != printed {
		t.Fatalf("complete messages = %d, want %d", got, printed)
	}
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("output unavailable")
}

func TestConsumeContinuesAfterWriteFailure(t *testing.T) {
	items := make(chan entry, 2)
	done := make(chan struct{})
	items <- entry{text: "cannot write"}
	items <- entry{done: done}
	close(items)
	consume(items, failingWriter{})
	select {
	case <-done:
	default:
		t.Fatal("consumer stopped after a write error")
	}
}
