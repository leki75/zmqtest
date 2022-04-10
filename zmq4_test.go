package main

import (
	"testing"

	"github.com/leki75/zmqtest/zmq4"
)

func BenchmarkZmq4Socket(b *testing.B) {
	b.Run("Recv-Inproc", zmq4.Bench(b, "inproc://zmq4-recv.inproc", true))
	b.Run("Send-Inproc", zmq4.Bench(b, "inproc://zmq4-recv.inproc", false))
	b.Run("Recv-IPC", zmq4.Bench(b, "ipc://zmq4-recv.ipc", true))
	b.Run("Send-IPC", zmq4.Bench(b, "ipc://zmq4-recv.ipc", false))
	b.Run("Recv-TCP", zmq4.Bench(b, "tcp://127.0.0.1:5555", true))
	b.Run("Send-TCP", zmq4.Bench(b, "tcp://127.0.0.1:5556", false))
}
