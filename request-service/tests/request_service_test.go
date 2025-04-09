package tests

import (
	"testing"

	"request-service/internal/service"
	"request-service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMarketService(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	service := service.NewMarketService(storage)
	assert.NotNil(t, service)
}

func TestMarketService_GetLatestData(t *testing.T) {
	dbURL := "postgres://postgres:postgres@localhost:5432/heist?sslmode=disable"
	storage, err := storage.NewStorage(dbURL)
	require.NoError(t, err)
	defer storage.Close()

	service := service.NewMarketService(storage)

	data, err := service.GetLatestData()
	require.NoError(t, err)
	assert.NotNil(t, data)

	for _, d := range data {
		assert.NotEmpty(t, d.Exchange)
		assert.NotEmpty(t, d.Symbol)
		assert.NotEmpty(t, d.Market)
		assert.Greater(t, d.Price, 0.0)
		assert.Greater(t, d.Volume, 0.0)
		assert.Greater(t, d.High, 0.0)
		assert.Greater(t, d.Low, 0.0)
		assert.NotEmpty(t, d.PriceChangePercent)
		assert.NotZero(t, d.Timestamp)
	}
}
