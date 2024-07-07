package setup

import (
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"math/rand"
	"time"
)

const delayCoef = 1500

func FanOut(input <-chan models.Message) <-chan models.Message {
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
// сообщения от пользователей отправляются с разной задержкой
func SetupProcess(users []models.User) ([]<-chan models.Message, error) {
	//инициалзируем каналы
	inputMsgChans := make([]<-chan models.Message, len(users))
	//разбросаем сообщения по каналам
	sendFunc := func(inputs []models.User) <-chan models.Message {
		out := make(chan models.Message)
		go func() {
			for _, user := range inputs {
				for _, msg := range user.Messages {
					time.Sleep(time.Duration(rand.Int()%delayCoef) * time.Millisecond)
					out <- msg
				}
			}
			close(out)
		}()
		return out
	}
	for i := 0; i < len(users); i++ {
		inputMsgChans[i] = FanOut(sendFunc(users))
	}
	return inputMsgChans, nil
}
