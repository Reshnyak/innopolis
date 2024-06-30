package main

import (
	"fmt"
	"sync"
)

// Принимает на вход каналы и соединяет их содержимое в один
func mergeChans[T any](inChans ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(len(inChans))
		for _, ch := range inChans {
			go func(c <-chan T) {
				defer wg.Done()
				for valIn := range c {
					out <- valIn
				}
			}(ch)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {

	nums1 := []int{1, 2, 3}
	nums2 := []int{4, 5, 6}
	results := make([]int, 0, len(nums1)+len(nums2))

	chan1 := sendAnyInChan(nums1...)
	chan2 := sendAnyInChan(nums2...)
	resChan := mergeChans(chan1, chan2)

	recv := recvAnyInChan(resChan, &results)

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
func recvAnyInChan[T any](in <-chan T, results *[]T) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for arg := range in {
			*results = append(*results, arg)
		}
	}()
	return done
}
