package store

type MockMeasurementStore struct {
	Items map[string]float64
}

func (m *MockMeasurementStore) AppendItem(
	measurement string,
	location string,
	field string,
	value float64,
) {
	m.Items[location] = value
}

func (m *MockMeasurementStore) Flush() {
	m.Items = make(map[string]float64)
}
