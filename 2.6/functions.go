package main

import (
	"sort"
)

func Filter[T any](slice []T, f func(T) bool) []T {
	var res []T
	for _, v := range slice {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[T1, T2 any](slice []T1, f func(T1) T2) []T2 {
	var res []T2
	for _, v := range slice {
		res = append(res, f(v))
	}
	return res
}

func Reduce[T1, T2 any](slice []T1, init T2, f func(T1, T2) T2) T2 {
	res := init
	for _, value := range slice {
		res = f(value, res)
	}
	return res
}
func (cs ControlSection) GetGrades() []int {
	gradeMap := make(map[int]bool)
	for _, value := range cs.Students {
		if !gradeMap[value.Grade] {
			gradeMap[value.Grade] = true
		}
	}
	res := make([]int, 0, len(gradeMap))
	for k, _ := range gradeMap {
		res = append(res, k)
	}
	sort.Ints(res)
	return res
}

// Метод возвращающий  среднее значение результатов среза для каждого грэйда и id предмета,
// используя функции высшего порядка
func (cs ControlSection) GetMeanObjectGrade(objId int, grade int) float32 {
	resByObj := Filter(cs.Results, func(res Result) bool {
		return res.ObjectID == objId
	})
	resGreade := Filter(resByObj, func(res Result) bool {
		for _, value := range cs.Students {
			if res.StudentID == value.ID {
				return res.ObjectID == objId
			}
		}
		return false
	})
	mean := Reduce(resGreade, 0.0, func(res Result, b float32) float32 {
		return float32(res.Result) + b
	})
	return mean / float32(len(resGreade))

}
