package main

import "testing"

func Benchmark_generatePrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generatePrimes()
	}
}
