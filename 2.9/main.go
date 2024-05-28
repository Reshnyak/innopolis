package main

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}
type Numbers[T Number] []T

// суммирование всех элементов
func (num Numbers[T]) accumulate() T {
	var res T
	for _, val := range num {
		res += val
	}
	return res
}

// произведение всех элементов
func (num Numbers[T]) product() T {
	var res T
	if len(num) > 0 {
		res = num[0]
		for _, val := range num[1:] {
			res += val
		}
	}
	return res
}

// сравнение с другим слайсом на равность
func (num Numbers[T]) Equal(other Numbers[T]) bool {

	if len(num) != len(other) {
		return false
	}
	for i, valNum := range num {
		if valNum != other[i] {
			return false
		}
	}
	return true
}

// проверка аргумента, является ли он элементом массива,
// выводит индекс и результат поиска
func (num Numbers[T]) FirstIndexOF(value T) (int, bool) {

	for ind, val := range num {
		if val == value {
			return ind, true
		}
	}
	return 0, false
}

// удаление элемента массива по значению
func (num Numbers[T]) DeleteByValue(value T) Numbers[T] {
	for ind, val := range num {
		if val == value {
			return append(num[:ind], num[ind+1:]...)
		}
	}
	return num
}

// удаление элемента массива по индексу
func (num Numbers[T]) DeleteByIndex(ind int) Numbers[T] {
	return append(num[:ind], num[ind+1:]...)
}
