package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var students Cache[int, Student]
var objects Cache[int, Object]

func main() {
	jsonFile, err := os.Open("dz3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	fmt.Println("File descriptor successfully created")

	var CS ContolSection
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &CS)
	students.Init()
	objects.Init()
	for _, v := range CS.Students {
		students.Set(v.ID, v)
	}
	for _, v := range CS.Objects {
		objects.Set(v.ID, v)
	}

	fmt.Println("___________________________________________________")
	fmt.Printf("Student name \t| Grade\t| Object\t| Resulte |\n")
	fmt.Println("___________________________________________________")
	for _, r := range CS.Results {

		stdPtr, _ := students.Get(r.StudentID)
		objPtr, _ := objects.Get(r.ObjectID)

		fmt.Printf("%s \t\t| ", stdPtr.Name)
		fmt.Printf(" %d\t| ", stdPtr.Grade)
		if len(objPtr.Name) < 5 {
			fmt.Printf(" %s\t\t| ", objPtr.Name)
		} else {
			fmt.Printf(" %s\t| ", objPtr.Name)
		}
		fmt.Printf(" %d\t  |\n", r.Result)
	}
}
