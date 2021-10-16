package chanman_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/yakuter/chanman"
)

func TestChanman(t *testing.T) {

	var callbackFn chanman.CallbackFn = func(data interface{}) error {
		t.Logf("callbackFn: %v", data)
		return nil
	}

	opts := &chanman.Options{
		CallbackFn: callbackFn,
		Limit:      20,
		Worker:     20,
		DataSize:   0,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := chanman.New(ctx, opts)

	go queue.Listen()

	for i := 0; i <= 20; i++ {
		if i < 5 {
			queue.Add(int64(i))
		} else {
			queue.Add(fmt.Sprintf("%d", i))
		}
		if i == 10 {
			queue.Quit()
		}
	}

	t.Logf("TestChanman done")
}
