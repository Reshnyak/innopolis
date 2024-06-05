package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	jsonFile, err := os.Open("dz3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	fmt.Println("File descriptor successfully created")

	var CS ControlSection
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &CS)
	StudMap := CS.GetStudentsMap()
	ObjMap := CS.GetObjectsMap()
	fmt.Println("________________________________________________________________")
	FormatPrint(5, "Student name", "Grade", "Object", "Result")
	fmt.Println("________________________________________________________________")
	for _, res := range CS.Results {
		name := StudMap[res.StudentID].Name
		grade := strconv.Itoa(StudMap[res.StudentID].Grade)
		object := ObjMap[res.ObjectID].Name
		FormatPrint(7, name, grade, object, strconv.Itoa(res.Result))
	}
}

// Форматированный вывод
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
