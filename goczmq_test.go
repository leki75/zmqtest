package main

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/leki75/zmqtest/common"
	"github.com/zeromq/goczmq"
)

func BenchmarkGoczmq(b *testing.B) {
	directions := map[string]common.TestFunc{
		"Pub": benchGoczmqPub,
		"Sub": benchGoczmqSub,
	}

	for dir, fn := range directions {
		for _, i := range []int{1, 10, 100, 1000} {
			b.Run(fmt.Sprintf("%s-%d", dir, i), fn(b, i))
		}
	}
}

func benchGoczmqPub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pubSock := goczmq.NewSock(goczmq.Pub)
		defer pubSock.Destroy()
		if _, err := pubSock.Bind(fmt.Sprintf("tcp://*:%d", num+5555)); err != nil {
			b.Fatal(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(num)
		for i := 0; i < num; i++ {
			go func() {
				subSock := goczmq.NewSock(goczmq.Sub)
				defer subSock.Destroy()

				if err := subSock.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
					panic(err)
				}
				wg.Done()

				for {
					if err := ctx.Err(); err != nil {
						return
					}
					if _, err := subSock.RecvMessage(); err != nil {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()

		sum := uint64(0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := pubSock.SendFrame(common.Msg, 0); err != nil {
				b.Fatal(err)
			}
			sum += uint64(len(common.Msg))
		}
		b.ReportMetric(float64(sum), "bytes/op")
	}
}

func benchGoczmqSub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ready := make(chan struct{})
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			pubSock := goczmq.NewSock(goczmq.Pub)
			defer pubSock.Destroy()
			if _, err := pubSock.Bind(fmt.Sprintf("tcp://*:%d", num+5555)); err != nil {
				panic(err)
			}
			wg.Done()

			<-ready
			for {
				if err := ctx.Err(); err != nil {
					return
				}
				if err := pubSock.SendFrame(common.Msg, 0); err != nil {
					panic(err)
				}
			}
		}()
		wg.Wait()

		for i := 0; i < num-1; i++ {
			go func() {
				subSock := goczmq.NewSock(goczmq.Sub)
				defer subSock.Destroy()
				if err := subSock.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
					panic(err)
				}

				for {
					if err := ctx.Err(); err != nil {
						return
					}
					if _, _, err := subSock.RecvFrame(); err != nil {
						panic(err)
					}
				}
			}()
		}

		sock := goczmq.NewSock(goczmq.Sub)
		defer sock.Destroy()

		if err := sock.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
			b.Fatal(err)
		}

		sum := 0
		ready <- struct{}{}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf, _, err := sock.RecvFrame()
			if err != nil {
				b.Fatal(err)
			}
			sum += len(buf)
		}
		b.ReportMetric(float64(sum), "bytes/op")
	}
}
