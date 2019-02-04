package core

import (
	"time"
)

type Timeouts struct {
	Connect time.Duration
	Read    time.Duration
	Write   time.Duration
}

func (tt Timeouts) RW() time.Duration {
	return tt.Read + tt.Write
}

func (tt Timeouts) Total() time.Duration {
	return tt.Connect + tt.RW()
}
