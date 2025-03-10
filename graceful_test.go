package graceful

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"
)

// TestNewManager 测试创建新的Manager实例
func TestNewManager(t *testing.T) {
	// 测试默认配置
	m := New()
	if m.timeout != time.Second*30 {
		t.Errorf("默认超时时间应为30秒，实际为%v", m.timeout)
	}
	if len(m.signals) != 2 {
		t.Errorf("默认信号数量应为2，实际为%d", len(m.signals))
	}

	// 测试自定义配置
	customTimeout := time.Second * 10
	customSignals := []os.Signal{syscall.SIGHUP}
	m = New(
		WithTimeout(customTimeout),
		WithSignals(customSignals...),
	)

	if m.timeout != customTimeout {
		t.Errorf("自定义超时时间应为%v，实际为%v", customTimeout, m.timeout)
	}
	if len(m.signals) != len(customSignals) {
		t.Errorf("自定义信号数量应为%d，实际为%d", len(customSignals), len(m.signals))
	}
}

// TestManagerGo 测试启动受管理的goroutine
func TestManagerGo(t *testing.T) {
	m := New(WithTimeout(time.Second))

	// 测试goroutine是否正常执行
	result := make(chan int, 1)
	m.CtxGo(func(ctx context.Context) {
		result <- 42
	})

	select {
	case val := <-result:
		if val != 42 {
			t.Errorf("goroutine返回值应为42，实际为%d", val)
		}
	case <-time.After(time.Second * 2):
		t.Error("goroutine执行超时")
	}
}

// TestManagerContext 测试获取Manager的上下文
func TestManagerContext(t *testing.T) {
	m := New()
	ctx := m.Context()
	if ctx == nil {
		t.Error("Context()返回的上下文不应为nil")
	}
}

// TestManagerShutdown 测试主动关闭所有goroutine
func TestManagerShutdown(t *testing.T) {
	m := New(WithTimeout(time.Second))

	// 启动一个goroutine，检测是否收到退出信号
	exitChan := make(chan struct{}, 1)
	m.CtxGo(func(ctx context.Context) {
		<-ctx.Done()
		exitChan <- struct{}{}
	})

	// 等待goroutine启动
	time.Sleep(time.Millisecond * 100)

	// 主动关闭
	m.Shutdown()

	// 检查goroutine是否收到退出信号
	select {
	case <-exitChan:
		// 成功收到退出信号
	case <-time.After(time.Second * 2):
		t.Error("goroutine未收到退出信号")
	}
}

// TestManagerWaitTimeout 测试等待超时情况
func TestManagerWaitTimeout(t *testing.T) {
	// 创建一个非常短的超时时间
	m := New(WithTimeout(time.Millisecond * 50))

	// 启动一个不会立即退出的goroutine
	m.CtxGo(func(ctx context.Context) {
		// 忽略ctx.Done()，模拟一个无法立即退出的goroutine
		time.Sleep(time.Second)
	})

	// 记录开始时间
	start := time.Now()

	// 主动关闭
	m.Shutdown()

	// 检查是否在超时时间附近返回
	duration := time.Since(start)
	if duration < time.Millisecond*50 {
		t.Errorf("应该等待至少50ms，实际等待了%v", duration)
	}
	if duration > time.Millisecond*200 {
		t.Errorf("应该在超时后立即返回，实际等待了%v", duration)
	}
}

// TestMultipleGoroutines 测试多个goroutine的情况
func TestMultipleGoroutines(t *testing.T) {
	m := New(WithTimeout(time.Second))

	// 计数器，用于记录已退出的goroutine数量
	counter := 0
	counterCh := make(chan struct{}, 5)

	// 启动5个goroutine
	for i := 0; i < 5; i++ {
		m.CtxGo(func(ctx context.Context) {
			<-ctx.Done()
			counterCh <- struct{}{}
		})
	}

	// 主动关闭
	m.Shutdown()

	// 等待所有goroutine退出或超时
	timeout := time.After(time.Second * 2)
Loop:
	for i := 0; i < 5; i++ {
		select {
		case <-counterCh:
			counter++
		case <-timeout:
			break Loop
		}
	}

	if counter != 5 {
		t.Errorf("应有5个goroutine退出，实际有%d个退出", counter)
	}
}
