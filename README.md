# Experiments with ZeroMQ

```sh
❯ go test -bench=. -benchmem ./...
goos: linux
goarch: amd64
pkg: github.com/leki75/zmqtest
cpu: AMD Ryzen 5 2600X Six-Core Processor           
BenchmarkGoChannel/Pub16-12               454153              2601 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Pub32-12               205632              6911 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Pub64-12               103726             13801 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Pub128-12               39790             36166 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Pub256-12               16857             76988 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Pub512-12                8148            168449 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub16-12               150996              7577 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub32-12                67592             17351 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub64-12                28483             40790 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub128-12                9577            133196 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub256-12                3955            319980 ns/op               0 B/op          0 allocs/op
BenchmarkGoChannel/Sub512-12                1938            647600 ns/op               0 B/op          0 allocs/op
BenchmarkZmq4/Pub16-12                   1889077               635.7 ns/op           332 B/op          7 allocs/op
BenchmarkZmq4/Pub32-12                   1788236               674.4 ns/op           379 B/op          8 allocs/op
BenchmarkZmq4/Pub64-12                   1744162               795.3 ns/op           453 B/op          9 allocs/op
BenchmarkZmq4/Pub128-12                  1961355               697.9 ns/op           391 B/op          8 allocs/op
BenchmarkZmq4/Pub256-12                  1962226               815.5 ns/op           446 B/op          9 allocs/op
BenchmarkZmq4/Pub512-12                  2244709               656.8 ns/op           350 B/op          7 allocs/op
BenchmarkZmq4/Sub16-12                    246169              5134 ns/op            2574 B/op         54 allocs/op
BenchmarkZmq4/Sub32-12                    127975              9015 ns/op            4929 B/op        103 allocs/op
BenchmarkZmq4/Sub64-12                     85738             17323 ns/op            9970 B/op        222 allocs/op
BenchmarkZmq4/Sub128-12                    55180             24664 ns/op           14167 B/op        320 allocs/op
BenchmarkZmq4/Sub256-12                    91666             34751 ns/op           20362 B/op        493 allocs/op
BenchmarkZmq4/Sub512-12                    84315             56840 ns/op           33458 B/op        878 allocs/op
PASS
ok      github.com/leki75/zmqtest       128.386s
```

**NOTE**: We increased the ZMQ_MAX_SOCKETS to support 512 TCP subscribers and use 2 I/O threads.
