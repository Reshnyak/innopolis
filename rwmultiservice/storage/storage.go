package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"log"
	"os"
	"sort"
	"sync"
)

type Storage interface {
	WriteMsg(msg models.Message) <-chan error
	GetAllFiles() ([]string, error)
	RegisterNewUser() (string, error)
	GetAllUsers() []models.User
}

func NewStorage(config *configs.Config) (Storage, error) {
	switch config.StorageType {
	case "FS":
		fl, err := models.NewFileNames(config)
		if err != nil {
			return nil, err
		}
		users, err := models.SetupUsers(fl, config)
		if err != nil {
			return nil, err
		}
		return &StorageFS{
			files: fl,
			path:  config.FilePath,
			users: users,
			mu:    new(sync.Mutex),
		}, nil
	case "DB":
		return nil, errors.New("we apologize that the storageDb is not implemented now. we are working on it")
	}
	return nil, errors.New("could not create or open storage")
}

type StorageFS struct {
	mu    *sync.Mutex
	path  string
	files *models.FileRepo
	users []models.User
}

func (sfs StorageFS) WriteMsg(msg models.Message) <-chan error {
	done := make(chan error)
	go func(mu *sync.Mutex) {
		defer close(done)
		mu.Lock()
		defer mu.Unlock()
		file, err := os.OpenFile("data/"+msg.FileId, os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			done <- fmt.Errorf("WriteMsg(%s) os.Open: %s", msg.FileId, err)
		}
		_, err = file.WriteString(msg.Data)
		if err != nil {
			done <- fmt.Errorf("WriteMsg(%s) file.WriteString: %s", msg.FileId, err)
		}
		done <- file.Close()
		log.Printf("write file is done:%s", msg.FileId)

	}(sfs.mu)

	return done
}
func (sfs *StorageFS) GetAllFiles() ([]string, error) {

	files := make([]string, 0, len(sfs.files.Names))
	//На данном этапе не требуется Lock, но если добавим POST на добавление нового файла
	sfs.mu.Lock()
	for f := range sfs.files.Names {
		files = append(files, f)
	}
	sfs.mu.Unlock()
	sort.Strings(files)
	return files, nil
}
func (sfs *StorageFS) RegisterNewUser() (string, error) {

	return "", nil
}

func (sfs *StorageFS) GetAllUsers() []models.User {

	countMessage := 0
	for _, usr := range sfs.users {
		countMessage += len(usr.Messages)
	}
	log.Printf("Users count:%d Message count:%d\n", len(sfs.users), countMessage)
	return sfs.users
}

// ////////////////////////////////////////////////////////
type ConfigDb struct {
}
type StorageDB struct {
	config *ConfigDb
	cache  *CacheMsg
	DB     *sql.DB
}

func (sdb *StorageDB) WriteMsg(msg models.Message) <-chan error {

	return nil
}
func (sdb *StorageDB) GetAllFiles() ([]string, error) {

	return []string{}, nil
}
func (sdb *StorageDB) RegisterNewUser() (string, error) {

	return "", nil
}
func (sdb *StorageDB) GetAllUsers() []models.User {
	return nil
}
