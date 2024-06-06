package main

import (
	"fmt"
	"sort"
	"strconv"
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

func (cs ControlSection) FindTopStudent() []Result {
	tops := Filter(cs.Results, func(result Result) bool {
		stud := Filter(cs.Results, func(r Result) bool {
			return r.StudentID == result.StudentID
		})
		if sum := Reduce(stud, 0, func(a Result, b int) int {
			return a.Result + b
		}); sum == len(stud)*5 {
			return true
		}
		return false
	})
	return tops
}

// Форматированный вывод
func (cs ControlSection) PrintResult(results []Result) {
	StudMap := cs.GetStudentsMap()
	ObjMap := cs.GetObjectsMap()
	fmt.Println("________________________________________________________________")
	FormatPrint(5, "Student name", "Grade", "Object", "Result")
	fmt.Println("________________________________________________________________")
	for _, res := range results {
		name := StudMap[res.StudentID].Name
		grade := strconv.Itoa(StudMap[res.StudentID].Grade)
		object := ObjMap[res.ObjectID].Name
		FormatPrint(7, name, grade, object, strconv.Itoa(res.Result))
	}
}
func FormatPrint(constraint int, colums ...string) {
	if len(colums) > 0 {
		for _, colum := range colums {
			if len(colum) < constraint {
				fmt.Printf(" %s\t\t| ", colum)
				continue
			}
			fmt.Printf(" %s\t| ", colum)
		}
		fmt.Println()
	}
}
