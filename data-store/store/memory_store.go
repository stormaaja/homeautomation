package store

import (
	"log"
	"os"
	"stormaaja/go-ha/data-store/tools"
)

type MemoryStore struct {
	Data          map[string]map[string]Measurement
	BackupEnabled bool
}

func (m *MemoryStore) LoadMemoryStore() error {
	err := tools.ReadJsonFile("memory-store.json", &m.Data)
	if err != nil {
		if err != os.ErrNotExist {
			log.Printf("Error loading memory store: %v", err)
		} else {
			log.Printf("Memory store file does not exist, starting with empty store")
		}
	}
	return nil
}

func (m MemoryStore) SaveMemoryStore() error {
	return tools.WriteJsonFile("memory-store.json", &m.Data)
}

func (m MemoryStore) GetMeasurement(
	measurementType string,
	key string,
) (Measurement, bool) {
	measurement, ok := m.Data[key][measurementType]
	return measurement, ok
}

func (m *MemoryStore) SetMeasurement(
	measurementType string,
	key string,
	measurement Measurement,
) {
	if m.Data[key] == nil {
		m.Data[key] = make(map[string]Measurement)
	}
	m.Data[key][measurementType] = measurement
	if m.BackupEnabled {
		err := m.SaveMemoryStore()
		if err != nil {
			log.Printf("Failed to backup memory store: %v", err)
		}
	}
}

func (m MemoryStore) FindMeasurements(
	queryParams map[string]string,
) []Measurement {
	measurements := []Measurement{}
	for _, deviceMeasurements := range m.Data {
		for _, measurement := range deviceMeasurements {
			if measurement.Matches(queryParams) {
				measurements = append(measurements, measurement)
			}
		}
	}
	return measurements
}
