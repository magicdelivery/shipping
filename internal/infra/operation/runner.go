package operation

import (
	"fmt"
	"math/rand"
	"time"

	"shipping/internal/infra/datetime"
)

const RETRY_COUNT_DEF = 5
const MAX_DELAY_DEF = 5 * time.Minute
const BASE_DELAY_DEF = 200 * time.Millisecond

type OperationFunc func() (interface{}, error)

type Runner struct {
	RetryCount uint8
	MaxDelay   time.Duration
	BaseDelay  time.Duration
	Sleeper    datetime.Sleeper
}

// Exponensial backoff algorithm
func (r *Runner) RunWithRetries(op OperationFunc) (interface{}, error) {
	retryMaxCnt := r.retryCount()
	delay := 0 * time.Millisecond
	var attempt uint8 = 0
	for ; attempt < retryMaxCnt; attempt++ {
		result, err := op()
		if err == nil {
			return result, nil
		}

		fmt.Printf("Operation failed: %v. Attempt#: %d, delay: %dms\n", err, attempt+1, delay/time.Millisecond)

		if attempt == retryMaxCnt { // err != nil
			return nil, err
		}

		delay := r.calcDelay(attempt)
		r.sleeper().Sleep(delay)
	}
	return nil, fmt.Errorf("max retries (%d) reached", retryMaxCnt)
}

func (r *Runner) calcDelay(attempt uint8) time.Duration {
	baseDelay := r.baseDelay()
	maxDelay := r.maxDelay()
	var factor time.Duration = 1 << attempt
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	delay := min(baseDelay*factor, maxDelay) + jitter
	return delay
}

func (r *Runner) maxDelay() time.Duration {
	if r.MaxDelay > 0 {
		return r.MaxDelay
	}
	return MAX_DELAY_DEF
}

func (r *Runner) retryCount() uint8 {
	if r.RetryCount > 0 {
		return r.RetryCount
	}
	return RETRY_COUNT_DEF
}

func (r *Runner) baseDelay() time.Duration {
	if r.BaseDelay > 0 {
		return r.BaseDelay
	}
	return BASE_DELAY_DEF
}

func (r *Runner) sleeper() datetime.Sleeper {
	if r.Sleeper != nil {
		return r.Sleeper
	}
	return &datetime.DefaultSleeper{}
}
