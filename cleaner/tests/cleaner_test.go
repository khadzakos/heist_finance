package tests

import (
	"context"
	"testing"
	"time"

	"cleaner/internal/cleaner"
	"cleaner/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCleaner(t *testing.T) {
	// Create a test storage
	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	// Test cleaner creation
	cleaner, err := cleaner.NewCleaner(storage)
	require.NoError(t, err)
	assert.NotNil(t, cleaner)
	assert.NotNil(t, cleaner.Db)
}

func TestCleaner_CleanOldData(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	cleaner, err := cleaner.NewCleaner(storage)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	go cleaner.CleanOldData()

	select {
	case <-ctx.Done():
		return
	}
}
