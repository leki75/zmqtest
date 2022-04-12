package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/leki75/zmqtest/common"
	"github.com/zeromq/goczmq"
)

func BenchmarkGoczmq(b *testing.B) {
	directions := map[string]common.TestFunc{
		// "Pub": benchGoczmqPub,
		"Sub": benchGoczmqSub,
	}

	for dir, fn := range directions {
		for _, i := range []int{1, 10, 100, 1000} {
			b.Run(fmt.Sprintf("%s-%d", dir, i), fn(b, i))
		}
	}
}

func benchGoczmqPub(b *testing.B, num int) func(*testing.B) {
	port := fmt.Sprintf("%d", num+10000)

	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		ready := make(chan struct{}, num)

		// start subscribers
		for i := 0; i < num; i++ {
			go func() {
				// connect and subscribe to the publisher
				subSock, err := goczmq.NewSub("tcp://127.0.0.1:"+port, "")
				if err != nil {
					panic(err)
				}
				defer subSock.Destroy()

				if buf, n, err := subSock.RecvFrame(); err != nil {
					panic(fmt.Sprintf("start: %v %v %v", buf, n, err))
				}
				ready <- struct{}{}

				for {
					if err := ctx.Err(); err != nil {
						break
					}
					if buf, n, err := subSock.RecvFrame(); err != nil {
						panic(fmt.Sprintf("loop: %v %v %v", buf, n, err))
					}
				}
				ready <- struct{}{}
			}()
		}

		// start publisher
		pubSock, err := goczmq.NewPub("tcp://*:" + port)
		if err != nil {
			b.Fatal(err)
		}
		defer pubSock.Destroy()

		// make sure that all subscribers are ready
		for i := 0; i < num; {
			if err := pubSock.SendFrame(common.Msg, 0); err != nil {
				b.Fatal(err)
			}

			select {
			case <-ready:
				i++
			default:
			}
		}

		defer func() {
			cancel()
			// make sure that all subscribers are logged out
			for i := 0; i < num; {
				if err := pubSock.SendFrame(common.Msg, 0); err != nil {
					b.Fatal(err)
				}

				select {
				case <-ready:
					i++
				default:
				}
			}
		}()

		b.ResetTimer()

		sum := uint64(0)
		for i := 0; i < b.N; i++ {
			if err := pubSock.SendFrame(common.Msg, 0); err != nil {
				b.Fatal(err)
			}
			sum += uint64(len(common.Msg))
		}
		b.ReportMetric(float64(sum)/float64(b.N), "bytes/op")
	}
}

func benchGoczmqSub(b *testing.B, num int) func(*testing.B) {
	port := fmt.Sprintf("%d", num+10000)

	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		ready := make(chan struct{}, num)

		go func() {
			pubSock := goczmq.NewSock(goczmq.Pub)
			defer pubSock.Destroy()
			if _, err := pubSock.Bind("tcp://*:" + port); err != nil {
				panic(err)
			}

			// make sure that all subscribers are ready
			for i := 0; i < num; {
				if err := pubSock.SendFrame(common.Msg, 0); err != nil {
					panic(err)
				}

				select {
				case <-ready:
					i++
				default:
				}
			}

			for {
				if err := ctx.Err(); err != nil {
					break
				}
				if err := pubSock.SendFrame(common.Msg, 0); err != nil {
					panic(err)
				}
			}

			// make sure that all subscribers are logged out
			for i := 0; i < num-1; {
				if err := pubSock.SendFrame(common.Msg, 0); err != nil {
					panic(err)
				}

				select {
				case <-ready:
					i++
				default:
				}
			}
		}()

		for i := 0; i < num-1; i++ {
			go func() {
				subSock, err := goczmq.NewSub("tcp://127.0.0.1:"+port, "")
				if err != nil {
					panic(err)
				}
				defer subSock.Destroy()

				if buf, n, err := subSock.RecvFrame(); err != nil {
					panic(fmt.Sprintf("start: %v %v %v", buf, n, err))
				}
				ready <- struct{}{}

				for {
					if err := ctx.Err(); err != nil {
						break
					}
					if buf, n, err := subSock.RecvFrame(); err != nil {
						panic(fmt.Sprintf("loop: %v %v %v", buf, n, err))
					}
				}
				ready <- struct{}{}
			}()
		}

		sock, err := goczmq.NewSub("tcp://127.0.0.1:"+port, "")
		if err != nil {
			panic(err)
		}
		defer sock.Destroy()

		b.ResetTimer()
		sum := 0
		for i := 0; i < b.N; i++ {
			buf, _, err := sock.RecvFrame()
			if err != nil {
				b.Fatal(err)
			}
			sum += len(buf)
		}
		b.ReportMetric(float64(sum), "bytes/op")

		cancel()
		ready <- struct{}{}
	}
}
