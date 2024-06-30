package main

//Функция ReadConsole реализована с позиции чтобы была возможность изменить redear на буфер или файл
//+ попытка вывести ошибку для её последующей обработки(тут вот пока не пришел к осознанию правильных практик)
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type Result struct {
	line string
	err  error
}

func NewResult(line string) Result {
	return Result{line: line}
}
func (r Result) AddError(err error) Result {
	r.err = err
	return r
}

// ReadAndSendString читает построчно данные и отправляет в канал out
// Функция реализована с позиции чтобы была возможность изменить redear на буфер или файл(например для тестов)
// + попытка вывести ошибку для её последующей обработки(тут вот пока не пришел к осознанию правильных практик)
func ReadAndSendString(rd io.Reader) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		reader := bufio.NewReader(rd)
		for {
			str, err := reader.ReadString('\n')
			resLine := NewResult(str)
			if err != nil && !errors.Is(err, io.EOF) {
				out <- resLine.AddError(fmt.Errorf("ReadAndSendString() reader.ReadString: %s", err))
				break
			}
			out <- resLine
		}
	}()
	return out
}

// Записывает данные из канала в файл
func WriteFile(fileName string, in <-chan Result) <-chan error {
	errorChan := make(chan error)
	go func() {
		defer close(errorChan)

		file, err := os.Create(fileName)
		if err != nil {
			errorChan <- fmt.Errorf("WriteFile(%s) os.Create: %s", fileName, err)
		}
		defer func() {
			err = file.Close()
			if err != nil {
				errorChan <- fmt.Errorf("WriteFile(%s) file.WriteString: %s", fileName, err)
			}
		}()

		for line := range in {
			if line.err != nil {
				errorChan <- fmt.Errorf("ReadAndSendString:%s", line.err)
			}
			_, err = file.WriteString(line.line)
			if err != nil {
				errorChan <- fmt.Errorf("WriteFile(%s) file.WriteString: %s", fileName, err)
			}
		}
	}()
	return errorChan
}
