# Graceful - Graceful Shutdown Library for Go

`graceful` is a lightweight Go library for managing goroutine lifecycles and implementing graceful shutdown. It provides a simple API to help developers safely terminate all goroutines when the application shuts down, ensuring proper resource cleanup.

## Features

- Simple and intuitive API
- Signal handling (e.g., SIGINT, SIGTERM)
- Configurable timeout mechanism
- Custom signal handling support
- Graceful goroutine termination

## Installation

```bash
go get github.com/kingcanfish/graceful
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kingcanfish/graceful"
)

func main() {
	// Create a new Manager instance
	manager := graceful.New(graceful.WithTimeout(5 * time.Second))

	// Start a worker goroutine
	manager.CtxGo(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Received shutdown signal, cleaning up...")
				// Perform cleanup
				time.Sleep(time.Second)
				fmt.Println("Cleanup completed")
				return
			default:
				// Do work
				fmt.Println("Working...")
				time.Sleep(time.Second)
			}
		}
	})

	// Wait for signal
	fmt.Println("Application started. Press Ctrl+C to exit")
	manager.Wait() // Blocks until signal is received
	fmt.Println("Application has been shut down gracefully")
}
```

## API Documentation

### Manager

`Manager` is responsible for managing goroutine lifecycles, including starting, monitoring, and shutting down goroutines.

```go
type Manager struct {
	// Internal fields
}
```

### Creating a Manager

```go
func New(options ...Option) *Manager
```

Creates a new Manager instance with configurable options.

### Configuration Options

```go
// Set the timeout duration for waiting goroutines to exit
func WithTimeout(timeout time.Duration) Option

// Set the signals to monitor
func WithSignals(signals ...os.Signal) Option
```

### Starting Goroutines

```go
// Start a goroutine without context
func (m *Manager) Go(f func())

// Start a goroutine with context
func (m *Manager) CtxGo(f func(ctx context.Context))
```

Starts a managed goroutine. The `CtxGo` version provides a context that will be canceled when the Manager initiates shutdown.

### Waiting for Signals

```go
func (m *Manager) Wait()
```

Blocks until a configured signal is received (default: SIGINT and SIGTERM), then notifies all goroutines to exit and waits for their completion.

### Manual Shutdown

```go
func (m *Manager) Shutdown()
```

Initiates shutdown manually without waiting for signals.

### Getting Context

```go
func (m *Manager) Context() context.Context
```

Returns the Manager's context, which can be used to derive child contexts.

## Best Practices

1. Regularly check context cancellation in goroutines
2. Perform necessary cleanup when receiving shutdown signals
3. Set appropriate timeout durations to prevent hanging
4. Use `Wait()` in the main function to handle signals
5. Prefer `CtxGo` over `Go` when you need shutdown notification

## License

MIT

## Powered By

Trae & Claude-3.5-Sonnet