package api

import (
	"encoding/json"
	"net/http"
	"request-service/internal/storage"
	"sort"
	"strings"
)

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func filterByMarket(data []storage.MarketData, market string) []storage.MarketData {
	var filtered []storage.MarketData
	for _, d := range data {
		if strings.Contains(strings.ToLower(d.Market), market) {
			filtered = append(filtered, d)
			break
		}
	}
	return filtered
}

func topAssets(data []storage.MarketData, n int) []storage.MarketData {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Price > data[j].Price
	})
	if len(data) > n {
		return data[:n]
	}
	return data
}
