package main

import (
	"math/rand"
	"testing"
)

const countElements = 20000

// var inputs = rand.Perm(countElements)

func BenchmarkMinEl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MinEl2(rand.Perm(countElements))
	}
}
func BenchmarkMinElBySort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MinElBySort(rand.Perm(countElements))
	}
}

func BenchmarkConcurencyMinEl(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MinElConcurrent(rand.Perm(countElements))
	}
}

func BenchmarkMinElThresHoldParallel(b *testing.B) {
	benchs := []struct {
		name      string
		threshold int
	}{
		{"1000", 1000},
		{"2000", 2000},
		{"4000", 4000},
		{"8000", 8000},
		{"10000", 10000},
		{"12000", 12000},
		{"16000", 16000},
		{"21000", 21000},
	}
	for _, bench := range benchs {
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				<-MinThresHoldParallel(bench.threshold, rand.Perm(countElements))
			}
		})
	}

}
