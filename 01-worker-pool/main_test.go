package workerpool

import (
	"errors"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func normalProcessor() (func(Job) error, *int) {
	var mu sync.Mutex
	var maxInFlight int
	var inFlight atomic.Int32

	return func(j Job) error {
		defer func() {
			inFlight.Add(-1)
		}()
		inFlight.Add(1)

		mu.Lock()
		flightNow := int(inFlight.Load())
		if flightNow > maxInFlight {
			maxInFlight = flightNow
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)
		return nil
	}, &maxInFlight
}

// Break after N processed
func breakableProcessor(n int) (func(Job) error, *int, *atomic.Int64) {
	var calls atomic.Int64
	var inFlight atomic.Int32
	var mu sync.Mutex
	var maxInFlight int

	return func(j Job) error {
		defer func() {
			inFlight.Add(-1)
		}()
		curCall := calls.Add(1)
		inFlight.Add(1)

		mu.Lock()
		flightNow := int(inFlight.Load())
		if flightNow > maxInFlight {
			maxInFlight = flightNow
		}
		mu.Unlock()

		if next := int(curCall); next > n {
			return errors.New("oopsie")
		}

		return nil
	}, &maxInFlight, &calls
}

func TestRun(t *testing.T) {
	jobs := []Job{"j", "o", "b", "i", "s", "b", "a", "d"}
	maxWorkers := 2

	gorutines := runtime.NumGoroutine()

	processer, maxInFlight := normalProcessor()

	err := RunV1(t.Context(), jobs, maxWorkers, processer)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if runtime.NumGoroutine() > gorutines {
		t.Error("not all gorutines are closed")
	}

	if *maxInFlight > maxWorkers {
		t.Errorf("expected gorutines in flight: %d, got: %d", maxWorkers, maxInFlight)
	}
}

func TestRunWithError(t *testing.T) {
	jobs := []Job{"1", "2", "3", "4", "5", "6", "7", "8"}
	maxWorkers := 2

	processor, maxInFlight, calls := breakableProcessor(2)

	err := RunV1(t.Context(), jobs, maxWorkers, processor)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if *maxInFlight > maxWorkers {
		t.Error("not all gorutines are closed")
	}

	if int(calls.Load()) >= len(jobs) {
		t.Errorf("calls not prematurely, jobs completed: %d", calls.Load())
	}

}

func TestRunThousandJobs(t *testing.T) {
	var jobs []Job
	maxWorkers := 4

	for i := range 1000 {
		jobs = append(jobs, Job(strconv.Itoa(i)))
	}

	processor, maxInFlight := normalProcessor()

	err := RunV1(t.Context(), jobs, maxWorkers, processor)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if *maxInFlight > maxWorkers {
		t.Errorf("expected gorutines in flight: %d, got: %d", maxWorkers, maxInFlight)
	}
}
