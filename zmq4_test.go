package main

import (
	"context"
	"testing"

	zmq "github.com/pebbe/zmq4"
)

var msgString = "12345678901234567890123456789012345678901234567890"

func BenchmarkZmq4Socket_recv_inproc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PULL)
	defer ch.Close()
	_ = ch.Bind("inproc://zmq4-recv.inproc")

	go func(ctx context.Context) {
		ch, _ := zmq.NewSocket(zmq.PUSH)
		defer ch.Close()
		_ = ch.Connect("inproc://zmq4-recv.inproc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Send(msgString, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Recv(0)
	}
}

func BenchmarkZmq4Socket_send_inproc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PUSH)
	defer ch.Close()
	_ = ch.Bind("inproc://zmq4-send.inproc")

	go func(ctx context.Context) {
		ch, _ := zmq.NewSocket(zmq.PULL)
		defer ch.Close()
		_ = ch.Connect("inproc://zmq4-send.inproc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Recv(0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Send(msgString, 0)
	}
}

func BenchmarkZmq4Socket_recv_ipc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PULL)
	defer ch.Close()
	_ = ch.Bind("ipc://zmq4-recv.ipc")

	go func(ctx context.Context) {
		ch, _ := zmq.NewSocket(zmq.PUSH)
		defer ch.Close()
		_ = ch.Connect("ipc://zmq4-recv.ipc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Send(msgString, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Recv(0)
	}
}

func BenchmarkZmq4Socket_send_ipc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PUSH)
	defer ch.Close()
	_ = ch.Bind("ipc://zmq4-send.ipc")

	go func() {
		ch, _ := zmq.NewSocket(zmq.PULL)
		defer ch.Close()
		_ = ch.Connect("ipc://zmq4-send.ipc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Recv(0)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Send(msgString, 0)
	}
}

func BenchmarkZmq4Socket_recv_tcp(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PULL)
	defer ch.Close()
	_ = ch.Bind("tcp://127.0.0.1:5555")

	go func(ctx context.Context) {
		ch, _ := zmq.NewSocket(zmq.PUSH)
		defer ch.Close()
		_ = ch.Connect("tcp://127.0.0.1:5555")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Send(msgString, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Recv(0)
	}
}

func BenchmarkZmq4Socket_send_tcp(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, _ := zmq.NewSocket(zmq.PUSH)
	defer ch.Close()
	_ = ch.Bind("tcp://127.0.0.1:5556")

	go func(ctx context.Context) {
		ch, _ := zmq.NewSocket(zmq.PULL)
		defer ch.Close()
		_ = ch.Connect("tcp://127.0.0.1:5556")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _ = ch.Recv(0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ch.Send(msgString, 0)
	}
}
