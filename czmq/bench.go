package czmq

import (
	"context"
	"testing"

	czmq "github.com/zeromq/goczmq"
)

var msg = []byte("12345678901234567890123456789012345678901234567890")

func Bench(b *testing.B, socket string, recv bool) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch := czmq.NewSock(czmq.Pull)
		defer ch.Destroy()
		_, _ = ch.Bind(socket)

		go func(ctx context.Context) {
			ch := czmq.NewSock(czmq.Push)
			defer ch.Destroy()
			_ = ch.Connect(socket)

			for {
				if err := ctx.Err(); err != nil {
					return
				}
				if recv {
					_ = ch.SendFrame(msg, 0)
				} else {
					_, _, _ = ch.RecvFrame()
				}
			}
		}(ctx)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if recv {
				_, _, _ = ch.RecvFrame()
			} else {
				_ = ch.SendFrame(msg, 0)
			}
		}
	}
}
