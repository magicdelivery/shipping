package operation

import (
	"errors"
	"shipping/internal/infra/datetime"
	"testing"
)

// ----------------------------------------------------------------------------------------------------

func TestRunWithRetriesForSuccess(t *testing.T) {
	spySleeper := datetime.NoDelaySpySleeper{}
	r := &Runner{
		Sleeper: &spySleeper,
	}

	// Test successful operation
	result, err := r.RunWithRetries(successfulOperation)
	if err != nil || result != "success" {
		t.Errorf("Expected success, got result: %v, error: %v", result, err)
	}

	if spySleeper.CallCount != 0 {
		t.Errorf("Do not expect retries, got %d", spySleeper.CallCount)
	}
}

// Mock OperationFunc that succeeds immediately
func successfulOperation() (interface{}, error) {
	return "success", nil
}

// ----------------------------------------------------------------------------------------------------

func TestRunWithRetriesForSuccessAfterFails(t *testing.T) {
	const FAIL_COUNT = 2
	spySleeper := datetime.NoDelaySpySleeper{}
	r := &Runner{
		Sleeper: &spySleeper,
	}

	// Test operation that fails twice before succeeding
	result, err := r.RunWithRetries(failingThenSuccessOperation(FAIL_COUNT))
	if err != nil || result != "success" {
		t.Errorf("Expected success after retries, got result: %v, error: %v", result, err)
	}

	if spySleeper.CallCount != FAIL_COUNT {
		t.Errorf("Expected %d retries, got %d", FAIL_COUNT, spySleeper.CallCount)
	}
}

// Mock OperationFunc that fails a specified number of times before succeeding
func failingThenSuccessOperation(failCount int) OperationFunc {
	attempt := 0
	return func() (interface{}, error) {
		if attempt < failCount {
			attempt++
			return nil, errors.New("temporary error")
		}
		return "success", nil
	}
}

// ----------------------------------------------------------------------------------------------------

func TestRunWithRetriesForFail(t *testing.T) {
	spySleeper := datetime.NoDelaySpySleeper{}
	r := &Runner{
		Sleeper: &spySleeper,
	}

	// Test operation that always fails
	result, err := r.RunWithRetries(alwaysFailingOperation)
	if err == nil || result != nil {
		t.Errorf("Expected failure, got result: %v, error: %v", result, err)
	}

	if spySleeper.CallCount != RETRY_COUNT_DEF {
		t.Errorf("Expected %d retries, got %d", RETRY_COUNT_DEF, spySleeper.CallCount)
	}
}

// Mock OperationFunc that always fails
func alwaysFailingOperation() (interface{}, error) {
	return nil, errors.New("permanent error")
}

// ----------------------------------------------------------------------------------------------------
