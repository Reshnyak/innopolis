package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func init() {
	getSlice()
	fmt.Println("_______________________________________")
}

const (
	maxValue     = 999
	maxSliceSize = 100
	minSliceSize = 10
)

var arr []int

func getSlice() []int {
	if len(arr) == 0 {
		count := minSliceSize + rand.Int()%(maxSliceSize-minSliceSize)
		slice := make([]int, count)
		for i, _ := range slice {
			slice[i] = rand.Int() % maxValue
		}
		arr = slice
	}
	return arr
}

func BenchmarkConcurrencyFindSubSort(b *testing.B) {
	// slice := []int{101, 89, 39, 565, 2, 5, 7, 22, 45, 34, 67, 6}
	for i := 0; i < b.N; i++ {
		FindMinSubFromElemsBySort(getSlice())
	}
}

func BenchmarkNoConcurrencyFindSub(b *testing.B) {
	// slice := []int{101, 89, 39, 565, 2, 5, 7, 22, 45, 34, 67, 6}
	for i := 0; i < b.N; i++ {
		FindMinSubFromElemsNoConcorrency(getSlice())
	}
}

func BenchmarkConcurrencyFindSub(b *testing.B) {
	// slice := []int{101, 89, 39, 565, 2, 5, 7, 22, 45, 34, 67, 6}
	for i := 0; i < b.N; i++ {
		FindMinSubFromElemsConcurrency(getSlice())
	}
}
