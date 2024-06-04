package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// Читает построчно данные и отправляет в канал out
func ReadConsole(rd io.Reader, out chan<- string) {
	reader := bufio.NewReader(rd)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		out <- str
	}
}

// Записывает данные из канала в файл
func WriteFile(fileName string, in <-chan string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)

	}
	defer file.Close()
	for line := range in {
		file.WriteString(line)
	}

}
