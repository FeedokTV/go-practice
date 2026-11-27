package vetandlinters

// // lostcancel: the cancel function returned by context.WithTimeout should be called, not discarded, to avoid a context leak (govet)
// func lostcancel() {
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	ctx.Done()
// }

// ---- Copy Locks -----
// type SafeCounter struct {
// 	mu sync.Mutex
// 	v  int
// }

// // copylocks: Inc passes lock by value: go-practice/09-vet-and-linters.SafeCounter contains sync.Mutex (govet)
// func (c SafeCounter) Inc() {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.v++
// }

//  printf: fmt.Printf format %d has arg "oopsie!" of wrong type string
// func printf() {
// 	fmt.Printf("%d", "oopsie!")
// }
