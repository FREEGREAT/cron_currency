package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scrapper.go/internal/model"
	"scrapper.go/internal/storage"
)

const ()

type handler struct {
	currency storage.CurrencyStorage
	pairs    storage.PairStorage
}

func NewHandler(curr storage.CurrencyStorage, pairs storage.PairStorage) *handler {
	return &handler{
		currency: curr,
		pairs:    pairs,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/api/pairs/", h.GetRatesHandler)
	router.POST("/api/pairs/", h.AddPairHandler)

}

func (h *handler) AddPairHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var pair model.Pair

	if err := json.NewDecoder(r.Body).Decode(&pair); err != nil {
		http.Error(w, "Invalid input:"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.pairs.AddPair(r.Context(), pair.Base, pair.Quote); err != nil {
		http.Error(w, "Failed to add pair", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Pair added")
	return
}

func (h *handler) GetRatesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var pair model.Pair

	if err := json.NewDecoder(r.Body).Decode(&pair); err != nil {
		http.Error(w, "Invalid input:"+err.Error(), http.StatusBadRequest)
		return
	}

	pairID, err := h.pairs.GetPairID(r.Context(), pair)
	if err != nil {
		http.Error(w, "Failed to get pairs:", http.StatusInternalServerError)
	}

	rates, err := h.currency.GetLatestRates(r.Context(), pairID)
	if err != nil {
		http.Error(w, "Failed to get rates", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(rates); err != nil {
		http.Error(w, "Failed yomayo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
	return

}
