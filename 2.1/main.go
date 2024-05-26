package main

import (
	"fmt"
	"sort"
)

// находит пересечение неопределенного количества слайсов типа int
func sliceIntersection(slices ...[]int) []int {
	var res []int

	// Создаем мапу для подсчета общего количества уникальных чисел
	countMap := make(map[int]int)
	for _, slice := range slices {
		// Создаем мапу уникальных чисел в слайсе
		uniqueMap := make(map[int]bool)
		for _, num := range slice {
			if !uniqueMap[num] {
				uniqueMap[num] = true
			}
		}
		//Добавляем в общую мапу
		for num := range uniqueMap {
			countMap[num]++
		}
	}

	for num, count := range countMap {
		// Если ключ в маппе встречался во всех слайсах то value равно длине
		if count == len(slices) {
			res = append(res, num)
		}
	}
	// Сортируем
	sort.Ints(res)

	return res
}

func main() {
	fmt.Println(sliceIntersection([]int{1, 2}, []int{3, 2}))
	fmt.Println(sliceIntersection([]int{1, 2, 3, 2}, []int{3, 2}))
	fmt.Println(sliceIntersection([]int{1, 2, 3, 2}, []int{3, 2}, []int{}))
}
