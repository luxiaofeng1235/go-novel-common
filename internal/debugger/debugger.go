/*
 * @Descripttion: 接口临时调试输出
 * @Author: red
 * @Date: 2026-11-13
 * @LastEditors: red
 * @LastEditTime: 2026-11-13
 */

// Package debugger 提供接口临时调试输出。
package debugger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const queueSize = 128

type entry struct {
	text string
	done chan struct{}
}

var (
	queue   = make(chan entry, queueSize)
	dropped atomic.Uint64
)

func init() {
	go consume(queue, os.Stdout)
}

// Hang 打印变量类型、内容和调用位置；队列满时跳过本条输出。
// 变量在调用时同步格式化，终端写入由后台协程处理。
// 跨协程共享的可变数据需由调用方保证读取安全。
func Hang(values ...any) {
	_, file, line, _ := runtime.Caller(1)

	// 先生成字符串，避免后台读取已经被业务修改的对象。
	var text strings.Builder
	fmt.Fprintf(&text, "\n[DEBUG] %s %s:%d\n", time.Now().Format("15:04:05.000"), file, line)
	for i, value := range values {
		fmt.Fprintf(&text, "[%d] (%T): %#v\n", i+1, value, value)
	}

	// 队列满时立即返回，不等待后台输出。
	select {
	case queue <- entry{text: text.String()}:
	default:
		dropped.Add(1)
	}
}

// Dropped 返回队列满时丢弃的调试消息数量。
func Dropped() uint64 {
	return dropped.Load()
}

// Flush 等待此前已入队的消息完成写入尝试，仅用于程序正常退出前收尾。
// 不要在每次接口调用时执行；写入错误会被忽略，等待超时返回 ctx.Err()。
func Flush(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	done := make(chan struct{})
	select {
	case queue <- entry{done: done}:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// consume 使用一个后台协程统一输出。
func consume(items <-chan entry, output io.Writer) {
	for item := range items {
		if item.done != nil {
			close(item.done)
			continue
		}

		// 调试输出写入失败时跳过，不抛出业务异常。
		_, _ = io.WriteString(output, item.text)
	}
}
