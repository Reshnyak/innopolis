package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	jsonFile, err := os.Open("dz3.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	fmt.Println("File descriptor successfully created")

	var CS ControlSection
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteValue, &CS)
	if err != nil {
		log.Fatal(err)
	}
	CS.PrintResult(CS.Results)
}
