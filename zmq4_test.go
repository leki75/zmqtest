package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/pebbe/zmq4"
)

var (
	zmqMessage = []byte("12345678901234567890123456789012345678901234567890")
	// zmqBuffer  = 1000 // the default socket buffer size

	zmqCtx *zmq4.Context
)

func init() {
	zmqCtx = newContext(2048)
	newProxy()
}

func BenchmarkZmq4(b *testing.B) {
	directions := map[string]func(*testing.B, int) func(*testing.B){
		"Pub": benchZmq4Pub,
		"Sub": benchZmq4Sub,
	}
	for dir, fn := range directions {
		for _, i := range []int{16, 32, 64, 128, 256, 512} {
			b.Run(fmt.Sprintf("%s%d", dir, i), fn(b, i))
		}
	}
}

func onError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func newContext(size int) *zmq4.Context {
	ctx, err := zmq4.NewContext()
	onError(err)
	onError(ctx.SetIoThreads(2))
	onError(ctx.SetMaxSockets(size))
	return ctx
}

func newProxy() {
	go func() {
		// publisher side
		pub, err := zmqCtx.NewSocket(zmq4.XSUB)
		onError(err)
		defer pub.Close()
		onError(pub.Bind("inproc://pubsub"))

		// subscriber side
		sub, err := zmqCtx.NewSocket(zmq4.XPUB)
		onError(err)
		defer sub.Close()
		onError(sub.Bind("tcp://*:10000"))

		// start proxy
		onError(zmq4.Proxy(pub, sub, nil))
	}()
}

func benchZmq4Pub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		// zmqCtx := newContext(num + 1)
		// defer onError(zmqCtx.Term())

		ctx, cancel := context.WithCancel(context.Background())

		ready := make(chan struct{}, num)
		done := make(chan struct{}, num)

		// subscribers
		for i := 0; i < num; i++ {
			go func() {
				sub, err := zmqCtx.NewSocket(zmq4.SUB)
				onError(err)
				defer sub.Close()

				onError(sub.Connect("tcp://127.0.0.1:10000"))
				sub.SetSubscribe("")

				// sync with publisher
				_, err = sub.RecvBytes(0)
				onError(err)
				ready <- struct{}{}

				// read messages
				for {
					if err := ctx.Err(); err != nil {
						break
					}
					_, err := sub.RecvBytes(0)
					onError(err)
				}

				// end
				done <- struct{}{}
			}()
		}

		// publisher
		pub, err := zmqCtx.NewSocket(zmq4.PUB)
		onError(err)
		defer pub.Close()

		onError(pub.Connect("inproc://pubsub"))
		// onError(pub.Bind("tcp://*:10000"))

		// wait for subscribers to start
		for i := 0; i < num; {
			_, err = pub.SendBytes(zmqMessage, 0)
			onError(err)
			select {
			case <-ready:
				i++
			default:
			}
		}

		// benchmark
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := pub.SendBytes(zmqMessage, 0)
			onError(err)
		}
		b.StopTimer()

		// wait for subscribers to stop
		cancel()
		for i := 0; i < num; {
			_, err = pub.SendBytes(zmqMessage, 0)
			onError(err)
			select {
			case <-done:
				i++
			default:
			}
		}
	}
}

func benchZmq4Sub(b *testing.B, num int) func(*testing.B) {
	return func(b *testing.B) {
		// zmqCtx := newContext(num + 1)
		// defer onError(zmqCtx.Term())

		ctx, cancel := context.WithCancel(context.Background())

		ready := make(chan struct{}, num)
		done := make(chan struct{}, num)

		// publisher
		go func() {
			pub, err := zmqCtx.NewSocket(zmq4.PUB)
			onError(err)
			defer pub.Close()

			onError(pub.Connect("inproc://pubsub"))
			// onError(pub.Bind("tcp://*:10000"))

			// wait for subscribers to start
			for i := 0; i < num; {
				_, err = pub.SendBytes(zmqMessage, 0)
				onError(err)
				select {
				case <-ready:
					i++
				default:
				}
			}

			// send messages
			for {
				if err := ctx.Err(); err != nil {
					break
				}
				_, err = pub.SendBytes(zmqMessage, 0)
				onError(err)
			}

			// waiting for subscribers to stop
			for i := 0; i < num-1; {
				_, err = pub.SendBytes(zmqMessage, 0)
				onError(err)
				select {
				case <-done:
					i++
				default:
				}
			}

			// end
			ready <- struct{}{}
		}()

		// subscribers
		for i := 0; i < num-1; i++ {
			go func() {
				sub, err := zmqCtx.NewSocket(zmq4.SUB)
				onError(err)
				defer sub.Close()

				onError(sub.Connect("tcp://127.0.0.1:10000"))
				sub.SetSubscribe("")

				// sync with publisher
				_, err = sub.RecvBytes(0)
				onError(err)
				ready <- struct{}{}

				// read messages
				for {
					if err := ctx.Err(); err != nil {
						break
					}
					_, err := sub.RecvBytes(0)
					onError(err)
				}

				// end
				done <- struct{}{}
			}()
		}

		// one more subscriber
		sub, err := zmqCtx.NewSocket(zmq4.SUB)
		onError(err)
		defer sub.Close()

		onError(sub.Connect("tcp://127.0.0.1:10000"))
		sub.SetSubscribe("")

		// sync with publisher
		_, err = sub.RecvBytes(0)
		onError(err)
		ready <- struct{}{}

		// benchmark
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := sub.RecvBytes(0)
			onError(err)
		}
		b.StopTimer()

		// wait for publisher to stop
		cancel()
		<-ready
	}
}
