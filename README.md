# Experiments with ZeroMQ

```sh
❯ go test -bench=. -benchmem ./*.go 
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGoChannel_recv-12       2114181               587.1 ns/op             0 B/op          0 allocs/op
BenchmarkGoChannel_send-12       2929124               407.9 ns/op             0 B/op          0 allocs/op
BenchmarkCzmqSocket/Recv-Inproc-12                629619              1633 ns/op             128 B/op          7 allocs/op
BenchmarkCzmqSocket/Send-Inproc-12               1275244               937.4 ns/op            51 B/op          4 allocs/op
BenchmarkCzmqSocket/Recv-IPC-12                   625003              1733 ns/op             128 B/op          7 allocs/op
BenchmarkCzmqSocket/Send-IPC-12                  1277078               944.8 ns/op            59 B/op          5 allocs/op
BenchmarkCzmqSocket/Recv-TCP-12                   626035              1945 ns/op             129 B/op          7 allocs/op
BenchmarkCzmqSocket/Send-TCP-12                  1254134               957.5 ns/op            61 B/op          5 allocs/op
BenchmarkZmq4Socket/Recv-Inproc-12                768044              1400 ns/op             176 B/op          4 allocs/op
BenchmarkZmq4Socket/Send-Inproc-12                884361              1430 ns/op             175 B/op          3 allocs/op
BenchmarkZmq4Socket/Recv-IPC-12                   822829              1325 ns/op             176 B/op          4 allocs/op
BenchmarkZmq4Socket/Send-IPC-12                   910080              1279 ns/op             175 B/op          3 allocs/op
BenchmarkZmq4Socket/Recv-TCP-12                   896509              1357 ns/op             176 B/op          4 allocs/op
BenchmarkZmq4Socket/Send-TCP-12                  1000000              1422 ns/op             171 B/op          3 allocs/op
PASS
ok      command-line-arguments  22.957s
```
