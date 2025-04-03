package api

import (
	"encoding/json"
	"net/http"
	"request-service/internal/service"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func filterByMarket(data []service.ResponseMarketData, market string) []service.ResponseMarketData {
	var filtered []service.ResponseMarketData
	for _, d := range data {
		if d.Market == market {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// func mainPageAssets(data []service.ResponseMarketData) []service.ResponseMarketData {

// }
