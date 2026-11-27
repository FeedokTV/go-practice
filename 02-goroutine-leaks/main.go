package leaks

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// ------- leaking gorutine -------

func naiveProducer(ctx context.Context, ch chan<- int) {
	if ctx.Err() != nil {
		return
	}
	time.Sleep(30 * time.Millisecond)
	ch <- rand.IntN(10)
}

func smartProducer(ctx context.Context, ch chan<- int) {
	sideJob := time.NewTimer(30 * time.Millisecond)
	defer sideJob.Stop()

	select {
	case <-sideJob.C:
	case <-ctx.Done():
		return
	}

	select {
	case ch <- rand.IntN(10):
	case <-ctx.Done():
		return
	}
}

func Leaking(ctx context.Context) {
	jobCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	ch := make(chan int)
	go naiveProducer(jobCtx, ch)

	select {
	case v := <-ch:
		fmt.Println("Got value:", v)
		return
	case <-jobCtx.Done():
		return
	}
}

func NotLeakingV1(ctx context.Context) {
	jobCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	ch := make(chan int, 1)
	go naiveProducer(jobCtx, ch)

	select {
	case v := <-ch:
		fmt.Println("Got value:", v)
		return
	case <-jobCtx.Done():
		return
	}
}

func NotLeakingV2(ctx context.Context) {
	jobCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	ch := make(chan int)
	go smartProducer(jobCtx, ch)

	select {
	case v := <-ch:
		fmt.Println("Got value:", v)
		return
	case <-jobCtx.Done():
		return
	}
}

// ------- errgroup -------

func Request(ctx context.Context, client *http.Client, url string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // close will free resources no matter of consuming error

	if resp.StatusCode == http.StatusInternalServerError {
		return errors.New("oopsie")
	}

	return nil
}

func FanOut(ctx context.Context, client *http.Client, urls []string) error {
	errGroup, reqCtx := errgroup.WithContext(ctx)

	for _, u := range urls {
		errGroup.Go(func() error { return Request(reqCtx, client, u) })
	}

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
