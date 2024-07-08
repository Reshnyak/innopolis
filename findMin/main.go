package main

import (
	"fmt"
	"math/rand"
	"sort"
)

const threshold = 4000

func MinElBySort(a []int) int {
	if len(a) == 0 {
		return 0
	}
	if len(a) == 1 {
		return a[0]
	}
	sort.Ints(a)
	return a[0]
}
func MinElConcurrent(a []int) int {
	resultChan := make(chan int)
	go minEl(a, resultChan)
	return <-resultChan
}
func minEl(a []int, resultChan chan<- int) {
	if len(a) == 0 {
		resultChan <- 0
		return
	}
	if len(a) == 1 {
		resultChan <- a[0]
		return
	}
	leftChan := make(chan int)
	rightChan := make(chan int)

	go minEl(a[:len(a)/2], leftChan)
	go minEl(a[len(a)/2:], rightChan)

	leftMin := <-leftChan
	rightMin := <-rightChan

	if leftMin <= rightMin {
		resultChan <- leftMin
	} else {
		resultChan <- rightMin
	}
}

func MinThresHoldParallel(thresHold int, a []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		if len(a) < thresHold {
			out <- MinEl2(a)
			return
		}

		leftChan := MinThresHoldParallel(thresHold, a[len(a)/2:])
		rightChan := MinThresHoldParallel(thresHold, a[:len(a)/2])

		leftMin := <-leftChan
		rightMin := <-rightChan

		if leftMin <= rightMin {
			out <- leftMin
		} else {
			out <- rightMin
		}
	}()

	return out
}

func main() {
	arr := rand.Perm(20000)
	min := <-MinThresHoldParallel(threshold, arr)
	fmt.Println(min) // Output: 1
}
func MinEl2(a []int) int {
	// только для первичной проверки
	if len(a) == 0 {
		return 0
	}
	if len(a) == 1 {
		return a[0]
	}
	t1 := MinEl2(a[:len(a)/2])
	t2 := MinEl2(a[len(a)/2:])
	if t1 <= t2 {
		return t1
	}
	return t2
}
