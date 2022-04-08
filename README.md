# Experiments with ZeroMQ

## Local test on MacOS X

```sh
❯ go test -bench=. -benchmem ./...
goos: darwin
goarch: amd64
pkg: github.com/leki75/zmqtest
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGoChannel_recv-12            	 4762549	       255.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoChannel_send-12            	 6585050	       176.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkCzmqSocket_recv_inproc-12    	 1387824	       798.6 ns/op	     128 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_inproc-12    	 1389370	       880.0 ns/op	     127 B/op	       6 allocs/op
BenchmarkCzmqSocket_recv_ipc-12       	 1415056	       851.6 ns/op	     128 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_ipc-12       	 1748611	       839.0 ns/op	     127 B/op	       6 allocs/op
BenchmarkCzmqSocket_recv_tcp-12       	 1356379	       957.0 ns/op	     128 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_tcp-12       	 1280512	       868.4 ns/op	     126 B/op	       6 allocs/op
BenchmarkZmq4Socket_recv_inproc-12    	 1649125	       647.5 ns/op	     304 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_inproc-12    	 1985466	       681.0 ns/op	     303 B/op	       5 allocs/op
BenchmarkZmq4Socket_recv_ipc-12       	 1854087	       715.4 ns/op	     304 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_ipc-12       	 2020258	       577.6 ns/op	     303 B/op	       5 allocs/op
BenchmarkZmq4Socket_recv_tcp-12       	 1431608	       844.2 ns/op	     305 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_tcp-12       	 1455072	       787.1 ns/op	     299 B/op	       5 allocs/op
PASS
ok  	github.com/leki75/zmqtest	28.108s
```

## Docker test on MacOS X, Linux container

```sh
❯ docker build -t zmqtest .
❯ docker run -t zmqtest
goos: linux
goarch: amd64
pkg: github.com/leki75/zmqtest
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkGoChannel_recv-6           	  464943	      2297 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoChannel_send-6           	  489412	      2198 ns/op	       0 B/op	       0 allocs/op
BenchmarkCzmqSocket_recv_inproc-6   	 1539637	       781.0 ns/op	     128 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_inproc-6   	 1523367	       795.1 ns/op	     127 B/op	       6 allocs/op
BenchmarkCzmqSocket_recv_ipc-6      	 1573741	       763.7 ns/op	     128 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_ipc-6      	 1562056	       773.0 ns/op	     127 B/op	       6 allocs/op
BenchmarkCzmqSocket_recv_tcp-6      	 1233158	       986.8 ns/op	     130 B/op	       7 allocs/op
BenchmarkCzmqSocket_send_tcp-6      	 1501854	       883.3 ns/op	     121 B/op	       6 allocs/op
BenchmarkZmq4Socket_recv_inproc-6   	 1675046	       715.2 ns/op	     304 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_inproc-6   	 1653790	       730.0 ns/op	     303 B/op	       5 allocs/op
BenchmarkZmq4Socket_recv_ipc-6      	 1676445	       699.0 ns/op	     304 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_ipc-6      	 1733248	       687.5 ns/op	     303 B/op	       5 allocs/op
BenchmarkZmq4Socket_recv_tcp-6      	 1727510	       674.2 ns/op	     305 B/op	       6 allocs/op
BenchmarkZmq4Socket_send_tcp-6      	 1791862	       672.3 ns/op	     300 B/op	       5 allocs/op
PASS
ok  	github.com/leki75/zmqtest	26.063s
```
