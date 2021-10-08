package chanman_test

import (
	"testing"

	"github.com/yakuter/chanman"
)

func TestChanman(t *testing.T) {
	t.Log("TestChanman")

	quitCh := make(chan struct{})
	chanmanQueue := chanman.New(quitCh, 4)

	callbackFn := func(data interface{}) {
		t.Logf("Received data: %v", data)
	}

	go chanmanQueue.Listen(callbackFn)

	chanmanQueue.Add("Test Data 1")
	chanmanQueue.Add("Test Data 2")
	quitCh <- struct{}{}
	chanmanQueue.Add("Test Data 3")
	chanmanQueue.Add("Test Data 4")
	chanmanQueue.Add("Test Data 5")

	<-quitCh
	t.Logf("TestChanman done")
}
