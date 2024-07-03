package models

import (
	"fmt"
	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"log"
	"math/rand"
	"os"
	"sync"
)

var defaulFileName = "default.txt"

type FileRepo struct {
	Names map[string]struct{}
	Path  string
	Count int
	mu    sync.RWMutex
}

func NewFileNames(config *configs.Config) (*FileRepo, error) {
	var fl FileRepo
	fl.Path = config.FilePath
	fl.Count = config.FilesCount
	err := fl.setupFiles()
	return &fl, err
}

// создаем файлы для имитации
func (fl *FileRepo) setupFiles() error {
	if len(fl.Names) != 0 {
		return nil
	}
	fl.Names = make(map[string]struct{})
	for i := 1; i <= fl.Count; i++ {
		filename := fmt.Sprintf("file%d.txt", i)
		fl.Names[filename] = struct{}{}
		file, err := os.Create("data/" + filename)
		if err != nil {
			return fmt.Errorf("setupFiles(%s) os.Create: %s", filename, err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("setupFiles(%s) close: %s", filename, err)
		}
	}
	log.Printf("%d files created ", len(fl.Names))
	return nil
}

// Проверяем существует ли файл с таким именем
func (fl *FileRepo) Exist(fileName string) bool {
	fl.mu.RLock()
	defer fl.mu.RUnlock()
	_, ok := fl.Names[fileName]
	return ok
}

// return random file name or default
func (fl *FileRepo) getRand() string {
	if len(fl.Names) == 0 {
		_ = fl.setupFiles()
	}
	i, ind := 0, rand.Int()%(len(fl.Names)-1)
	fl.mu.RLock()
	defer fl.mu.RUnlock()
	for key, _ := range fl.Names {
		if i == ind {
			return key
		}
		i++
	}
	return defaulFileName
}
