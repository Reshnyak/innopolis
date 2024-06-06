package main

import (
	"fmt"
	"sort"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Object struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Result struct {
	ObjectID  int `json:"object_id"`
	StudentID int `json:"student_id"`
	Result    int `json:"result"`
}
type ControlSection struct {
	Students []Student
	Objects  []Object
	Results  []Result
	grades   []int
}

// Метод возвращающий мапу студентов с ID - key
func (cs ControlSection) GetStudentsMap() map[int]Student {
	res := make(map[int]Student)
	for _, stud := range cs.Students {
		res[stud.ID] = stud
	}
	return res
}

// Метод возвращающий мапу предметов с ID - key
func (cs ControlSection) GetObjectsMap() map[int]Object {
	res := make(map[int]Object)
	for _, obj := range cs.Objects {
		res[obj.ID] = obj
	}
	return res
}

// Метод возвращающий мапу средних результатов по предметам и грэйдам
// Ключи первого уровня ObjectID,  второго грэйды
func (cs ControlSection) GetObjectGradesMeanMap() map[int]map[int]float32 {
	studMap := cs.GetStudentsMap()
	resMap := make(map[int]map[int]float32)
	for _, obj := range cs.Objects {
		resMap[obj.ID] = make(map[int]float32)
		for _, grade := range cs.GetSortGrades() {
			resGreade := Filter(cs.Results, func(res Result) bool {
				return res.ObjectID == obj.ID && studMap[res.StudentID].Grade == grade
			})
			sum := Reduce(resGreade, 0.0, func(res Result, b float32) float32 {
				return float32(res.Result) + b
			})
			if len(resGreade) > 0 {
				resMap[obj.ID][grade] = sum / float32(len(resGreade))
			}
		}
	}
	return resMap
}

// Возвращает срез отсортированных в порядке возрастания грэйдов
func (cs *ControlSection) GetSortGrades() []int {
	if len(cs.grades) == 0 {
		gradeMap := make(map[int]struct{})
		for _, value := range cs.Students {
			if _, ok := gradeMap[value.Grade]; !ok {
				gradeMap[value.Grade] = struct{}{}
			}
		}
		cs.grades = make([]int, 0, len(gradeMap))
		for grade := range gradeMap {
			cs.grades = append(cs.grades, grade)
		}
		sort.Ints(cs.grades)
	}
	return cs.grades
}

// GetMeanObjectGrade
// Форматированный вывод сводных данных по предметам с использованием функций высшего порядка
func (cs ControlSection) PrintMeanObjectsByFunctions() {
	grades := cs.GetSortGrades()
	objectsGradeMeans := cs.GetObjectGradesMeanMap()
	for _, obj := range cs.Objects {
		var total float32
		fmt.Println("_________________________")
		if len(obj.Name) < 8 {
			fmt.Printf("%s\t \t | Mean\t|\n", obj.Name)
		} else {
			fmt.Printf("%s\t | Mean\t|\n", obj.Name)
		}
		fmt.Println("_________________________")
		for _, grade := range grades {
			mean := objectsGradeMeans[obj.ID][grade]
			fmt.Printf("%d grade \t | %.1f\t|\n", grade, mean)
			total += mean
		}
		fmt.Println("_________________________")
		fmt.Printf("mean \t\t | %d\t|\n", int(total)/len(grades))
		fmt.Println("_________________________")
	}
}
