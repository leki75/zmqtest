package zmq4

import (
	"context"
	"testing"

	zmq "github.com/pebbe/zmq4"
)

var msg = []byte("12345678901234567890123456789012345678901234567890")

func Bench(b *testing.B, socket string, recv bool) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var flag zmq.Type
		if recv {
			flag = zmq.PULL
		} else {
			flag = zmq.PUSH
		}

		ch, _ := zmq.NewSocket(flag)
		defer ch.Close()
		_ = ch.Bind(socket)

		go func() {
			var flag zmq.Type
			if recv {
				flag = zmq.PUSH
			} else {
				flag = zmq.PULL
			}

			ch, _ := zmq.NewSocket(flag)
			defer ch.Close()
			_ = ch.Connect(socket)

			for {
				if err := ctx.Err(); err != nil {
					return
				}
				if recv {
					_, _ = ch.SendBytes(msg, 0)
				} else {
					_, _ = ch.RecvBytes(0)
				}
			}
		}()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if recv {
				_, _ = ch.RecvBytes(0)
			} else {
				_, _ = ch.SendBytes(msg, 0)
			}
		}
	}
}
