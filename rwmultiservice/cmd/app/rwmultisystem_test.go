package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Reshnyak/innopolis/rwmultiservice/setup"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"github.com/Reshnyak/innopolis/rwmultiservice/middleware"
	"github.com/Reshnyak/innopolis/rwmultiservice/storage"
)

func ProcessTestCase(userChanCount int, cfg *configs.Config, files []string, messages []models.Message) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rwms := &RWMultiSystem{
		config:  cfg,
		storage: storage.NewStorageFS(cfg),
		cache:   storage.NewCache(),
	}
	//coздадим файлы
	err := rwms.storage.File().CreateNewFiles(files...)
	if err != nil {
		return fmt.Errorf("processTestCase: CreateNewFiles  %s", err)
	}
	//инициалзируем каналы
	inputMsgChans := make([]<-chan models.Message, userChanCount)
	//разбросаем сообщения по каналам
	sendFunc := func(inputs []models.Message) <-chan models.Message {
		out := make(chan models.Message)
		go func() {
			for _, msg := range inputs {
				out <- msg
			}
			close(out)
		}()
		return out
	}
	sendChan := sendFunc(messages)
	for i := 0; i < userChanCount; i++ {
		inputMsgChans[i] = setup.FanOut(sendChan)
	}

	return rwms.Process(ctx, inputMsgChans)
}

// ### Сценарий 1: Успешная запись
// 1. Пользователь отправляет сообщение с правильным токеном в канал записи.
// 2. Сообщение кешируется.
// 3. Воркер через заданный интервал времени (секунда) берет сообщение из кеша.
// 4. Воркер записывает сообщение в целевой файл.
// 5. Кеш очищается для этого файла.
func TestRWSystemCase1(t *testing.T) {
	cfg, _ := configs.GetConfig("default")
	cfg.FilePath = ""
	cfg.FilesCount = 1
	cfg.WorkersCount = 1
	token, _ := middleware.UserTokens.Generate()
	files := []string{"file11.txt"}
	messages := []models.Message{{
		Token:  token,
		FileId: "file11.txt",
		Data:   "123",
	}}
	want := []string{"123"}

	err := ProcessTestCase(1, cfg, files, messages)
	if err != nil {
		t.Errorf("ProcessTestCase:%s", err)
		return
	}
	for i, file := range files {
		got, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("could not readfile:%s", err)
		}
		if !strings.Contains(string(got), want[i]) {
			t.Errorf("got msg %s want:%s", string(got), want[i])
		}
		_ = os.Remove(file)
	}
}

// ### Сценарий 2: Неверный токен
// 1. Пользователь отправляет сообщение с неправильным токеном в канал записи.
// 2. Сообщение проверяется на валидность токена.
// 3. Сообщение не кешируется и отбрасывается.
func TestRWSystemCase2(t *testing.T) {
	cfg, _ := configs.GetConfig("default")
	cfg.FilePath = ""
	cfg.FilesCount = 1
	cfg.WorkersCount = 1
	_, _ = middleware.UserTokens.Generate()
	files := []string{"file11.txt"}
	messages := []models.Message{{
		Token:  "token",
		FileId: "file11.txt",
		Data:   "123",
	}}

	err := ProcessTestCase(1, cfg, files, messages)
	if err != nil {
		t.Errorf("ProcessTestCase:%s", err)
		return
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		for _, msg := range messages {
			if msg.FileId == file && strings.Contains(string(data), msg.Data) {
				t.Errorf("validation token service is incorrect:%s", err)
			}
		}
		_ = os.Remove(file)
	}
}

//### Сценарий 3: Остановка приложения (Graceful Shutdown)
//1. Приложение получает сигнал остановки.
//2. Воркер проходит по кешу и записывает все оставшиеся сообщения в соответствующие файлы.
//3. Приложение завершает работу.

//Достаточно остановить приложение по Ctrl+C

// ### Сценарий 4: Высокая нагрузка
// 1. Пользователи массово отправляют сообщения в каналы записи.
// 2. Обработчики записывают сообщения в кеш.
// 3. Воркер масштабируется
// Большую нагрузку можно получить запустив приложение с конфигом в котором большое количество сообщений и файлов
func TestRWSystemCase4(t *testing.T) {
	cfg, _ := configs.GetConfig("default")
	cfg.FilePath = ""
	cfg.FilesCount = 4
	cfg.WorkersCount = 4
	token, _ := middleware.UserTokens.Generate()

	files := []string{"file11.txt", "file12.txt", "file13.txt", "file14.txt"}
	messages := []models.Message{
		{
			Token:  token,
			FileId: "file11.txt",
			Data:   "123",
		},
		{
			Token:  token,
			FileId: "file12.txt",
			Data:   "456",
		},
		{
			Token:  token,
			FileId: "file13.txt",
			Data:   "789",
		},
		{
			Token:  token,
			FileId: "file14.txt",
			Data:   "qwerty",
		},
	}
	want := []string{"123", "456", "789", "qwerty"}

	err := ProcessTestCase(4, cfg, files, messages)
	if err != nil {
		t.Errorf("ProcessTestCase:%s", err)
		return
	}
	for i, file := range files {
		got, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("could not readfile:%s", err)
		}
		if !strings.Contains(string(got), want[i]) {
			t.Errorf("got msg %s want:%s", string(got), want[i])
		}

		_ = os.Remove(file)
	}
}

// ### Cценарий 5: Файл с одновременной записью
// 1. Корректная запись данных в один и тот же файл, когда несколько пользователей одновременно отправляют сообщения для него.
// 2. Синхронизация, чтобы избежать конфликтов и потери данных.
func TestRWSystemCase5(t *testing.T) {

	cfg, _ := configs.GetConfig("default")
	cfg.FilePath = ""
	cfg.FilesCount = 4
	cfg.WorkersCount = 4
	token, _ := middleware.UserTokens.Generate()

	files := []string{"file11.txt"}
	messages := []models.Message{
		{
			Token:  token,
			FileId: "file11.txt",
			Data:   "123",
		},
		{
			Token:  token,
			FileId: "file11.txt",
			Data:   "456",
		},
		{
			Token:  token,
			FileId: "file11.txt",
			Data:   "789",
		},
		{
			Token:  token,
			FileId: "file11.txt",
			Data:   "qwerty",
		},
	}
	want := []string{"123", "456", "789", "qwerty"}

	err := ProcessTestCase(4, cfg, files, messages)
	if err != nil {
		t.Errorf("ProcessTestCase:%s", err)
		return
	}
	for i, file := range files {
		got, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("could not readfile:%s", err)
		}
		if !strings.Contains(string(got), want[i]) {
			t.Errorf("got msg %s want:%s", string(got), want[i])
		}
		_ = os.Remove(file)
	}
}

//### Cценарий 6: Сбор работы воркера
//1. Если воркер сталкивается с ошибкой при записи данных в файл (недоступность файла, отказ диска, etc), система должна предпринять меры по обработке этой ситуации.
//2. Ретраи, получается
// Тестировал следующим образом:
// Запускал приложение и в терминале менял путь к файлам mv data data2
// Появлялось сообщение о повторе. Возвращал в терминале верный путь
