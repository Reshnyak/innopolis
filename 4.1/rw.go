package main

//Функция ReadConsole реализована с позиции чтобы была возможность изменить redear на буфер или файл
//+ попытка вывести ошибку для её последующей обработки(тут вот пока не пришел к осознанию правильных практик)
import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Читает построчно данные и отправляет в канал out
// Функция ReadConsole реализована с позиции чтобы была возможность изменить redear на буфер или файл(например для тестов)
// + попытка вывести ошибку для её последующей обработки(тут вот пока не пришел к осознанию правильных практик)
func ReadConsole(rd io.Reader, out chan<- string, errChan chan<- error) {
	go func() {
		reader := bufio.NewReader(rd)
		for {
			str, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				errChan <- fmt.Errorf("ReadConsole() reader.ReadString: %s", err)
				break
			}
			out <- str
		}
	}()
}

// Записывает данные из канала в файл
func WriteFile(fileName string, in <-chan string, errChan chan<- error) {
	go func() {
		file, err := os.Create(fileName)
		if err != nil {
			errChan <- fmt.Errorf("WriteFile(%s) os.Create: %s", fileName, err)
		}

		for line := range in {
			_, err := file.WriteString(line)
			if err != nil {
				errChan <- fmt.Errorf("WriteFile(%s) file.WriteString: %s", fileName, err)
			}
		}
		errChan <- file.Close()
	}()
}
