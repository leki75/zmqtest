package main

import (
	"testing"

	"github.com/leki75/zmqtest/czmq"
)

func BenchmarkCzmqSocket(b *testing.B) {
	b.Run("Recv-Inproc", czmq.Bench(b, "inproc://czmq-recv.inproc", true))
	b.Run("Send-Inproc", czmq.Bench(b, "inproc://czmq-recv.inproc", false))
	b.Run("Recv-IPC", czmq.Bench(b, "ipc://czmq-recv.ipc", true))
	b.Run("Send-IPC", czmq.Bench(b, "ipc://czmq-recv.ipc", false))
	b.Run("Recv-TCP", czmq.Bench(b, "tcp://127.0.0.1:5557", true))
	b.Run("Send-TCP", czmq.Bench(b, "tcp://127.0.0.1:5558", false))
}
