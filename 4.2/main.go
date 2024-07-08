package main

import (
	"fmt"
	"math/big"
)

func main() {
	inputs := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 61, 14, 15, 16, 17, 18, 19, 23, 29, 30, 31}
	primes := make([]int64, 0, len(inputs))
	composites := make([]int64, 0, len(inputs))
	primeChan, compositeChan := SelectNums(inputs)

	outPrime := AddNums(primeChan, &primes)
	outComposite := AddNums(compositeChan, &composites)

	<-outPrime
	fmt.Printf("Primes:%v\n", primes)
	<-outComposite
	fmt.Printf("Composites:%v\n", composites)

}

// Разделяет массив чисел на два канала: простых и составных чисел
func SelectNums(nums []int64) (<-chan int64, <-chan int64) {
	prime := make(chan int64)
	composite := make(chan int64)
	go func() {
		defer close(prime)
		defer close(composite)
		for _, num := range nums {
			if big.NewInt(num).ProbablyPrime(0) {
				prime <- num
				continue
			}
			composite <- num
		}
	}()
	return prime, composite
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
