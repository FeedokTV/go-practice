# 08 - Profiling and benchmarks

## Results

### First

```
BenchmarkBuildStringV1-16   103.4 ns/op   72 B/op   5 allocs/op
BenchmarkBuildStringV2-16    57.53 ns/op  56 B/op   3 allocs/op
```

As we can see, the first version (with the `+=` operation) was much slower and took much more resources. It took 103.4 nanoseconds per operation, needed 72 bytes per operation and made 5 allocations per operation.

The second one took 57.53 nanoseconds and 56 bytes per operation and made 3 allocations per operation.

The first one is much heavier (let's say roughly "twice as heavy").

And the Builder has its own logic.

If we run `go build -gcflags=-m ./...`, we can see:

```
08-profiling-benchmarks/main.go:5:20: leaking param: strs to result ~r0 level=1

08-profiling-benchmarks/main.go:15:20: strs does not escape
08-profiling-benchmarks/main.go:19:16: "strings: illegal use of non-zero Builder copied by value" escapes to heap
08-profiling-benchmarks/main.go:19:16: append escapes to heap
```

V2 has only three calls because it uses the Builder internally, and fewer because `WriteString` actually copies into the Builder's own internal buffer.

This is because strings are immutable in Go, and each `+=` makes a new string at every step. We have 5 allocations because the first step is not an allocation: it just points the `res` header at the bytes of "This". So `res` has a pointer to the same underlying bytes as `strs[0]`.

The three allocations are made by `append`, because the underlying byte slice in the Builder is growing. This is exactly what "append escapes to heap" is trying to say. The allocations there are allocations of a new underlying array for the slice.

So in that case, if we read the `-gcflags=-m` output: the compiler says "leaking param: strs to result", which means that data from `strs` may leak into the result. In the benchmark we don't see it, but if the passed array has one string, the result string will point to the bytes of `strs[0]`, so they won't be deleted until the result is deleted. (Am I right, or what does it cause?)

`strs does not escape` because we copy the string header with `append`.

### Second

The pprof heap profile shows sampled heap memory allocations:

```
Showing nodes accounting for 1628.23kB, 100% of 1628.23kB total
      flat  flat%   sum%        cum   cum%
 1115.23kB 68.49% 68.49%  1115.23kB 68.49%  main.LoaderHandler
     513kB 31.51%   100%      513kB 31.51%  runtime.mallocgc
         0     0%   100%  1115.23kB 68.49%  net/http.(*ServeMux).ServeHTTP
         0     0%   100%  1115.23kB 68.49%  net/http.(*conn).serve
```

The command `go tool pprof -inuse_space -top` analyzes a heap profile and prints a text table of the functions holding the largest amount of live RAM.

Its result is below:

```
Showing nodes accounting for 152.13MB, 98.06% of 155.13MB total
Dropped 39 nodes (cum <= 0.78MB)
      flat  flat%   sum%        cum   cum%
   49.51MB 31.92% 31.92%    49.51MB 31.92%  net/textproto.readMIMEHeader
   28.51MB 18.38% 50.29%    92.02MB 59.32%  net/http.readRequestLimit
   23.51MB 15.15% 65.45%   130.03MB 83.82%  net/http.(*conn).readRequest
      12MB  7.74% 73.18%       12MB  7.74%  context.withCancel (inline)
      11MB  7.09% 80.27%       11MB  7.09%  net/url.parse
       6MB  3.87% 84.14%       10MB  6.45%  net.(*conn).Read
    4.66MB  3.00% 87.14%     7.16MB  4.62%  main.LoaderHandler
    2.50MB  1.61% 88.76%     2.50MB  1.61%  unicode/utf16.Encode
    2.02MB  1.30% 90.06%     2.02MB  1.30%  compress/flate.(*fastGen).addBlock
       2MB  1.29% 91.35%        2MB  1.29%  net/textproto.(*Reader).ReadLine (inline)
       2MB  1.29% 92.64%       14MB  9.03%  context.WithCancel
       2MB  1.29% 93.92%        2MB  1.29%  internal/syscall/windows.errnoErr (inline)
    1.50MB  0.97% 94.89%     1.50MB  0.97%  sync.(*Pool).pinSlow
    1.50MB  0.97% 95.86%     1.50MB  0.97%  internal/poll.init.func7
    1.28MB  0.83% 96.68%     2.42MB  1.56%  compress/flate.NewWriter (inline)
    1.14MB  0.73% 97.42%     1.14MB  0.73%  compress/flate.newFastEnc (inline)
       1MB  0.65% 98.06%     2.14MB  1.38%  runtime/pprof.(*profileBuilder).emitLocation
         0     0% 98.06%     2.02MB  1.30%  compress/flate.(*Writer).Close (inline)
```

We can see that `main.LoaderHandler` allocated 1115.23kB of memory, while `runtime.mallocgc` allocated 513kB.

If we call `go tool pprof -top -lines http://localhost:6060/debug/pprof/heap`, we can see:

```
 1115.23kB 68.49% 68.49%  1115.23kB 68.49%  main.LoaderHandler .../go-practice\08-profiling-benchmarks\main.go:42
```

Our line is:

```go
garbage = append(garbage, &n)
```

After the function finishes, and even if the connection is dropped, each of the variables is still alive, because `garbage` still holds a reference to it.
