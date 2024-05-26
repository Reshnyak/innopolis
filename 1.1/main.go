package main

import (
	"flag"
	"fmt"
)

func main() {
	// Установка флагов
	filename := flag.String("file", "problems.csv", "Quiz CSV file'")
	isRand := flag.Bool("random", false, "random order")
	partSize := flag.Int("part", 100, "reading the file in parts of N lines")

	flag.Parse()

	totalQuestions, correctAnswers := ProcessCSV(*filename, *isRand, *partSize)
	incorrectAnswers := totalQuestions - correctAnswers

	fmt.Printf("Correct answers: %d\t", correctAnswers)
	fmt.Printf("Incorrect answers: %d\n", incorrectAnswers)
}
