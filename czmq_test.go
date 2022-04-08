package main

import (
	"context"
	"testing"

	zmq "github.com/zeromq/goczmq"
)

var msgBytes = []byte("12345678901234567890123456789012345678901234567890")

func BenchmarkCzmqSocket_recv_inproc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Pull)
	defer ch.Destroy()
	_, _ = ch.Bind("inproc://czmq-recv.inproc")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Push)
		defer ch.Destroy()
		_ = ch.Connect("inproc://czmq-recv.inproc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_ = ch.SendFrame(msgBytes, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ch.RecvFrame()
	}
}

func BenchmarkCzmqSocket_send_inproc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Push)
	defer ch.Destroy()
	_, _ = ch.Bind("inproc://czmq-send.inproc")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Pull)
		defer ch.Destroy()
		_ = ch.Connect("inproc://czmq-send.inproc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _, _ = ch.RecvFrame()
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ch.SendFrame(msgBytes, 0)
	}
}

func BenchmarkCzmqSocket_recv_ipc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Pull)
	defer ch.Destroy()
	_, _ = ch.Bind("ipc://czmq-recv.ipc")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Push)
		defer ch.Destroy()
		_ = ch.Connect("ipc://czmq-recv.ipc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_ = ch.SendFrame(msgBytes, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ch.RecvFrame()
	}
}

func BenchmarkCzmqSocket_send_ipc(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Push)
	defer ch.Destroy()
	_, _ = ch.Bind("ipc://czmq-send.ipc")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Pull)
		defer ch.Destroy()
		_ = ch.Connect("ipc://czmq-send.ipc")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _, _ = ch.RecvFrame()
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ch.SendFrame(msgBytes, 0)
	}
}

func BenchmarkCzmqSocket_recv_tcp(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Pull)
	defer ch.Destroy()
	_, _ = ch.Bind("tcp://127.0.0.1:5557")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Push)
		defer ch.Destroy()
		_ = ch.Connect("tcp://127.0.0.1:5557")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_ = ch.SendFrame(msgBytes, 0)
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ch.RecvFrame()
	}
}

func BenchmarkCzmqSocket_send_tcp(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := zmq.NewSock(zmq.Push)
	defer ch.Destroy()
	_, _ = ch.Bind("tcp://127.0.0.1:5558")

	go func(ctx context.Context) {
		ch := zmq.NewSock(zmq.Pull)
		defer ch.Destroy()
		_ = ch.Connect("tcp://127.0.0.1:5558")

		for {
			if err := ctx.Err(); err != nil {
				return
			}
			_, _, _ = ch.RecvFrame()
		}
	}(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ch.SendFrame(msgBytes, 0)
	}
}
