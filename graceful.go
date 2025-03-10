// Package graceful provides a graceful shutdown manager for Go applications.
// It helps manage the lifecycle of goroutines and ensures they are properly cleaned up
// when the application receives termination signals or needs to shut down.
//
// The package offers a Manager type that coordinates goroutine lifecycles by:
// - Managing goroutine startup and shutdown
// - Handling OS signals (e.g., SIGINT, SIGTERM)
// - Providing timeout-based shutdown
// - Ensuring clean resource cleanup
package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Manager handles graceful shutdown of goroutines in an application.
// It provides mechanisms to start goroutines, monitor their lifecycle,
// and ensure they shut down cleanly when the application needs to terminate.
type Manager struct {
	ctx        context.Context    // Context for coordinating goroutine lifecycle
	cancelFunc context.CancelFunc // Function to cancel the context
	wg         sync.WaitGroup     // WaitGroup for tracking active goroutines
	timeout    time.Duration      // Maximum time to wait for goroutines to exit
	signals    []os.Signal        // OS signals to monitor for shutdown
}

// Option defines a function type for configuring Manager instances.
// It follows the functional options pattern for flexible configuration.
type Option func(*Manager)

// WithTimeout returns an Option that sets the maximum duration to wait
// for goroutines to exit during shutdown. If goroutines do not exit within
// this time, the manager will proceed with shutdown anyway.
//
// Example:
//
//	manager := graceful.New(graceful.WithTimeout(5 * time.Second))
func WithTimeout(timeout time.Duration) Option {
	return func(m *Manager) {
		m.timeout = timeout
	}
}

// WithSignals returns an Option that sets which OS signals the manager should
// monitor for triggering graceful shutdown. By default, the manager monitors
// SIGINT and SIGTERM.
//
// Example:
//
//	manager := graceful.New(graceful.WithSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP))
func WithSignals(signals ...os.Signal) Option {
	return func(m *Manager) {
		m.signals = signals
	}
}

// New creates a new Manager instance with the provided options.
// It initializes the manager with default settings that can be overridden
// through the provided options.
//
// Default settings:
// - Timeout: 30 seconds
// - Signals: SIGINT and SIGTERM
//
// Example:
//
//	manager := graceful.New(
//		graceful.WithTimeout(5 * time.Second),
//		graceful.WithSignals(syscall.SIGINT, syscall.SIGTERM),
//	)
func New(options ...Option) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		ctx:        ctx,
		cancelFunc: cancel,
		timeout:    time.Second * 30,                             // Default timeout: 30 seconds
		signals:    []os.Signal{syscall.SIGINT, syscall.SIGTERM}, // Default signals
	}

	for _, option := range options {
		option(m)
	}

	return m
}

// Go starts a new managed goroutine using the manager's context.
// The provided function will be executed in a new goroutine and will
// automatically use the manager's context. This is a convenience method
// that simplifies goroutine creation when no custom context is needed.
//
// Example:
//
//	manager.Go(func() {
//		// Do work
//		// The function will be stopped when manager initiates shutdown
//	})
func (m *Manager) Go(f func()) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		f()
	}()
}

// Wait blocks until one of the monitored signals is received, then initiates
// graceful shutdown. It notifies all managed goroutines to exit and waits
// for them to complete or for the timeout to expire.
//
// This method is typically called in the main function after starting all
// goroutines.
//
// Example:
//
//	func main() {
//		manager := graceful.New()
//		// Start goroutines...
//		manager.Wait() // Block until signal received
//	}
func (m *Manager) Wait() {
	// Create a signal channel
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, m.signals...)

	// Wait for signal
	<-sigCh

	// Stop receiving signals
	signal.Stop(sigCh)

	// Notify all goroutines to exit and wait for completion
	m.waitForGoroutines()
}

// Shutdown initiates graceful shutdown without waiting for signals.
// It notifies all managed goroutines to exit and waits for them to complete
// or for the timeout to expire.
//
// This method is useful when you need to programmatically shut down the
// application.
//
// Example:
//
//	if err != nil {
//		// Handle error and shut down
//		manager.Shutdown()
//	}
func (m *Manager) Shutdown() {
	// Notify all goroutines to exit and wait for completion
	m.waitForGoroutines()
}

// waitForGoroutines handles the graceful shutdown process by canceling
// the context and waiting for all goroutines to exit or for the timeout
// to expire.
func (m *Manager) waitForGoroutines() {
	// Notify all goroutines to exit
	m.cancelFunc()

	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	// Wait for all goroutines to exit or timeout
	c := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		// All goroutines have exited
	case <-timeoutCtx.Done():
		// Timeout occurred
	}
}

// Context returns the manager's context, which is canceled when shutdown
// begins. This context can be used to derive child contexts or passed
// directly to functions that accept a context.
//
// Example:
//
//	ctx := manager.Context()
//	// Use ctx to create child contexts or pass to functions
func (m *Manager) Context() context.Context {
	return m.ctx
}

// CtxGo starts a new managed goroutine. The provided function receives a context
// that will be canceled when the manager initiates shutdown. Goroutines should
// monitor this context and exit when it's canceled.
//
// Example:
//
//	manager.CtxGo(func(ctx context.Context) {
//		for {
//			select {
//			case <-ctx.Done():
//				// Clean up and exit
//				return
//			default:
//				// Do work
//			}
//		}
//	})
func (m *Manager) CtxGo(f func(ctx context.Context)) {
	m.Go(func() {
		f(m.ctx)
	})
}
