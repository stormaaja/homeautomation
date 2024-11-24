package genericroutes

import (
	"net/http"
	"stormaaja/go-ha/data-store/store"
)

type StoreRoute struct {
	MeasurementStores []store.MeasurementStore
}

func (s StoreRoute) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s StoreRoute) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/measurements/flush" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	for _, store := range s.MeasurementStores {
		store.Flush()
	}
	w.WriteHeader(http.StatusOK)
}
