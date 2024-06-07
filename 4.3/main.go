package main

import (
	"fmt"
	"sync"
)

func mergeChans[T any](in1, in2 <-chan T, out chan T) <-chan struct{} {
	done := make(chan struct{})
	defer close(done)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for valIn := range in1 {
			out <- valIn
		}
	}()
	go func() {
		defer wg.Done()
		for valIn := range in2 {
			out <- valIn
		}
	}()
	wg.Wait()
	close(out)
	return done
}
func main() {
	nums1 := []int{1, 2, 3}
	nums2 := []int{4, 5, 6}
	results := make([]int, 0, len(nums1)+len(nums2))
	resChan := make(chan int)

	chan1 := sendAnyInChan(nums1...)
	chan2 := sendAnyInChan(nums2...)

	recv := recvAnyInChan(resChan, &results)
	<-mergeChans(chan1, chan2, resChan)

	<-recv
	fmt.Printf("Merge results:%v", results)

}

// Вспомогательная возможно опасная функция для отправки в канал
func sendAnyInChan[T any](args ...T) <-chan T {
	out := make(chan T)
	go func() {
		defer func() {
			close(out)
		}()
		for _, arg := range args {
			out <- arg
		}
	}()
	return out
}

// Вторая вспомогательная возможно опасная функция для приема из канала
func recvAnyInChan[T any](in chan T, results *[]T) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for arg := range in {
			*results = append(*results, arg)
		}
	}()
	return done
}
