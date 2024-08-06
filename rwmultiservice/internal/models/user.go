package models

import (
	"fmt"
	"math/rand"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/middleware"
	"github.com/Reshnyak/innopolis/rwmultiservice/utils"
)

type Users []User

type User struct {
	UserToken string
	Messages  []Message
}

func NewUsers() *Users {
	return &Users{}
}

// Create users with message and token
func (u *Users) SetupUsers(config *configs.Config) error {
	*u = make(Users, config.UsersCount)
	for i := 0; i < config.UsersCount; i++ {
		token, err := middleware.UserTokens.Generate()
		if err != nil {
			return fmt.Errorf("could not create user: %s", err)
		}
		messages, err := setupMessages(token, config)
		if err != nil {
			return fmt.Errorf("could not create messages: %s", err)
		}
		(*u)[i] = User{
			UserToken: token,
			Messages:  messages,
		}
	}
	return nil
}

// Create random message
func setupMessages(token string, config *configs.Config) ([]Message, error) {
	currentCount := 1 + rand.Int()%(config.MessageMaxCount-1)
	messages := make([]Message, currentCount)
	for i := 0; i < len(messages); i++ {
		messages[i] = Message{
			Token:  token,
			FileId: fmt.Sprintf("file%d.txt", 1+rand.Int()%(config.FilesCount-1)),
			Data:   utils.GetRandomString(config.MessageMaxLen),
		}
	}
	return messages, nil
}
