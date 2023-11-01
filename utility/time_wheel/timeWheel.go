// Package time_wheel 类似时间轮结构，用于添加、移除和执行定时任务。
// 你可以创建一个 DefaultTimeWheel 实例，添加任务，启动时间轮，并在指定时间触发任务的执行。
// 当不再需要时间轮时，可以停止它以释放资源
package time_wheel

import (
	"context"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

// TaskHandler 通用的函数类型，用于处理定时任务
type TaskHandler[T any] func(*DefaultTimeWheel[T], string, T)

// DefaultTimeWheel 时间轮定义
type DefaultTimeWheel[T any] struct {
	ctx      context.Context
	slot     cmap.ConcurrentMap[string, *gtimer.Timer]
	quitChan chan struct{}
	taskChan chan bool
	onTick   TaskHandler[T]
}

// NewTimeWheel 创建一个时间轮实例
func NewTimeWheel[T any](handler TaskHandler[T]) *DefaultTimeWheel[T] {
	timeWheel := &DefaultTimeWheel[T]{
		ctx:      gctx.New(),
		slot:     cmap.New[*gtimer.Timer](),
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
			t.slot.IterCb(func(_ string, v *gtimer.Timer) {
				v.Close()
			})
			return
		}
	}
}

// Stop 停止时间轮，它向退出通道 quitChan 发送关闭信号，触发 Start 方法中的退出操作，停止时间轮的运行
func (t *DefaultTimeWheel[T]) Stop() {
	close(t.quitChan)
}

// AddOnce 添加执行一次的定时任务
func (t *DefaultTimeWheel[T]) AddOnce(key string, delay time.Duration, value T) {
	timer := gtimer.New()
	timer.AddOnce(t.ctx, delay, func(ctx context.Context) {
		t.onTick(t, key, value)
	})

	if t.slot.Has(key) {
		t.slot.Remove(key)
	}

	t.slot.Set(key, timer)
}

// Remove 移除指定名称的任务
func (t *DefaultTimeWheel[T]) Remove(name string) {
	if timer, ok := t.slot.Get(name); ok {
		timer.Close()
		t.slot.Remove(name)
	}

}
