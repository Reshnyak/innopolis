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
type ContolSection struct {
	Students []Student
	Objects  []Object
	Results  []Result
}

// Метод возвращающий имя студента по его id
func (cs ContolSection) GetStudentNameByID(id int) string {
	for _, s := range cs.Students {
		if s.ID == id {
			return s.Name
		}
	}
	return ""
}

// Метод возвращающий grade студента по его id
func (cs ContolSection) GetStudentGradeByID(id int) int {
	for _, s := range cs.Students {
		if s.ID == id {
			return s.Grade
		}
	}
	return 0
}

// Метод возвращающий наименование предмета по id
func (cs ContolSection) GetObjectNameByID(id int) string {
	for _, obj := range cs.Objects {
		if obj.ID == id {
			return obj.Name
		}
	}
	return ""
}

// Метод возвращающий мапу средних значений результатов для каждого грэйда по id предмета
func (cs ContolSection) GetObjectGradeMeanById(objId int) map[int]float32 {
	type MeanGrade struct{ sum, count int }
	gradeMap := make(map[int]MeanGrade)
	for _, res := range cs.Results {
		if res.ObjectID == objId {
			grade := cs.GetStudentGradeByID(res.StudentID)
			gradeMap[grade] = MeanGrade{gradeMap[grade].sum + res.Result, gradeMap[grade].count + 1}
		}
	}
	resMap := make(map[int]float32, len(gradeMap))
	for k, v := range gradeMap {
		resMap[k] = float32(v.sum) / float32(v.count)
	}
	return resMap
}

// Форматированный вывод =) ...надо попробовать повторить c template
func (cs ContolSection) PrintControlSection() {
	fmt.Println("___________________________________________________")
	fmt.Printf("Student name \t| Grade\t| Object\t| Resulte |\n")
	fmt.Println("___________________________________________________")
	for _, r := range cs.Results {
		fmt.Printf("%s \t\t| ", cs.GetStudentNameByID(r.StudentID))
		fmt.Printf(" %d\t| ", cs.GetStudentGradeByID(r.StudentID))
		objName := cs.GetObjectNameByID(r.ObjectID)
		if len(objName) < 5 {
			fmt.Printf(" %s\t\t| ", cs.GetObjectNameByID(r.ObjectID))
		} else {
			fmt.Printf(" %s\t| ", cs.GetObjectNameByID(r.ObjectID))
		}
		fmt.Printf(" %d\t  |\n", r.Result)
	}

}

// Форматированный вывод сводных данных по предметам
func (cs ContolSection) PrintMeanObjects() {
	var total float32
	for _, obj := range cs.Objects {
		fmt.Println("_________________________")
		if len(obj.Name) < 8 {
			fmt.Printf("%s\t \t | Mean\t|\n", obj.Name)
		} else {
			fmt.Printf("%s\t | Mean\t|\n", obj.Name)
		}
		fmt.Println("_________________________")
		gradeMeans := cs.GetObjectGradeMeanById(obj.ID)
		slice := make([]int, 0, len(gradeMeans))
		for k := range gradeMeans {
			slice = append(slice, k)
		}
		sort.Ints(slice)
		for _, k := range slice {
			fmt.Printf("%d grade \t | %.1f\t|\n", k, gradeMeans[k])
			total += gradeMeans[k]
		}
		fmt.Println("_________________________")
		fmt.Printf("mean \t\t | %d\t|\n", int(total)/len(slice))
		fmt.Println("_________________________")
	}
}
