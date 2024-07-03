package setup

import (
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"math/rand"
	"time"
)

const delaySend = 1500

// SetupProcess имитация отправки сообщений пользователями. Пользователи сгенерированы в storage
// на каждого пользователя создается канал...лучше сделать n воркеров, но имитация не самое главное в задании =)
// сообщения от пользователей отправляются с разной задержкой
func SetupProcess(users []models.User) ([]chan models.Message, error) {
	inputMsgChans := make([]chan models.Message, len(users))
	for i := 0; i < len(users); i++ {
		inputMsgChans[i] = make(chan models.Message)
	}
	for i, user := range users {
		go func(i int, u models.User) {
			defer close(inputMsgChans[i])
			for _, msg := range u.Messages {
				time.Sleep(time.Duration(rand.Int()%delaySend) * time.Millisecond)
				inputMsgChans[i] <- msg
			}
		}(i, user)
	}

	return inputMsgChans, nil
}
