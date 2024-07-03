package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

const (
	maxWorkers = 4
)

func main() {
	slice := []int{101, 89, 39, 565, 2, 5, 8, 22, 45, 34, 67, 6}
	fmt.Println(FindMinSubFromElemsBySort(slice))
}

// Быстро сортируем и проходим с двух сторон ищем минимальную разницу
func FindMinSubFromElemsBySort(elems []int) int {
	var minDist int
	if len(elems) > 1 {
		getDist := func(a, b int) int {
			dist := a - b
			if dist < 0 {
				dist *= -1
			}
			return dist
		}

		sort.Ints(elems)

		minDist = getDist(elems[0], elems[1])
		for l, r := 2, len(elems)-1; l < r; l, r = l+1, r-1 {
			dist := getDist(elems[r-1], elems[r])
			if minDist > dist {
				minDist = dist
			}

			dist = getDist(elems[l-1], elems[l])
			if minDist > dist {
				minDist = dist
			}
		}
	}
	return int(minDist)
}

func worker(mu *sync.Mutex, wg *sync.WaitGroup, startPos <-chan int, sl []int, minDist *float64) {

	go func() {
		defer wg.Done()
		for pos := range startPos {
			for j := pos + 1; j < len(sl); j++ {
				dist := math.Abs(float64(sl[pos] - sl[j]))
				// fmt.Printf("[%d:%d][%d-%d] =  dist:%f\n", pos, j, sl[pos], sl[j], dist)
				mu.Lock()
				*minDist = math.Min(dist, *minDist)
				mu.Unlock()
			}
		}
	}()

}
func FindMinSubFromElemsConcurrency(elems []int) int {

	var minDist float64
	if len(elems) > 1 {
		inCh := make(chan int)
		mu := new(sync.Mutex)
		wg := new(sync.WaitGroup)
		minDist = math.Abs(float64(elems[0] - elems[1]))
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			worker(mu, wg, inCh, elems, &minDist)
		}
		for i, _ := range elems {
			inCh <- i
		}
		close(inCh)
		wg.Wait()
	}
	return int(minDist)
}
func FindMinSubFromElemsNoConcorrency(elems []int) int {

	var minDist float64
	if len(elems) > 1 {
		minDist = math.Abs(float64(elems[0] - elems[1]))
		for i := 0; i < len(elems); i++ {
			for j := i + 1; j < len(elems); j++ {
				dist := math.Abs(float64(elems[i] - elems[j]))
				minDist = math.Min(dist, minDist)
			}
		}
	}
	return int(minDist)
}
