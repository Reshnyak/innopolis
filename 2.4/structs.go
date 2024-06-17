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
}

type MeanGrade struct {
	Sum   int
	Count int
	Mean  float32
}

func (mg *MeanGrade) CalcMean() {
	if mg.Count > 0 {
		mg.Mean = float32(mg.Sum) / float32(mg.Count)
	}
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

// Метод возвращающий мапу структур MeanGrade результатов по предметам и грэйдам
// Ключи первого уровня ObjectID,  второго грэйды
func (cs ControlSection) GetObjectGradesSum() map[int]map[int]MeanGrade {
	resMap := make(map[int]map[int]MeanGrade)
	for _, obj := range cs.Objects {
		resMap[obj.ID] = make(map[int]MeanGrade)
	}
	studMap := cs.GetStudentsMap()
	for _, res := range cs.Results {
		obj := res.ObjectID
		grade := studMap[res.StudentID].Grade
		//Исправлено - не в одну строку, если в одну то 0 обязателен.
		resMap[obj][grade] = MeanGrade{
			Sum:   resMap[obj][grade].Sum + res.Result,
			Count: resMap[obj][grade].Count + 1,
		}
	}
	for objKey, gradeMap := range resMap {
		for gradeKey, mg := range gradeMap {
			mg.CalcMean()
			resMap[objKey][gradeKey] = mg
		}
	}
	return resMap
}

// Форматированный вывод сводных данных по предметам
func (cs ControlSection) PrintMeanObjects() {
	objectsGradeSums := cs.GetObjectGradesSum()

	for _, obj := range cs.Objects {
		var total float32
		var count int
		fmt.Println("_________________________")
		if len(obj.Name) < 8 {
			fmt.Printf("%s\t \t | Mean\t|\n", obj.Name)
		} else {
			fmt.Printf("%s\t | Mean\t|\n", obj.Name)
		}
		fmt.Println("_________________________")
		grades := make([]int, 0, len(objectsGradeSums[obj.ID]))
		for grade := range objectsGradeSums[obj.ID] {
			grades = append(grades, grade)
		}
		sort.Ints(grades)
		for _, grade := range grades {
			fmt.Printf("%d grade \t | %.1f\t|\n", grade, objectsGradeSums[obj.ID][grade].Mean)
			total += float32(objectsGradeSums[obj.ID][grade].Sum)
			count += objectsGradeSums[obj.ID][grade].Count
		}
		if count > 0 {
			fmt.Println("_________________________")
			fmt.Printf("mean \t\t | %3.f\t|\n", total/float32(count))
			fmt.Println("_________________________")
		}
	}
}
