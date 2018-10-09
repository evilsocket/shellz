package core

import (
	"errors"
	"time"
)

var (
	TimeoutError = errors.New("timeout")
)

type Timeouts struct {
	Connect time.Duration
	Read    time.Duration
	Write   time.Duration
}

func WithTimeout(tm time.Duration, cb func() interface{}) (error, interface{}) {
	timeout := time.After(tm)
	done := make(chan interface{})
	go func() {
		done <- cb()
	}()

	select {
	case <-timeout:
		return TimeoutError, nil
	case res := <-done:
		if res != nil {
			if e, ok := res.(error); ok {
				return e, nil
			}
		}
		return nil, res
	}
}
