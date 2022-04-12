package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/leki75/zmqtest/common"
)

func BenchmarkGoChannel(b *testing.B) {
	directions := map[string]common.TestFunc{
		"Pub": benchChannelPub,
		"Sub": benchChannelSub,
	}
	for dir, fn := range directions {
		for _, i := range []int{1, 10, 100, 1000} {
			b.Run(fmt.Sprintf("%s-%d", dir, i), fn(b, i))
		}
	}
}

func benchChannelSub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pub := make([]chan []byte, 0, num)
		sub := make(chan []byte)
		pub = append(pub, sub)

		for i := 0; i < num-1; i++ { // start num - 1 subscribers
			ch := make(chan []byte)
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case _, ok := <-ch:
						if !ok {
							return
						}
					}
				}
			}()

			pub = append(pub, ch)
		}

		go func(ctx context.Context) { // start publisher
			for {
				for j := 0; j < num; j++ {
					select {
					case <-ctx.Done():
						return
					case pub[j] <- common.Msg:
					}
				}
			}
		}(ctx)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			<-sub
		}
	}
}

func benchChannelPub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		pub := make([]chan []byte, 0, num)
		for i := 0; i < num; i++ {
			ch := make(chan []byte)
			defer close(ch)

			go func(ch chan []byte) {
				for range ch {
				}
			}(ch)

			pub = append(pub, ch)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < num; j++ {
				pub[j] <- common.Msg
			}
		}
	}
}
