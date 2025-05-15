package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"scrapper.go/internal/model"
	"scrapper.go/internal/service"
)

const ()

type handler struct {
	storageService service.StorageService
}

func NewHandler(service service.StorageService) *handler {
	return &handler{
		storageService: service,
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

	if err := h.storageService.AddPair(r.Context(), pair.Base, pair.Quote); err != nil {
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

	pairID, err := h.storageService.GetPairID(r.Context(), pair)
	if err != nil {
		http.Error(w, "Failed to get pairs:", http.StatusInternalServerError)
	}

	rates, err := h.storageService.GetLatestRates(r.Context(), pairID)
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
