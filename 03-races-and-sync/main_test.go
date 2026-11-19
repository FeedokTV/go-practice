package racesandsync

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestRacyCounter(t *testing.T) {
	shouldStr := os.Getenv("GOPRACTICE_03_TEST_RACY")

	if should, err := strconv.ParseBool(shouldStr); !should || err != nil {
		return
	}

	inc, _ := RacyCounter()

	var wg sync.WaitGroup

	for range 3 {
		wg.Go(func() {
			for range 3 {
				inc()
			}
		})
	}

	wg.Wait()

	// Waiting for the race checker to say it's word
}

func TestMutexCounter(t *testing.T) {
	inc, mu, counter := MutexCounter()

	var wg sync.WaitGroup

	for range 3 {
		wg.Go(func() {
			for range 3 {
				inc()
			}
		})
	}

	wg.Wait()

	mu.Lock() // Not necessary here, but this how it should be done!
	defer mu.Unlock()
	if *counter != 9 {
		t.Errorf("counter got: %d, expected: 9. probably race", *counter)
	}

}

func TestAtomicCounter(t *testing.T) {
	inc, counter := AtomicCounter()

	var wg sync.WaitGroup

	for range 3 {
		wg.Go(func() {
			for range 3 {
				inc()
			}
		})
	}

	wg.Wait()

	if counter.Load() != 9 {
		t.Errorf("counter got: %d, expected: 9. probably race", counter.Load())
	}

}
