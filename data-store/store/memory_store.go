package store

type MemoryStore struct {
	Data map[string]map[string]Measurement
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
