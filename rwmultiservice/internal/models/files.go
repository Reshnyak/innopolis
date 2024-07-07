package models

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
)

type FileRepo struct {
	Names map[string]struct{}
	Path  string
	Count int
	mu    *sync.RWMutex
}

// err := fl.setupFiles()
func NewFileNames(config *configs.Config) *FileRepo {
	return &FileRepo{
		mu:    new(sync.RWMutex),
		Names: make(map[string]struct{}),
		Path:  config.FilePath,
		Count: config.FilesCount,
	}
}

// создаем файлы для имитации
func (fl *FileRepo) SetupFiles() error {
	if len(fl.Names) != 0 {
		return nil
	}
	fl.Names = make(map[string]struct{})
	for i := 1; i <= fl.Count; i++ {
		filename := fmt.Sprintf("file%d.txt", i)
		fl.Names[filename] = struct{}{}
		file, err := os.Create(fl.Path + filename)
		if err != nil {
			return fmt.Errorf("setupFiles(%s) os.Create: %s", filename, err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("setupFiles(%s) close: %s", filename, err)
		}
		log.Printf("setupFiles(%s) created", filename)
	}
	log.Printf("%d files created ", len(fl.Names))
	return nil
}
func (fl *FileRepo) CreateNewFiles(fileNames ...string) error {
	for _, fileName := range fileNames {
		file, err := os.Create(fl.Path + fileName)
		if err != nil {
			return fmt.Errorf("setupFiles(%s) os.Create: %s", fileName, err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("setupFiles(%s) close: %s", fileName, err)
		}
		fl.Names[fileName] = struct{}{}
	}
	return nil
}

// Проверяем существует ли файл с таким именем
func (fl *FileRepo) Exist(fileName string) bool {
	fl.mu.RLock()
	defer fl.mu.RUnlock()
	_, ok := fl.Names[fileName]
	return ok
}
