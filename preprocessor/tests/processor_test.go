package tests

import (
	"testing"

	"preprocessor/internal/config"
	"preprocessor/internal/processor"
	"preprocessor/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProcessor(t *testing.T) {
	cfg := &config.Config{
		Preprocessor: config.PreprocessorConfig{
			Queue: "test-queue",
		},
	}

	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	processor, err := processor.NewProcessor(cfg, storage)
	require.NoError(t, err)
	assert.NotNil(t, processor)

	worker := processor.AddWorker()
	assert.NotNil(t, worker)
	assert.Equal(t, 0, worker.Id)
}

func TestProcessor_ProcessWorker(t *testing.T) {
	cfg := &config.Config{
		Preprocessor: config.PreprocessorConfig{
			Queue: "test-queue",
		},
	}

	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	processor, err := processor.NewProcessor(cfg, storage)
	require.NoError(t, err)

	worker := processor.AddWorker()
	assert.NotNil(t, worker)
	assert.Equal(t, 0, worker.Id)
	assert.NotNil(t, worker.Jobs)
	assert.NotNil(t, worker.Db)
	assert.NotNil(t, worker.Processor)
}

func TestProcessor_RemoveWorker(t *testing.T) {
	cfg := &config.Config{
		Preprocessor: config.PreprocessorConfig{
			Queue: "test-queue",
		},
	}

	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	processor, err := processor.NewProcessor(cfg, storage)
	require.NoError(t, err)

	worker := processor.AddWorker()
	assert.NotNil(t, worker)
	processor.RemoveWorker(worker.Id)

	newWorker := processor.AddWorker()
	assert.NotNil(t, newWorker)
	assert.Equal(t, 0, newWorker.Id)
}
