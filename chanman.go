package chanman

import (
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

// NewChanman creates a new Chanman instance
func New(quitCh chan struct{}, limit int) *Chanman {
	return &Chanman{
		quitCh:  quitCh,
		queueCh: make(chan interface{}),
		limit:   limit,
	}
}

// Listen starts listening for queue items
func (cm *Chanman) Listen(callbackFn func(interface{})) {
	defer closeQueueCh(cm.queueCh)

	for {
		select {
		case <-cm.quitCh:
			logger.Info("Quit signal received. Stopped listening.")
			close(cm.quitCh)
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
		return
	}

	if cm.IsLimitExceeded() {
		logger.Errorf("Failed to add %q. Queue limit (%d) exceeded", data, cm.limit)
		closeQueueCh(cm.queueCh)
		return
	}

	// It is safe to add item to queue
	cm.queueCh <- data
}

// IsLimitExceeded returns true if the channel limit has been reached
func (cm *Chanman) IsLimitExceeded() bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	return cm.count > cm.limit
}

// CloseCh closes a channel gracefully
func closeQueueCh(c chan interface{}) {
	if !isQueueChClosed(c) {
		close(c)
	}
}

// IsChClosed returns true if a channel is already closed
func isQueueChClosed(c chan interface{}) bool {
	select {
	case <-c:
		return true
	default:
	}
	return false
}
