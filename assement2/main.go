package main

import (
	"fmt"
	"log"
	"math"
	"sync"
)

const workersCount int = 4

func EvalSequence(matrix [][]int, userAnswer []int) int {
	if err := validation(matrix, userAnswer); err != nil {
		log.Println(err)
		return 0
	}
	//если в матрице очень большие веса, то при вычислении макс балла будет переполнение с уходом в отрицательные числа
	//в условии выбора максимума отрицательное отбросится, но это может быть не тот результат на который расчитываем
	maxGrade := calMaxGrade(matrix)
	if maxGrade == math.MaxInt {
		log.Println("maxGrade overflow ")
	}
	userGrade := calcUserGrade(matrix, userAnswer)
	if maxGrade == 0 {
		return 0
	}
	return userGrade * 100 / maxGrade
}

// Т.к. нам надо просчитать пути для каждой вершины, а задачи эти не взаимосвязаны
// и при больших матрицах могут быть затратны, то выполним их конкурентно -  по бэнчмарку выигрываем =)
func calMaxGrade(matrix [][]int) int {
	maxGrade := 0
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	gmx := &GraphMtx2{
		adjMatrix: matrix,
	}
	inputs := make(chan int)
	grades := make([]<-chan int, workersCount)
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		grades[i] = gmx.findMaxGradeByStartNode(inputs)
	}

	for _, grade := range grades {
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				mu.Lock()
				if maxGrade < val {
					maxGrade = val
				}
				mu.Unlock()
			}
		}(grade)
	}
	for i := 0; i < len(matrix); i++ {
		inputs <- i
	}
	close(inputs)
	wg.Wait()
	return maxGrade
}

func validation(matrix [][]int, userAnswer []int) error {
	if len(matrix) == 0 {
		return fmt.Errorf("matrix is empty")
	}
	for i := 0; i < len(matrix); i++ {
		if len(matrix[i]) != len(matrix) {
			return fmt.Errorf("the matrix is not square")
		}
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] < 0 {
				return fmt.Errorf("the weight of the matrix should not be negative")
			}
			if i == j && matrix[i][j] != 0 {
				return fmt.Errorf("there should be no loops in the matrix")
			}
		}
		answers := make(map[int]struct{})
		for _, answer := range userAnswer {
			if answer < 0 || answer > len(matrix)-1 {
				return fmt.Errorf("the answers values should not exceed the range of the matrix")
			}
			if _, ok := answers[answer]; ok {
				return fmt.Errorf("answer not unique")
			} else {
				answers[answer] = struct{}{}
			}

		}
	}
	return nil
}
func calcUserGrade(matrix [][]int, userAnswer []int) int {
	userGrade := 0
	for i := 1; i < len(userAnswer); i++ {
		userGrade += matrix[userAnswer[i-1]][userAnswer[i]]
	}
	return userGrade
}
func main() {
	mtx1 := [][]int{

		{0, 1, 0, 0, 0, 0, 0},
		{1, 0, 1, 1, 2, 0, 0},
		{0, 1, 0, 0, 0, 1, 1},
		{0, 1, 0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0, 0, 1},
		{0, 0, 1, 1, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 0},
	}

	gMtx := GraphMtx{
		adjMatrix: mtx1,
	}
	fmt.Printf("<<<%d>>>\n", EvalSequence(mtx1, []int{0, 1, 4, 6, 2, 5, 3}))
	fmt.Printf("MaxGrade:%d", gMtx.CalcMaxGrade())
	fmt.Println("\n___________________________________")

	fmt.Printf("MaxGrade:%d", calMaxGrade(mtx1))
}
