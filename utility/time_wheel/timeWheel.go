// Package time_wheel 类似时间轮结构，用于添加、移除和执行定时任务。
// 你可以创建一个 DefaultTimeWheel 实例，添加任务，启动时间轮，并在指定时间触发任务的执行。
// 当不再需要时间轮时，可以停止它以释放资源
package time_wheel

import (
	"context"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
)

// TaskHandler 通用的函数类型，用于处理定时任务
type TaskHandler[T any] func(*DefaultTimeWheel[T], string, T)

// DefaultTimeWheel 时间轮定义
type DefaultTimeWheel[T any] struct {
	ctx      context.Context
	cron     *gcron.Cron
	quitChan chan struct{}
	taskChan chan bool
	onTick   TaskHandler[T]
}

// NewTimeWheel 创建一个时间轮实例
func NewTimeWheel[T any](handler TaskHandler[T]) *DefaultTimeWheel[T] {
	cron := gcron.New()

	timeWheel := &DefaultTimeWheel[T]{
		ctx:      gctx.New(),
		cron:     cron,
		quitChan: make(chan struct{}),
		onTick:   handler,
	}
	return timeWheel
}

// Start 时间轮的启动方法，它包含一个无限循环，通过监听退出通道
func (t *DefaultTimeWheel[T]) Start() {
	for {
		select {
		case <-t.quitChan:
			entries := t.cron.Entries()
			for _, v := range entries {
				t.cron.Remove(v.Name)
			}
			return
		}
	}
}

// Stop 停止时间轮，它向退出通道 quitChan 发送关闭信号，触发 Start 方法中的退出操作，停止时间轮的运行
func (t *DefaultTimeWheel[T]) Stop() {
	close(t.quitChan)
}

// Add 添加定时任务
func (t *DefaultTimeWheel[T]) Add(key, delay string, value T) {
	_, err := t.cron.Add(t.ctx, delay, func(ctx context.Context) {
		t.onTick(t, key, value)
	}, key)
	if err != nil {
		panic(err)
	}
}

// AddOnce 添加执行一次的定时任务
func (t *DefaultTimeWheel[T]) AddOnce(key string, delay time.Duration, value T) {
	t.cron.DelayAddOnce(t.ctx, delay, "* * * * * *", func(ctx context.Context) {
		t.onTick(t, key, value)
	}, key)
}

// Remove 移除指定名称的任务
func (t *DefaultTimeWheel[T]) Remove(name string) {
	t.cron.Remove(name)
}
