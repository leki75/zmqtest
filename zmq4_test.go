package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/leki75/zmqtest/common"
	"github.com/pebbe/zmq4"
)

func BenchmarkZmq4(b *testing.B) {
	directions := map[string]common.TestFunc{
		"Pub": benchZmq4Pub,
		"Sub": benchZmq4Sub,
	}
	for dir, fn := range directions {
		for _, i := range []int{1, 10, 100, 1000} {
			b.Run(fmt.Sprintf("%s-%d", dir, i), fn(b, i))
		}
	}
}

func benchZmq4Pub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch, err := zmq4.NewSocket(zmq4.PUB)
		if err != nil {
			b.Fatal(err)
		}
		defer ch.Close()
		if err := ch.Bind(fmt.Sprintf("tcp://*:%d", num+5555)); err != nil {
			b.Fatal(err)
		}

		for i := 0; i < num; i++ {
			go func() {
				ch, err := zmq4.NewSocket(zmq4.SUB)
				if err != nil {
					panic(err)
				}
				defer ch.Close()
				if err := ch.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
					panic(err)
				}

				for {
					if err := ctx.Err(); err != nil {
						return
					}
					if _, err := ch.RecvBytes(0); err != nil {
						panic(err)
					}
				}
			}()
		}

		sum := uint64(0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			n, err := ch.SendBytes(common.Msg, 0)
			if err != nil {
				b.Fatal(err)
			}
			sum += uint64(n)
		}
		b.ReportMetric(float64(sum), "bytes/op")
	}
}

func benchZmq4Sub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			ch, err := zmq4.NewSocket(zmq4.PUB)
			if err != nil {
				panic(err)
			}
			defer ch.Close()
			if err := ch.Bind(fmt.Sprintf("tcp://*:%d", num+5555)); err != nil {
				panic(err)
			}

			for {
				if err := ctx.Err(); err != nil {
					return
				}
				if _, err = ch.SendBytes(common.Msg, 0); err != nil {
					panic(err)
				}
			}
		}()

		for i := 0; i < num-1; i++ {
			go func() {
				ch, err := zmq4.NewSocket(zmq4.SUB)
				if err != nil {
					panic(err)
				}
				defer ch.Close()
				if err := ch.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
					panic(err)
				}

				for {
					if err := ctx.Err(); err != nil {
						return
					}
					if _, err := ch.RecvBytes(0); err != nil {
						panic(err)
					}
				}
			}()
		}

		ch, err := zmq4.NewSocket(zmq4.SUB)
		if err != nil {
			b.Fatal(err)
		}
		defer ch.Close()
		if err := ch.Connect(fmt.Sprintf("tcp://127.0.0.1:%d", num+5555)); err != nil {
			b.Fatal(err)
		}

		sum := uint64(0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf, err := ch.RecvBytes(0)
			if err != nil {
				b.Fatal(err)
			}
			sum += uint64(len(buf))
		}
		b.ReportMetric(float64(sum), "bytes/op")
	}
}
