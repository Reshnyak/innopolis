package models

import (
	"fmt"
	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/middleware"
	"math/rand"
)

type User struct {
	UserToken string
	Messages  []Message
	count     int
}

var text = []byte("1234567890     ABCDEFGHIKLMNOPQRSTVXYZ")

// Create users with message and token
func SetupUsers(fl *FileRepo, config *configs.Config) ([]User, error) {
	res := make([]User, config.UsersCount)
	for u := 0; u < config.UsersCount; u++ {
		token, err := middleware.UserTokens.Generate()
		if err != nil {
			return nil, fmt.Errorf("could not create user: %s", err)
		}
		messages, err := setupMessages(token, fl, config)
		if err != nil {
			return nil, fmt.Errorf("could not create messages: %s", err)
		}
		res[u] = User{
			UserToken: token,
			Messages:  messages,
		}
	}
	return res, nil
}

// Create random message
func setupMessages(token string, fl *FileRepo, config *configs.Config) ([]Message, error) {
	currentCount := 1 + rand.Int()%(config.MessageMaxCount-1)
	messages := make([]Message, currentCount)
	for i := 0; i < len(messages); i++ {

		data := make([]byte, 1+rand.Int()%(config.MessageMaxLen-1))
		rand.Shuffle(len(text), func(i, j int) {
			text[i], text[j] = text[j], text[i]
		})
		copy(data, text)
		data = append(data, []byte("\n")...)
		messages[i] = Message{
			Token:  token,
			FileId: fl.getRand(),
			Data:   string(data),
		}
	}
	return messages, nil
}
