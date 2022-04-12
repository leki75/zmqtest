package main

import (
	"context"
	"fmt"
	"testing"
)

var (
	channelBuffer  = 1000
	channelMessage = []byte("12345678901234567890123456789012345678901234567890")
)

func BenchmarkGoChannel(b *testing.B) {
	directions := map[string]func(*testing.B, int) func(*testing.B){
		"Pub": benchChannelPub,
		"Sub": benchChannelSub,
	}
	for dir, fn := range directions {
		for _, n := range []int{16, 32, 64, 128, 256, 512} {
			b.Run(fmt.Sprintf("%s%d", dir, n), fn(b, n))
		}
	}
}

func benchChannelSub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		subscribers := make([]chan []byte, 0, num)

		for i := 0; i < num-1; i++ { // start num - 1 subscribers
			ch := make(chan []byte, channelBuffer)

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

			subscribers = append(subscribers, ch)
		}

		sub := make(chan []byte, channelBuffer)
		subscribers = append(subscribers, sub)

		go func(ctx context.Context) { // start publisher
			for {
				for j := 0; j < num; j++ {
					select {
					case <-ctx.Done():
						return
					case subscribers[j] <- channelMessage:
					default:
						// TODO: should we remove default or or not?
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
			ch := make(chan []byte, channelBuffer)
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
				select {
				case pub[j] <- channelMessage:
				default:
				}
			}
		}
	}
}
