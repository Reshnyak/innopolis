package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// ReadConsole - чтение из консоли и отправка в канал. Завершение происходит при завершении контекста
func ReadConsole(ctx context.Context, wg *sync.WaitGroup, rd io.Reader, out chan<- string) {
	go func() {
		defer close(out)
		defer wg.Done()
		reader := bufio.NewReader(rd)
		for {
			select {
			case <-ctx.Done():
				log.Println("Reading from console finished")
				break
			default:
				str, err := reader.ReadString('\n')
				if err != nil && err != io.EOF {
					log.Println(fmt.Errorf("ReadConsole() reader.ReadString: %s", err))
					break
				}
				out <- str
			}
		}
	}()
}

// WriteFile Записывает данные из канала в файл
func WriteFile(ctx context.Context, wg *sync.WaitGroup, fileName string, in <-chan string) {
	go func() {
		defer wg.Done()
		file, err := os.Create(fileName)
		if err != nil {
			log.Println(fmt.Errorf("WriteFile(%s) os.Create: %s", fileName, err))
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Writing to file finished")
				if err = file.Close(); err != nil {
					log.Println(fmt.Errorf("WriteFile(%s) file.Close: %s", fileName, err))
				}
				break
			default:
				for line := range in {
					_, err := file.WriteString(line)
					if err != nil {
						log.Println(fmt.Errorf("WriteFile(%s) file.WriteString: %s", fileName, err))
					}
				}
			}
		}

	}()
}
