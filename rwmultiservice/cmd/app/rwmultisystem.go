package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/retry"
	"github.com/Reshnyak/innopolis/rwmultiservice/middleware"
	"github.com/Reshnyak/innopolis/rwmultiservice/storage"
)

type RWMultiSystem struct {
	config  *configs.Config
	storage storage.Storage
	cache   *storage.CacheMsg
}

func New(c *configs.Config) (*RWMultiSystem, error) {
	st, err := storage.NewStorage(c)
	if err != nil {
		return nil, err
	}
	return &RWMultiSystem{
		config:  c,
		storage: st,
		cache:   storage.NewCache(),
	}, nil
}

// fanInValidate - проверяет сообщения из входных каналов на соответствие токену
// в дальнейшем можно проверить и на существование файла с таким именем
func fanInValidate(inputs []chan models.Message) <-chan models.Message {
	log.Printf("fanInValidate len inputs:%d", len(inputs))
	wg := new(sync.WaitGroup)
	out := make(chan models.Message)
	output := func(c <-chan models.Message) {
		for msg := range c {
			if middleware.UserTokens.IsValid(msg.Token) {
				out <- msg
			}
		}
		wg.Done()
	}
	wg.Add(len(inputs))
	for _, res := range inputs {
		go output(res)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func (rwms *RWMultiSystem) worker(wg *sync.WaitGroup, inputs <-chan []models.Message) {
	go func() {
		defer wg.Done()
		retry := retry.NewRetrie(1*time.Second, 3)
		for messages := range inputs {
			for _, msg := range messages {
				err := retry.Run(msg, rwms.storage.WriteMsg)
				if err != nil {
					log.Printf("worker can not write message in file%s: %s", msg.FileId, err)
				}
			}
		}
	}()
}
func (rwms *RWMultiSystem) startWorkers() <-chan struct{} {
	done := make(chan struct{})
	defer close(done)
	wg := new(sync.WaitGroup)
	cacheChan := make(chan []models.Message)
	go func() {
		for _, fileName := range rwms.cache.GetKeys() {
			messages, ok := rwms.cache.Get(fileName)
			if ok {
				cacheChan <- messages
			}
		}
		close(cacheChan)
	}()

	for w := 0; w < rwms.config.WorkersCount; w++ {
		wg.Add(1)
		rwms.worker(wg, cacheChan)
	}
	wg.Wait()
	return done
}
func (rwms *RWMultiSystem) Process(ctx context.Context, inputs []chan models.Message) error {
	ticker := time.NewTicker(rwms.config.WorkerDuration)
	//запустим проверку сообщений из входных каналов
	go func() {
		for msg := range fanInValidate(inputs) {
			rwms.cache.Set(msg.FileId, msg)
		}
	}()

	for {
		select {
		case <-ticker.C:
			go func() {
				rwms.startWorkers()
			}()
			fmt.Printf("ticket len cache:%d\n", rwms.cache.Len())
		case <-ctx.Done():
			fmt.Printf("shutdown:saving message from cache to files... changes for %d files are waiting to be written \n", rwms.cache.Len())
			<-rwms.startWorkers()
			fmt.Printf("shutdown:saved\n")
			return nil
		}
	}
}
