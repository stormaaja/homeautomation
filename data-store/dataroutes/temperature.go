package dataroutes

import (
	"fmt"
	"net/http"
	"stormaaja/go-ha/data-store/store"
	"strings"
)

type TemperatureRoute struct {
	Store store.DataStore
}

func ParseId(path string) string {
	splittedPath := strings.Split(path, "/")
	if len(splittedPath) < 3 {
		return ""
	}

	return splittedPath[3]
}

func IsValidValueType(path string) bool {
	splittedPath := strings.Split(path, "/")
	return len(splittedPath) == 5 && splittedPath[4] == "temperature"
}

func (t TemperatureRoute) HandleGet(w http.ResponseWriter, r *http.Request) {
	if !IsValidValueType(r.URL.Path) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	sensorId := ParseId(r.URL.Path)

	if sensorId == "" {
		http.Error(w, "Invalid sensor id", http.StatusBadRequest)
		return
	}

	// Get the temperature from the store
	temperature, success := t.Store.GetFloat(sensorId)
	if !success {
		http.Error(w, "Sensor not found", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%f", temperature)
}

func (t TemperatureRoute) HandlePost(w http.ResponseWriter, r *http.Request) {
	if !IsValidValueType(r.URL.Path) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	sensorId := ParseId(r.URL.Path)

	if sensorId == "" {
		http.Error(w, "Invalid sensor id", http.StatusBadRequest)
		return
	}

	var temperature float64
	_, error := fmt.Fscanf(r.Body, "%f", &temperature)

	if error != nil {
		http.Error(w, "Invalid temperature", http.StatusBadRequest)
		return
	}

	t.Store.SetFloat(sensorId, temperature)
	w.WriteHeader(http.StatusCreated)
}
