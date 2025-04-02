package service

import (
	"request-service/internal/storage"
)

type MarketService struct {
	storage *storage.Storage
}

func NewMarketService(s *storage.Storage) *MarketService {
	return &MarketService{storage: s}
}

func (s *MarketService) GetLatestData() ([]storage.MarketData, error) {
	return s.storage.GetLatestMarketData()
}
