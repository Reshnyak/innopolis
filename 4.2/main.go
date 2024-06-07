package main

import (
	"fmt"
	"math/big"
)

func main() {
	inputs := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 61, 14, 15, 16, 17, 18, 19, 23, 29, 30, 31}

	primeChan := make(chan int64)
	primes := make([]int64, 0, len(inputs))

	compositeChan := make(chan int64)
	composites := make([]int64, 0, len(inputs))

	snChan := SelectNums(inputs, primeChan, compositeChan)
	outPrime := AddNums(primeChan, &primes)
	outComposite := AddNums(compositeChan, &composites)
	<-snChan
	<-outPrime
	fmt.Printf("Primes:%v\n", primes)
	<-outComposite
	fmt.Printf("Composites:%v\n", composites)
}

// Разделяет массив чисел на два канала: простых и составных чисел
func SelectNums(nums []int64, prime, composite chan<- int64) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, num := range nums {
			if big.NewInt(num).ProbablyPrime(0) {
				prime <- num
				continue
			}
			composite <- num
		}
		close(prime)
		close(composite)
	}()
	return done
}

// Добавляет числа из канала в срез
func AddNums(numsChan <-chan int64, nums *[]int64) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for num := range numsChan {
			*nums = append(*nums, num)
		}
	}()
	return done
}
