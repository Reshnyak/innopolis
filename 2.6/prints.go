package main

import (
	"fmt"
	"sort"
)

// GetMeanObjectGrade
// Форматированный вывод сводных данных по предметам с использованием функций высшего порядка
func (cs ControlSection) PrintMeanObjectsByFunctions() {

	for _, obj := range cs.Objects {
		var total float32
		fmt.Println("_________________________")
		if len(obj.Name) < 8 {
			fmt.Printf("%s\t \t | Mean\t|\n", obj.Name)
		} else {
			fmt.Printf("%s\t | Mean\t|\n", obj.Name)
		}
		fmt.Println("_________________________")

		slice := cs.GetGrades()

		for _, k := range slice {
			mean := cs.GetMeanObjectGrade(obj.ID, k)
			fmt.Printf("%d grade \t | %.1f\t|\n", k, mean)
			total += mean
		}
		fmt.Println("_________________________")
		fmt.Printf("mean \t\t | %d\t|\n", int(total)/len(slice))
		fmt.Println("_________________________")
	}
}

// Форматированный вывод сводных данных по предметам
func (cs ControlSection) PrintMeanObjects() {

	for _, obj := range cs.Objects {
		var total float32
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
