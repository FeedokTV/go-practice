package racesandsync

import (
	"sync"
	"sync/atomic"
)

func RacyCounter() (func(), *int) {
	var counter int

	return func() {
		counter++
	}, &counter
}

func MutexCounter() (func(), *sync.Mutex, *int) {
	var mu sync.Mutex
	var counter int

	return func() {
		mu.Lock()
		defer mu.Unlock()
		counter++
	}, &mu, &counter
}

func AtomicCounter() (func(), *atomic.Int64) {
	var counter atomic.Int64

	return func() {
		counter.Add(1)
	}, &counter
}
