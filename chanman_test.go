package chanman_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/yakuter/chanman"
)

func TestChanman(t *testing.T) {
	callbackFn := func(data interface{}) error {
		t.Logf("Processed data: %v", data)
		return nil
	}

	opts := &chanman.Options{
		CallbackFn: callbackFn,
		Limit:      19,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := chanman.New(ctx, opts)

	go queue.Listen()

	for i := 0; i <= 20; i++ {
		queue.Add(fmt.Sprintf("job-%d", i))
		if i == 10 {
			queue.Quit()
		}
	}

	t.Logf("TestChanman done")
}
