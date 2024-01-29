package main

import "testing"

func BenchmarkCountingLines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = start()
	}
}
