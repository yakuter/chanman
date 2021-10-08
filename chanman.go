package chanman

import (
	"fmt"
	"sync"
)

var logger Logger = NewBuiltinLogger()

// Chanman is a channel queue manager
type Chanman struct {
	quitCh  chan struct{}
	queueCh chan interface{}
	count   int
	limit   int
	mu      sync.Mutex
}

// New creates a new Chanman instance
func New(quitCh chan struct{}, limit int) *Chanman {
	return &Chanman{
		quitCh:  quitCh,
		queueCh: make(chan interface{}),
		limit:   limit,
	}
}

// Listen starts listening for queue items
func (cm *Chanman) Listen(callbackFn func(interface{})) {
	for {
		select {
		case <-cm.quitCh:
			// logger.Info("Quit signal received. Stopped listening.")
			cm.quit()
			fmt.Println("Quit signal received. Stopped listening.")
			return
		case data, ok := <-cm.queueCh:
			if ok {
				callbackFn(data)
			}
		}
	}
}

// Add adds a new item to the queue
func (cm *Chanman) Add(data interface{}) {
	cm.mu.Lock()
	cm.count += 1
	cm.mu.Unlock()

	if isQueueChClosed(cm.queueCh) {
		logger.Errorf("Failed to add item to queue. Channel is already closed.")
		cm.quit()
		return
	}

	if cm.isLimitExceeded() {
		logger.Errorf("Failed to add %q. Queue limit (%d) exceeded", data, cm.limit)
		cm.quit()
		return
	}

	// It is safe to add item to queue
	cm.queueCh <- data
}

// isLimitExceeded returns true if the channel limit has been reached
func (cm *Chanman) isLimitExceeded() bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.count > cm.limit
}

// isLimitExceeded returns true if the channel limit has been reached
func (cm *Chanman) quit() {
	cm.closeQuitCh()
	cm.closeQueueCh()
}

// closeQueueCh closes a channel gracefully
func (cm *Chanman) closeQueueCh() {
	if !isQueueChClosed(cm.queueCh) {
		close(cm.queueCh)
	}
}

// closeQuitCh closes a channel gracefully
func (cm *Chanman) closeQuitCh() {
	if !isQuitChClosed(cm.quitCh) {
		close(cm.quitCh)
	}
}

// isQueueChClosed returns true if a channel is already closed
func isQueueChClosed(c chan interface{}) bool {
	select {
	case <-c:
		return true
	default:
	}
	return false
}

// isQueueChClosed returns true if a channel is already closed
func isQuitChClosed(c chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
	}
	return false
}
