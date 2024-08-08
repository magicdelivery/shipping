package datetime

import "time"

type Sleeper interface {
	Sleep(time.Duration)
}

type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// Sleepers for unit tests

type NoDelaySpySleeper struct {
	CallCount uint8
}

func (d *NoDelaySpySleeper) Sleep(duration time.Duration) {
	time.Sleep(0)
	d.CallCount++
}
