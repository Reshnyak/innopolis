package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Работа с CSV файлом. Открытие и чтение пачками на тот случай, если файл большого размера
func ProcessCSV(filename string, isRand bool, partSize int) (int64, int64) {
	var totalQuestions, correctAnswers int64
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var lines [][]string
	fmt.Println("Please answer the questions:")
	for {
		// Читаем файл построчно
		line, err := reader.Read()
		if err != nil {
			break
		}
		lines = append(lines, line)

		//Обрабатываем пачку строк как наберется нужное количество
		if len(lines) >= partSize {

			totalQuestions += int64(len(lines))
			correctAnswers += processLines(lines, isRand)

			lines = nil
		}
	}

	// Обработка оставшихся строк
	if len(lines) > 0 {
		totalQuestions += int64(len(lines))
		correctAnswers += processLines(lines, isRand)
	}

	return totalQuestions, correctAnswers
}
func processLines(lines [][]string, isRand bool) int64 {
	var correctAnswers int64
	if isRand {
		rand.Shuffle(len(lines), func(i, j int) {
			lines[i], lines[j] = lines[j], lines[i]
		})
	}
	for _, line := range lines {
		question := line[0]
		answer := strings.TrimSpace(line[1]) // убираем unicode.IsSpace
		answer = strings.ToLower(answer)     //все в нижний, с ответом пользователя будем делать также

		fmt.Printf("%s = ", question)

		var userAnswer string
		fmt.Scanln(&userAnswer)
		userAnswer = strings.TrimSpace(userAnswer)
		userAnswer = strings.ToLower(userAnswer) //все в нижний
		if answer == userAnswer {                // вообще как вариант можно strings.EqualFold и не переводить всех в нижний
			correctAnswers++
		}
	}

	return correctAnswers
}

/*	*/
