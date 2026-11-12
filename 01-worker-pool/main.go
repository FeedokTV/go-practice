package workerpool

import (
	"context"
	"sync"
)

type Job string

// Approach usually comes in API's like REST when client give's us a complete
// list of jobs
// It's appropriate to get chan from caller when it's queues or streams (Kafka/gRPC stream)

func worker(ctx context.Context, jobs <-chan Job, process func(Job) error, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()
	var err error

	for j := range jobs {
		if ctx.Err() != nil {
			return
		}
		err = process(j)
		if err != nil {
			select {
			case errors <- err:
			case <-ctx.Done():
			}
			return
		}
	}

}

func RunV1(ctx context.Context, in []Job, maxWorkers int, process func(Job) error) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan Job, 10)
	errors := make(chan error, 1)

	var wg sync.WaitGroup

	for range maxWorkers {
		wg.Add(1)
		go worker(runCtx, jobs, process, &wg, errors)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
	loop:
		for _, j := range in {

			select {
			case jobs <- j:
			case <-runCtx.Done():
				break loop
			}

		}
		close(jobs)
	}()

	// Signal that everything finished succesfully
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case err := <-errors: // Signal that error was occured
		cancel()
		wg.Wait()
		close(errors)
		return err
	}

	return nil
}
