package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
)

type UserReposytory interface {
	RegisterNewUser(users ...models.User)
	GetAllUsers() []models.User
	User() *models.Users
}
type FileReposytory interface {
	WriteMsg(msg models.Message) <-chan error
	GetAllFiles() ([]string, error)
	File() *models.FileRepo
}

type Storage interface {
	FileReposytory
	UserReposytory
}

func NewStorage(config *configs.Config) (Storage, error) {
	switch config.StorageType {
	case "FS":
		return NewStorageFS(config), nil
	case "DB":
		return nil, errors.New("we apologize that the storageDb is not implemented now. we are working on it")
	}
	return nil, errors.New("could not create or open storage")
}

type StorageFS struct {
	mu    *sync.Mutex
	path  string
	files *models.FileRepo
	users *models.Users
}

func NewStorageFS(config *configs.Config) *StorageFS {
	return &StorageFS{
		files: models.NewFileNames(config),
		path:  config.FilePath,
		users: models.NewUsers(),
		mu:    new(sync.Mutex),
	}
}
func (sfs StorageFS) File() *models.FileRepo {
	return sfs.files
}
func (sfs StorageFS) User() *models.Users {
	return sfs.users
}
func (sfs StorageFS) WriteMsg(msg models.Message) <-chan error {
	done := make(chan error)
	go func(mu *sync.Mutex) {
		defer close(done)
		mu.Lock()
		defer mu.Unlock()
		file, err := os.OpenFile(sfs.path+msg.FileId, os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			done <- fmt.Errorf("WriteMsg(%s) os.Open: %s", msg.FileId, err)
			return
		}
		defer func() {
			done <- file.Close()
		}()
		_, err = file.WriteString(msg.Data)
		if err != nil {
			done <- fmt.Errorf("WriteMsg(%s) file.WriteString: %s", msg.FileId, err)
			return
		}

		log.Printf("write file is done:%s", msg.FileId)

	}(sfs.mu)

	return done
}
func (sfs *StorageFS) GetAllFiles() ([]string, error) {

	files := make([]string, 0, len(sfs.files.Names))
	//На данном этапе не требуется Lock, но если добавим добавление новых файлов
	sfs.mu.Lock()
	for f := range sfs.files.Names {
		files = append(files, f)
	}
	sfs.mu.Unlock()
	sort.Strings(files)
	return files, nil
}
func (sfs *StorageFS) RegisterNewUser(users ...models.User) {
	sfs.mu.Lock()
	defer sfs.mu.Unlock()
	*sfs.users = append(*sfs.users, users...)

}

func (sfs *StorageFS) GetAllUsers() []models.User {
	sfs.mu.Lock()
	defer sfs.mu.Unlock()
	countMessage := 0
	for _, usr := range *sfs.users {
		countMessage += len(usr.Messages)
	}
	log.Printf("Users count:%d Message count:%d\n", len(*sfs.users), countMessage)
	return *sfs.users
}

// ////////////////////////////////////////////////////////
type ConfigDb struct {
}
type StorageDB struct {
	config *ConfigDb
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
