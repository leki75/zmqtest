package main

import "testing"

func BenchmarkGoChannel_recv(b *testing.B) {
	ch := make(chan string)
	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case ch <- "12345678901234567890123456789012345678901234567890":
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkGoChannel_send(b *testing.B) {
	ch := make(chan string)
	defer close(ch)

	go func() {
		for range ch {
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- "12345678901234567890123456789012345678901234567890"
	}
}
