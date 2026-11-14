package leaks

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func aliveGorutinesEqualBaseNum(base int) bool {
	ticker := time.NewTicker(time.Millisecond * 50)
	timer := time.NewTimer(time.Second * 2)

	for {
		select {
		case <-ticker.C:
			if runtime.NumGoroutine() <= base {
				return true
			}
		case <-timer.C:
			return false
		}
	}
}

func TestLeaking(t *testing.T) {
	gorutines := runtime.NumGoroutine()

	for range 5 {
		Leaking(t.Context())
	}

	if aliveGorutinesEqualBaseNum(gorutines) {
		t.Error("naive producer should leak")
	}

	if runtime.NumGoroutine() != gorutines+5 {
		t.Error("amount of gorutiness is less than expected")
	}
}

func TestNotLeaking(t *testing.T) {
	gorutines := runtime.NumGoroutine()

	NotLeakingV1(t.Context())

	if !aliveGorutinesEqualBaseNum(gorutines) {
		t.Error("not leaking function is leaking")
	}

}

func TestNotLeakingV2(t *testing.T) {
	gorutines := runtime.NumGoroutine()

	NotLeakingV2(t.Context())

	if !aliveGorutinesEqualBaseNum(gorutines) {
		t.Error("smart producer is leaking")
	}

}

// ----- Test Fan Out -----

func TestFanOut(t *testing.T) {
	var canceledErrors atomic.Int64

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/hold" {
			<-r.Context().Done()
		}

		if r.URL.Path == "/err" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status":"ok"}`))
		}

		if errors.Is(r.Context().Err(), context.Canceled) {
			canceledErrors.Add(1)
		}

	}))

	urls := []string{ts.URL + "/hold", ts.URL + "/hold", ts.URL + "/err"}

	err := FanOut(t.Context(), ts.Client(), urls)
	if err.Error() != "oopsie!" {
		t.Errorf("expected error: %s, got: %v", "oopsie!", err.Error())
	}

	ts.Close() // Blocks until every handler has returned

	if int(canceledErrors.Load()) != 2 {
		t.Errorf("expected canceled errors: 2, got: %d", canceledErrors.Load())
	}
}
