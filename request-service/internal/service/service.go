package service

import (
	"request-service/internal/storage"
)

type MarketService struct {
	storage  *storage.Storage
	response []ResponseMarketData
}

func NewMarketService(storage *storage.Storage) *MarketService {
	return &MarketService{
		storage:  storage,
		response: make([]ResponseMarketData, 0),
	}
}

func (s *MarketService) GetLatestData() ([]ResponseMarketData, error) {
	lastData, err := s.storage.GetLatestMarketData()
	if err != nil {
		return nil, err
	}

	s.response = make([]ResponseMarketData, 0)
	for _, d := range lastData {
		s.response = append(s.response, ResponseMarketData{
			Exchange:           d.Exchange,
			Symbol:             d.Symbol,
			Market:             d.Market,
			Price:              float64(d.Price) / 1e3,
			Volume:             float64(d.Volume) / 1e3,
			High:               float64(d.High) / 1e3,
			Low:                float64(d.Low) / 1e3,
			PriceChangePercent: d.PriceChangePercent + "%",
			Timestamp:          d.Timestamp,
		})
	}

	return s.response, nil
}
