package setup

import (
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"math/rand"
	"sync"
	"time"
)

const delayCoef = 500

func FanOut(input <-chan models.Message) chan models.Message {
	out := make(chan models.Message)
	go func() {
		for v := range input {
			out <- v
		}
		close(out)
	}()
	return out
}

// SetupProcess имитация отправки сообщений пользователями. Пользователи сгенерированы в storage
// на каждого пользователя создается канал...лучше сделать n воркеров, ведь пользователей может быть очень много
// сообщения от пользователей отправляются с разной задержкой по времени
func SetupProcess(users []models.User) ([]<-chan models.Message, error) {

	//инициалзируем каналы
	inputMsgChans := make([]<-chan models.Message, len(users))
	//разбросаем сообщения по каналам
	sendFunc := func(inputs []models.User) <-chan models.Message {
		out := make(chan models.Message)
		wg := new(sync.WaitGroup)

		for _, user := range inputs {
			wg.Add(1)
			go func(w *sync.WaitGroup, u models.User) {
				defer w.Done()
				for _, msg := range u.Messages {
					time.Sleep(time.Duration(rand.Int()%delayCoef) * time.Millisecond)
					out <- msg
				}
			}(wg, user)
		}
		go func() {
			wg.Wait()
			close(out)
		}()

		return out
	}
	sendChan := sendFunc(users)
	for i := 0; i < len(users); i++ {
		inputMsgChans[i] = FanOut(sendChan)
	}
	return inputMsgChans, nil
}
