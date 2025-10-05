package store

import (
	"testing"
)

func TestMemoryStore_GetMeasurement(t *testing.T) {
	store := MemoryStore{
		Data: map[string]map[string]Measurement{
			"key1": {
				"temperature": Measurement{Value: 23.5},
			},
		},
	}

	measurement, ok := store.GetMeasurement("temperature", "key1")
	if !ok {
		t.Errorf("expected measurement to be found")
	}
	if measurement.Value != 23.5 {
		t.Errorf("expected measurement value to be 23.5, got %v", measurement.Value)
	}

	_, ok = store.GetMeasurement("humidity", "key1")
	if ok {
		t.Errorf("expected measurement to not be found")
	}
}

func TestMemoryStore_SetMeasurement(t *testing.T) {
	store := MemoryStore{
		Data: make(map[string]map[string]Measurement),
	}

	measurement := Measurement{Value: 23.5}
	store.SetMeasurement("temperature", "key1", measurement)

	storedMeasurement, ok := store.GetMeasurement("temperature", "key1")
	if !ok {
		t.Errorf("expected measurement to be found")
	}
	if storedMeasurement.Value != 23.5 {
		t.Errorf("expected measurement value to be 23.5, got %v", storedMeasurement.Value)
	}
}
