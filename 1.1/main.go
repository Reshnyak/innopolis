package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	// Установка флагов
	filename := flag.String("file", "problems.csv", "Quiz CSV file'")
	isRand := flag.Bool("random", false, "random order")
	partSize := flag.Int("part", 100, "reading the file in parts of N lines")

	flag.Parse()

	totalQuestions, correctAnswers, err := ProcessCSV(*filename, *isRand, *partSize)
	if err != nil {
		log.Fatal(err)
	}
	incorrectAnswers := totalQuestions - correctAnswers

	fmt.Printf("Correct answers: %d\t", correctAnswers)
	fmt.Printf("Incorrect answers: %d\n", incorrectAnswers)
}
