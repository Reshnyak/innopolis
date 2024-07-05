package storage

import (
	"github.com/Reshnyak/innopolis/rwmultiservice/configs"
	"github.com/Reshnyak/innopolis/rwmultiservice/internal/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWriteMsg(t *testing.T) {
	cfg := configs.DefaultInitialize()
	cfg.FilePath = ""
	file, err := os.Create("file1.txt")
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)
	store, err := NewStorage(cfg)
	assert.NoError(t, err)
	msg := models.Message{
		Token:  "123",
		FileId: "file1.txt",
		Data:   "happy lama",
	}
	err = <-store.WriteMsg(msg)
	assert.NoError(t, err)
	data, err := os.ReadFile("file1.txt")
	assert.NoError(t, err)
	assert.Equal(t, msg.Data, string(data))
	_ = os.Remove("file1.txt")
}
