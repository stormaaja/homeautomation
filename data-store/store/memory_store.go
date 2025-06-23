package store

type MemoryStore struct {
	Data map[string]map[string]any
}

func (m MemoryStore) GetMeasurement(
	measurementType string,
	key string,
) (Measurement, bool) {
	measurement, ok := m.Data[key][measurementType].(Measurement)
	return measurement, ok
}

func (m *MemoryStore) SetMeasurement(
	measurementType string,
	key string,
	measurement Measurement,
) {
	if m.Data[key] == nil {
		m.Data[key] = make(map[string]any)
	}
	m.Data[key][measurementType] = measurement
}
