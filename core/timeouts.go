package core

import "time"

type Timeouts struct {
	Connect time.Duration
	Read    time.Duration
	Write   time.Duration
}
